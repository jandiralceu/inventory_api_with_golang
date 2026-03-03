package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jandiralceu/inventory_api_with_golang/internal/dto"
	"github.com/jandiralceu/inventory_api_with_golang/internal/models"
	"github.com/jandiralceu/inventory_api_with_golang/internal/pkg"
	"github.com/jandiralceu/inventory_api_with_golang/internal/repository"
)

// InventoryService defines the business logic contract for inventory management.
// It orchestrates operations between repositories and handles cache invalidation.
type InventoryService interface {
	// Create registers a new inventory record for a product in a warehouse.
	// Validates if product and warehouse exist first.
	Create(ctx context.Context, req dto.CreateInventoryRequest) (*models.Inventory, error)
	// Update modifies non-quantity fields of an existing inventory record.
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (*models.Inventory, error)
	// Delete removes an inventory record and purges related cache entries.
	Delete(ctx context.Context, id uuid.UUID) error
	// FindByID retrieves a single inventory record, using cache when available.
	FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error)
	// FindAll returns a paginated list of inventory records with optional filters.
	FindAll(ctx context.Context, req dto.GetInventoryListRequest) (*dto.InventoryListResponse, error)
	// AddStock increases the physical quantity of an item.
	AddStock(ctx context.Context, id uuid.UUID, quantity int) error
	// RemoveStock decreases the physical quantity of an item after verifying availability.
	RemoveStock(ctx context.Context, id uuid.UUID, quantity int) error
	// ReserveStock earmarks a quantity for a future transaction without decreasing physical stock.
	ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
	// ReleaseStock removes a reservation, making the quantity available again.
	ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error
	// GetTransactionHistory retrieves a paginated audit log of stock movements.
	GetTransactionHistory(ctx context.Context, req dto.TransactionListRequest) (dto.TransactionListResponse, error)
}

type inventoryService struct {
	repo            repository.InventoryRepository
	productRepo     repository.ProductRepository
	warehouseRepo   repository.WarehouseRepository
	transactionRepo repository.InventoryTransactionRepository
	cache           pkg.CacheManager
}

const (
	_inventoryCachePrefix = "inventory:"
	_inventoryCacheTTL    = 15 * time.Minute
)

// NewInventoryService initializes the inventory service with its dependencies.
func NewInventoryService(
	repo repository.InventoryRepository,
	productRepo repository.ProductRepository,
	warehouseRepo repository.WarehouseRepository,
	transactionRepo repository.InventoryTransactionRepository,
	cache pkg.CacheManager,
) InventoryService {
	return &inventoryService{
		repo:            repo,
		productRepo:     productRepo,
		warehouseRepo:   warehouseRepo,
		transactionRepo: transactionRepo,
		cache:           cache,
	}
}

// Create registers a new inventory record for a product in a warehouse.
// It validates if product and warehouse exist first and logs an initial OPENING transaction if quantity > 0.
func (s *inventoryService) Create(ctx context.Context, req dto.CreateInventoryRequest) (*models.Inventory, error) {
	// Check if product exists
	if _, err := s.productRepo.FindByID(ctx, req.ProductID); err != nil {
		return nil, err
	}

	// Check if warehouse exists
	if _, err := s.warehouseRepo.FindByID(ctx, req.WarehouseID); err != nil {
		return nil, err
	}

	inventory := &models.Inventory{
		ProductID:        req.ProductID,
		WarehouseID:      req.WarehouseID,
		Quantity:         req.Quantity,
		ReservedQuantity: 0,
		LocationCode:     req.LocationCode,
		MinQuantity:      req.MinQuantity,
		MaxQuantity:      req.MaxQuantity,
		Metadata:         req.Metadata,
		Version:          1,
	}

	if err := s.repo.Create(ctx, inventory); err != nil {
		return nil, err
	}

	// Record initial transaction if quantity > 0
	if req.Quantity > 0 {
		_ = s.logTransaction(ctx, inventory, req.Quantity, "OPENING", "Initial stock setup")
	}

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return inventory, nil
}

