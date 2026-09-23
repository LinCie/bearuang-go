package warehouses

import "time"

/*
======================================
Warehouses
======================================
*/

const (
	warehouseStatusActive   = "active"
	warehouseStatusInactive = "inactive"
	warehouseStatusArchived = "archived"
)

// Warehouse represents a physical or logical inventory location.
type Warehouse struct {
	ID          string     `db:"id"`
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Address     *string    `db:"address"`
	Status      string     `db:"status"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}
