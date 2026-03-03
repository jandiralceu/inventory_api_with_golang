package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/apperrors"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestInventoryService_Create(t *testing.T) {
	repo := new(MockInventoryRepository)
	productRepo := new(MockProductRepository)
	warehouseRepo := new(MockWarehouseRepository)
	transactionRepo := new(MockInventoryTransactionRepository)
	cache := new(MockCacheManager)
	svc := NewInventoryService(repo, productRepo, warehouseRepo, transactionRepo, cache)

	productID := uuid.New()
	warehouseID := uuid.New()
	req := dto.CreateInventoryRequest{
		ProductID:    productID,
		WarehouseID:  warehouseID,
		Quantity:     100,
		LocationCode: "A-1",
	}

	t.Run("Success", func(t *testing.T) {
		productRepo.On("FindByID", mock.Anything, productID).Return(&models.Product{ID: productID}, nil).Once()
		warehouseRepo.On("FindByID", mock.Anything, warehouseID).Return(&models.Warehouse{ID: warehouseID}, nil).Once()
		repo.On("Create", mock.Anything, mock.MatchedBy(func(i *models.Inventory) bool {
			return i.ProductID == productID && i.WarehouseID == warehouseID && i.Quantity == 100
		})).Return(nil).Once()
		cache.On("DeletePrefix", mock.Anything, "inventory:").Return(nil).Once()

		inv, err := svc.Create(context.Background(), req)

		assert.NoError(t, err)
		assert.NotNil(t, inv)
		assert.Equal(t, 1, inv.Version)
		repo.AssertExpectations(t)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		productRepo.On("FindByID", mock.Anything, productID).Return((*models.Product)(nil), apperrors.ErrNotFound).Once()

		inv, err := svc.Create(context.Background(), req)

		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, inv)
	})

	t.Run("WarehouseNotFound", func(t *testing.T) {
		productRepo.On("FindByID", mock.Anything, productID).Return(&models.Product{ID: productID}, nil).Once()
		warehouseRepo.On("FindByID", mock.Anything, warehouseID).Return((*models.Warehouse)(nil), apperrors.ErrNotFound).Once()

		inv, err := svc.Create(context.Background(), req)

		assert.ErrorIs(t, err, apperrors.ErrNotFound)
		assert.Nil(t, inv)
	})
}

func TestInventoryService_FindByID(t *testing.T) {
	repo := new(MockInventoryRepository)
	productRepo := new(MockProductRepository)
	warehouseRepo := new(MockWarehouseRepository)
	transactionRepo := new(MockInventoryTransactionRepository)
	cache := new(MockCacheManager)
	svc := NewInventoryService(repo, productRepo, warehouseRepo, transactionRepo, cache)

	id := uuid.New()
	expected := &models.Inventory{ID: id, Quantity: 50}

	t.Run("CacheHit", func(t *testing.T) {
		cache.On("Get", mock.Anything, fmt.Sprintf("inventory:id:%s", id), mock.AnythingOfType("*models.Inventory")).
			Run(func(args mock.Arguments) {
				dest := args.Get(2).(*models.Inventory)
				*dest = *expected
			}).Return(nil).Once()

		inv, err := svc.FindByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, expected.Quantity, inv.Quantity)
		repo.AssertNotCalled(t, "FindByID", mock.Anything, mock.Anything)
	})

	t.Run("CacheMiss", func(t *testing.T) {
		cache.On("Get", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("not found")).Once()
		repo.On("FindByID", mock.Anything, id).Return(expected, nil).Once()
		cache.On("Set", mock.Anything, mock.Anything, expected, 15*time.Minute).Return(nil).Once()

		inv, err := svc.FindByID(context.Background(), id)

		assert.NoError(t, err)
		assert.Equal(t, expected.ID, inv.ID)
	})
}

func TestInventoryService_StockOperations(t *testing.T) {
	repo := new(MockInventoryRepository)
	productRepo := new(MockProductRepository)
	warehouseRepo := new(MockWarehouseRepository)
	transactionRepo := new(MockInventoryTransactionRepository)
	cache := new(MockCacheManager)
	svc := NewInventoryService(repo, productRepo, warehouseRepo, transactionRepo, cache)

	id := uuid.New()
	existing := &models.Inventory{ID: id, Quantity: 10, ReservedQuantity: 2, Version: 1}

	t.Run("AddStock_Success", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, id).Return(existing, nil).Once()
		repo.On("UpdateStock", mock.Anything, id, 5, existing.Version).Return(nil).Once()
		cache.On("DeletePrefix", mock.Anything, "inventory:").Return(nil).Once()

		err := svc.AddStock(context.Background(), id, 5)

		assert.NoError(t, err)
	})

	t.Run("RemoveStock_Insufficient", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, id).Return(existing, nil).Once()

		err := svc.RemoveStock(context.Background(), id, 20)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
	})

	t.Run("ReserveStock_Success", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, id).Return(existing, nil).Once()
		repo.On("UpdateReservedStock", mock.Anything, id, 3, existing.Version).Return(nil).Once()
		cache.On("DeletePrefix", mock.Anything, "inventory:").Return(nil).Once()

		err := svc.ReserveStock(context.Background(), id, 3)

		assert.NoError(t, err)
	})

	t.Run("ReleaseStock_Excessive", func(t *testing.T) {
		repo.On("FindByID", mock.Anything, id).Return(existing, nil).Once()

		err := svc.ReleaseStock(context.Background(), id, 10)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot release more than reserved")
	})
}

func TestInventoryService_GetTransactionHistory(t *testing.T) {
	repo := new(MockInventoryRepository)
	productRepo := new(MockProductRepository)
	warehouseRepo := new(MockWarehouseRepository)
	transactionRepo := new(MockInventoryTransactionRepository)
	cache := new(MockCacheManager)
	svc := NewInventoryService(repo, productRepo, warehouseRepo, transactionRepo, cache)

	t.Run("Success", func(t *testing.T) {
		req := dto.TransactionListRequest{
			PaginationRequest: dto.PaginationRequest{Limit: 10, Page: 1},
		}

		productID := uuid.New()
		userID := uuid.New()
		mockTransactions := []models.InventoryTransaction{
			{
				ID:              uuid.New(),
				ProductID:       productID,
				UserID:          &userID,
				TransactionType: "ADJUSTMENT",
				Product:         &models.Product{ID: productID, Name: "Test Product"},
				User:            &models.User{ID: userID, Name: "Test User"},
			},
		}

		transactionRepo.On("FindAll", mock.Anything, mock.MatchedBy(func(f repository.TransactionListFilter) bool {
			return f.Pagination.Limit == 10
		})).Return(mockTransactions, int64(1), nil).Once()

		resp, err := svc.GetTransactionHistory(context.Background(), req)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), resp.Total)
		assert.Len(t, resp.Data, 1)
		assert.Equal(t, "ADJUSTMENT", resp.Data[0].TransactionType)
		assert.NotNil(t, resp.Data[0].Product)
		assert.Equal(t, "Test Product", resp.Data[0].Product.Name)
		assert.NotNil(t, resp.Data[0].User)
		assert.Equal(t, "Test User", resp.Data[0].User.Name)
		transactionRepo.AssertExpectations(t)
	})
}