// Update modifies non-quantity fields of an existing inventory record.
// If quantity is updated directly, it logs an ADJUSTMENT transaction.
func (s *inventoryService) Update(ctx context.Context, id uuid.UUID, req dto.UpdateInventoryRequest) (*models.Inventory, error) {
	inventory, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.LocationCode != "" {
		inventory.LocationCode = req.LocationCode
	}
	if req.MinQuantity != nil {
		inventory.MinQuantity = *req.MinQuantity
	}
	if req.MaxQuantity != nil {
		inventory.MaxQuantity = req.MaxQuantity
	}
	if req.Metadata != nil {
		inventory.Metadata = req.Metadata
	}
	// Quantity and ReservedQuantity are typically updated via specific methods (AddStock, etc.)
	// But we can allow direct updates if needed
	var qtyDiff int
	if req.Quantity != nil {
		qtyDiff = *req.Quantity - inventory.Quantity
		inventory.Quantity = *req.Quantity
	}
	if req.ReservedQuantity != nil {
		inventory.ReservedQuantity = *req.ReservedQuantity
	}

	if err := s.repo.Update(ctx, inventory); err != nil {
		return nil, err
	}

	// Record adjustment if quantity changed directly
	if qtyDiff != 0 {
		_ = s.logTransaction(ctx, inventory, qtyDiff, "ADJUSTMENT", "Manual update in Inventory record")
	}

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return inventory, nil
}

// Delete removes an inventory record and purges related cache entries.
func (s *inventoryService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return nil
}

// FindByID retrieves a single inventory record, using cache when available.
func (s *inventoryService) FindByID(ctx context.Context, id uuid.UUID) (*models.Inventory, error) {
	cacheKey := fmt.Sprintf("%sid:%s", _inventoryCachePrefix, id)
	var inventory models.Inventory

	if err := s.cache.Get(ctx, cacheKey, &inventory); err == nil {
		return &inventory, nil
	}

	inv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	_ = s.cache.Set(ctx, cacheKey, inv, _inventoryCacheTTL)
	return inv, nil
}

// FindAll returns a paginated list of inventory records with optional filters.
func (s *inventoryService) FindAll(ctx context.Context, req dto.GetInventoryListRequest) (*dto.InventoryListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.Limit
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filter := repository.InventoryListFilter{
		ProductID:   req.ProductID,
		WarehouseID: req.WarehouseID,
		LowStock:    req.LowStock,
		Page:        page,
		Limit:       limit,
		Sort:        req.Sort,
		Order:       req.Order,
	}

	inventories, total, err := s.repo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.InventoryResponse, len(inventories))
	for i, inv := range inventories {
		responses[i] = dto.InventoryResponse{
			ID:                inv.ID,
			ProductID:         inv.ProductID,
			WarehouseID:       inv.WarehouseID,
			Quantity:          inv.Quantity,
			ReservedQuantity:  inv.ReservedQuantity,
			AvailableQuantity: inv.AvailableQuantity(),
			LocationCode:      inv.LocationCode,
			MinQuantity:       inv.MinQuantity,
			MaxQuantity:       inv.MaxQuantity,
			Version:           inv.Version,
			LastCountedAt:     inv.LastCountedAt,
			Metadata:          inv.Metadata,
			CreatedAt:         inv.CreatedAt,
			UpdatedAt:         inv.UpdatedAt,
			Product:           inv.Product,
			Warehouse:         inv.Warehouse,
		}
	}

	return &dto.InventoryListResponse{
		PaginatedResponse: dto.PaginatedResponse[dto.InventoryResponse]{
			Data:  responses,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	}, nil
}

// AddStock increases the physical quantity of an item using optimistic locking and logs an IN transaction.
func (s *inventoryService) AddStock(ctx context.Context, id uuid.UUID, quantity int) error {
	inv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	err = s.repo.UpdateStock(ctx, id, quantity, inv.Version)
	if err != nil {
		return err
	}

	// Record transaction
	_ = s.logTransaction(ctx, inv, quantity, "IN", "Stock added via AddStock API")

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return nil
}

// RemoveStock decreases the physical quantity of an item and logs an OUT transaction.
// Returns an error if available quantity is insufficient.
func (s *inventoryService) RemoveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	inv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if inv.AvailableQuantity() < quantity {
		return fmt.Errorf("insufficient stock: available %d, required %d", inv.AvailableQuantity(), quantity)
	}

	err = s.repo.UpdateStock(ctx, id, -quantity, inv.Version)
	if err != nil {
		return err
	}

	// Record transaction
	_ = s.logTransaction(ctx, inv, -quantity, "OUT", "Stock removed via RemoveStock API")

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return nil
}

