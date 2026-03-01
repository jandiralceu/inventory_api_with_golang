package repository

// PaginationParams defines the common pagination parameters for all repositories.
type PaginationParams struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

// GetOffset calculates the offset for database queries.
func (p PaginationParams) GetOffset() int {
	if p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit
}

// GetOrderBy returns a string suitable for GORM's Order method.
func (p PaginationParams) GetOrderBy() string {
	if p.Sort == "" {
		return "created_at DESC"
	}
	return p.Sort + " " + p.Order
}