// ReserveStock earmarks a quantity for a future transaction without decreasing physical stock.
// Logs a RESERVE transaction (quantity change 0 as physical balance doesn't change).
func (s *inventoryService) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	inv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if inv.AvailableQuantity() < quantity {
		return fmt.Errorf("insufficient stock for reservation: available %d, required %d", inv.AvailableQuantity(), quantity)
	}

	err = s.repo.UpdateReservedStock(ctx, id, quantity, inv.Version)
	if err != nil {
		return err
	}

	// Reservations don't change physical quantity balance immediately,
	// but we log it as a non-quantity-changing transaction for audit.
	_ = s.logTransaction(ctx, inv, 0, "RESERVE", fmt.Sprintf("Reserved %d units", quantity))

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return nil
}

// ReleaseStock removes a reservation, making the quantity available again.
// Logs a RELEASE transaction.
func (s *inventoryService) ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	inv, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if inv.ReservedQuantity < quantity {
		return fmt.Errorf("cannot release more than reserved: reserved %d, required %d", inv.ReservedQuantity, quantity)
	}

	err = s.repo.UpdateReservedStock(ctx, id, -quantity, inv.Version)
	if err != nil {
		return err
	}

	_ = s.logTransaction(ctx, inv, 0, "RELEASE", fmt.Sprintf("Released %d reserved units", quantity))

	_ = s.cache.DeletePrefix(ctx, _inventoryCachePrefix)
	return nil
}

// GetTransactionHistory retrieves a paginated audit log of stock movements for the specific filters.
func (s *inventoryService) GetTransactionHistory(ctx context.Context, req dto.TransactionListRequest) (dto.TransactionListResponse, error) {
	filter := repository.TransactionListFilter{
		InventoryID:     req.InventoryID,
		ProductID:       req.ProductID,
		WarehouseID:     req.WarehouseID,
		UserID:          req.UserID,
		TransactionType: req.TransactionType,
		StartDate:       req.StartDate,
		EndDate:         req.EndDate,
		Pagination: repository.PaginationParams{
			Page:  req.GetPage(),
			Limit: req.GetLimit(),
			Sort:  req.GetSort("created_at"),
			Order: req.GetOrder(),
		},
	}

	transactions, total, err := s.transactionRepo.FindAll(ctx, filter)
	if err != nil {
		return dto.TransactionListResponse{}, err
	}

	responses := make([]dto.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = dto.TransactionResponse{
			ID:              tx.ID,
			InventoryID:     tx.InventoryID,
			ProductID:       tx.ProductID,
			WarehouseID:     tx.WarehouseID,
			UserID:          tx.UserID,
			QuantityChange:  tx.QuantityChange,
			QuantityBalance: tx.QuantityBalance,
			TransactionType: tx.TransactionType,
			ReferenceID:     tx.ReferenceID,
			Reason:          tx.Reason,
			CreatedAt:       tx.CreatedAt,
			Product:         tx.Product,
			Warehouse:       tx.Warehouse,
			User:            tx.User,
		}
	}

	return dto.TransactionListResponse{
		PaginatedResponse: dto.NewPaginatedResponse(responses, total, filter.Pagination.Page, filter.Pagination.Limit),
	}, nil
}

// Helpers

// getUserIDFromCtx extracts the userID uuid from the context if present.
func (s *inventoryService) getUserIDFromCtx(ctx context.Context) *uuid.UUID {
	val := ctx.Value("userID")
	if val == nil {
		return nil
	}
	id, ok := val.(uuid.UUID)
	if !ok {
		return nil
	}
	return &id
}

// logTransaction is an internal helper to persist a new InventoryTransaction record.
func (s *inventoryService) logTransaction(ctx context.Context, inv *models.Inventory, delta int, txType string, reason string) error {
	userID := s.getUserIDFromCtx(ctx)

	// Balance after change
	balance := inv.Quantity + delta

	transaction := &models.InventoryTransaction{
		InventoryID:     inv.ID,
		ProductID:       inv.ProductID,
		WarehouseID:     inv.WarehouseID,
		UserID:          userID,
		QuantityChange:  delta,
		QuantityBalance: balance,
		TransactionType: txType,
		Reason:          reason,
	}

	return s.transactionRepo.Create(ctx, nil, transaction)
}
