package inventory

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jmoiron/sqlx"

	"bearuang-go/internal/database"
	"bearuang-go/internal/module/products"
	"bearuang-go/internal/module/warehouses"
)

// TestInventoryAndLifecycleIntegration is a PostgreSQL integration test requiring an already-migrated test database.
func TestInventoryAndLifecycleIntegration(t *testing.T) {
	db := openLifecycleTestDB(t)

	warehouseRepository := warehouses.NewPostgresRepository(db)
	productRepository := products.NewPostgresRepository(db)
	inventoryRepository := NewPostgresRepository(db)
	ctx := context.Background()

	unused := newLifecycleFixture(t, db)
	if err := warehouseRepository.Delete(ctx, unused.warehouseID); err != nil {
		t.Fatalf("delete unused warehouse: %v", err)
	}
	if err := productRepository.DeleteVariant(ctx, unused.productID, unused.variantID); err != nil {
		t.Fatalf("delete unused variant: %v", err)
	}
	if err := productRepository.Delete(ctx, unused.productID); err != nil {
		t.Fatalf("delete product with unused variants: %v", err)
	}

	used := newLifecycleFixture(t, db)
	balance, inbound, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: used.warehouseID,
		VariantID:   used.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    6,
	})
	if err != nil {
		t.Fatalf("adjustment in: %v", err)
	}
	if balance.Quantity != 6 || inbound.BalanceAfter != 6 {
		t.Fatalf("unexpected adjustment-in result: balance=%d movement=%d", balance.Quantity, inbound.BalanceAfter)
	}
	assertLifecycleSnapshot(t, db, used.variantID, lifecycleSnapshot{Stock: 6, Balance: 6, Movements: 1})

	_, err = db.ExecContext(ctx, `UPDATE PRODUCT_VARIANTS SET stock = -1 WHERE id = $1`, used.variantID)
	var pgError *pgconn.PgError
	if !errors.As(err, &pgError) || pgError.ConstraintName != "product_variant_stock_nonnegative_check" {
		t.Fatalf("negative aggregate stock error = %v, want product_variant_stock_nonnegative_check", err)
	}

	_, _, err = inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: used.warehouseID,
		VariantID:   used.variantID,
		Type:        movementTypeAdjustmentOut,
		Quantity:    7,
	})
	if !errors.Is(err, errInsufficientStock) {
		t.Fatalf("insufficient adjustment error = %v, want %v", err, errInsufficientStock)
	}
	assertLifecycleSnapshot(t, db, used.variantID, lifecycleSnapshot{Stock: 6, Balance: 6, Movements: 1})

	if err := warehouseRepository.Delete(ctx, used.warehouseID); err == nil || !strings.Contains(err.Error(), "warehouse has inventory history") {
		t.Fatalf("delete warehouse with stock error = %v, want warehouse_in_use", err)
	}
	if err := productRepository.DeleteVariant(ctx, used.productID, used.variantID); err == nil || !strings.Contains(err.Error(), "product variant has inventory history") {
		t.Fatalf("delete variant with stock error = %v, want variant_in_use", err)
	}
	if err := productRepository.Delete(ctx, used.productID); err == nil || !strings.Contains(err.Error(), "product contains variants with inventory history") {
		t.Fatalf("delete product with stock error = %v, want product_in_use", err)
	}

	balance, outbound, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: used.warehouseID,
		VariantID:   used.variantID,
		Type:        movementTypeAdjustmentOut,
		Quantity:    6,
	})
	if err != nil {
		t.Fatalf("adjustment out: %v", err)
	}
	if balance.Quantity != 0 || outbound.BalanceAfter != 0 {
		t.Fatalf("unexpected adjustment-out result: balance=%d movement=%d", balance.Quantity, outbound.BalanceAfter)
	}
	assertLifecycleSnapshot(t, db, used.variantID, lifecycleSnapshot{Stock: 0, Balance: 0, Movements: 2})

	if err := warehouseRepository.Delete(ctx, used.warehouseID); err == nil || !strings.Contains(err.Error(), "warehouse has inventory history") {
		t.Fatalf("delete warehouse with zero balance and history error = %v, want warehouse_in_use", err)
	}
	if err := productRepository.DeleteVariant(ctx, used.productID, used.variantID); err == nil || !strings.Contains(err.Error(), "product variant has inventory history") {
		t.Fatalf("delete variant with zero stock and history error = %v, want variant_in_use", err)
	}
	if err := productRepository.Delete(ctx, used.productID); err == nil || !strings.Contains(err.Error(), "product contains variants with inventory history") {
		t.Fatalf("delete product with zero stock and history error = %v, want product_in_use", err)
	}

	warehouse, err := warehouseRepository.GetByID(ctx, used.warehouseID)
	if err != nil {
		t.Fatalf("get warehouse for archive: %v", err)
	}
	warehouse.Status = "archived"
	if err := warehouseRepository.Update(ctx, warehouse); err != nil {
		t.Fatalf("archive warehouse: %v", err)
	}
	if _, err := inventoryRepository.GetBalance(ctx, used.warehouseID, used.variantID); err != nil {
		t.Fatalf("read archived warehouse balance: %v", err)
	}
	movements, err := inventoryRepository.GetMovements(ctx, &used.warehouseID, &used.variantID)
	if err != nil || len(movements) != 2 {
		t.Fatalf("read archived warehouse movements: count=%d err=%v", len(movements), err)
	}
	if _, err := inventoryRepository.GetMovementByID(ctx, inbound.ID); err != nil {
		t.Fatalf("read archived warehouse movement by ID: %v", err)
	}
	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: used.warehouseID,
		VariantID:   used.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    1,
	}); !errors.Is(err, errWarehouseNotActive) {
		t.Fatalf("adjust archived warehouse error = %v, want %v", err, errWarehouseNotActive)
	}
}

// TestWarehouseStatusLifecycleIntegration verifies warehouse status rules against an already-migrated PostgreSQL test database.
func TestWarehouseStatusLifecycleIntegration(t *testing.T) {
	db := openLifecycleTestDB(t)
	warehouseRepository := warehouses.NewPostgresRepository(db)
	inventoryRepository := NewPostgresRepository(db)
	ctx := context.Background()

	withStock := newLifecycleFixture(t, db)
	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: withStock.warehouseID,
		VariantID:   withStock.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    5,
	}); err != nil {
		t.Fatalf("adjustment in: %v", err)
	}
	warehouse, err := warehouseRepository.GetByID(ctx, withStock.warehouseID)
	if err != nil {
		t.Fatalf("get warehouse for active update: %v", err)
	}
	warehouse.Name = "Renamed warehouse"
	if err := warehouseRepository.Update(ctx, warehouse); err != nil {
		t.Fatalf("update active warehouse with stock: %v", err)
	}
	for _, status := range []string{"inactive", "archived"} {
		if err := setWarehouseLifecycleStatus(warehouseRepository, ctx, withStock.warehouseID, status); err == nil || !strings.Contains(err.Error(), "warehouse must have zero stock") {
			t.Fatalf("set warehouse status to %s with stock error = %v, want warehouse_has_stock", status, err)
		}
		warehouse, err := warehouseRepository.GetByID(ctx, withStock.warehouseID)
		if err != nil || warehouse.Status != "active" {
			t.Fatalf("warehouse status after rejected %s transition = %+v, err=%v; want active", status, warehouse, err)
		}
		assertLifecycleSnapshot(t, db, withStock.variantID, lifecycleSnapshot{Stock: 5, Balance: 5, Movements: 1})
	}

	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: withStock.warehouseID,
		VariantID:   withStock.variantID,
		Type:        movementTypeAdjustmentOut,
		Quantity:    5,
	}); err != nil {
		t.Fatalf("adjustment out: %v", err)
	}
	if err := setWarehouseLifecycleStatus(warehouseRepository, ctx, withStock.warehouseID, "archived"); err != nil {
		t.Fatalf("archive zero-stock warehouse with history: %v", err)
	}
	movements, err := inventoryRepository.GetMovements(ctx, &withStock.warehouseID, &withStock.variantID)
	if err != nil || len(movements) != 2 {
		t.Fatalf("read archived warehouse history: count=%d err=%v", len(movements), err)
	}
	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: withStock.warehouseID,
		VariantID:   withStock.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    1,
	}); !errors.Is(err, errWarehouseNotActive) {
		t.Fatalf("adjust archived warehouse error = %v, want %v", err, errWarehouseNotActive)
	}

	inactive := newLifecycleFixture(t, db)
	if err := setWarehouseLifecycleStatus(warehouseRepository, ctx, inactive.warehouseID, "inactive"); err != nil {
		t.Fatalf("inactivate zero-stock warehouse: %v", err)
	}
	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: inactive.warehouseID,
		VariantID:   inactive.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    1,
	}); !errors.Is(err, errWarehouseNotActive) {
		t.Fatalf("adjust inactive warehouse error = %v, want %v", err, errWarehouseNotActive)
	}
	if err := setWarehouseLifecycleStatus(warehouseRepository, ctx, inactive.warehouseID, "active"); err != nil {
		t.Fatalf("reactivate zero-stock warehouse: %v", err)
	}
	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: inactive.warehouseID,
		VariantID:   inactive.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    1,
	}); err != nil {
		t.Fatalf("adjust reactivated warehouse: %v", err)
	}
}

// TestWarehouseStatusAdjustmentRace verifies PostgreSQL row locking prevents stock creation during warehouse retirement.
func TestWarehouseStatusAdjustmentRace(t *testing.T) {
	db := openLifecycleTestDB(t)
	fixture := newLifecycleFixture(t, db)
	inventoryRepository := NewPostgresRepository(db)
	warehouseRepository := warehouses.NewPostgresRepository(db)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	start := make(chan struct{})
	adjustmentResult := make(chan error, 1)
	statusResult := make(chan error, 1)
	go func() {
		<-start
		_, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
			WarehouseID: fixture.warehouseID,
			VariantID:   fixture.variantID,
			Type:        movementTypeAdjustmentIn,
			Quantity:    3,
		})
		adjustmentResult <- err
	}()
	go func() {
		<-start
		statusResult <- setWarehouseLifecycleStatus(warehouseRepository, ctx, fixture.warehouseID, "archived")
	}()
	close(start)

	adjustmentErr := <-adjustmentResult
	statusErr := <-statusResult
	warehouse, err := warehouseRepository.GetByID(ctx, fixture.warehouseID)
	if err != nil {
		t.Fatalf("get warehouse after status race: %v", err)
	}
	switch {
	case adjustmentErr == nil:
		if statusErr == nil || !strings.Contains(statusErr.Error(), "warehouse must have zero stock") {
			t.Fatalf("archive after adjustment error = %v, want warehouse_has_stock", statusErr)
		}
		if warehouse.Status != "active" {
			t.Fatalf("warehouse status after adjustment won race = %q, want active", warehouse.Status)
		}
		assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{Stock: 3, Balance: 3, Movements: 1})
	case errors.Is(adjustmentErr, errWarehouseNotActive):
		if statusErr != nil {
			t.Fatalf("archive after status transition won race: %v", statusErr)
		}
		if warehouse.Status != "archived" {
			t.Fatalf("warehouse status after retirement won race = %q, want archived", warehouse.Status)
		}
		assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{})
	default:
		t.Fatalf("unexpected adjustment result in warehouse status race: %v", adjustmentErr)
	}
}

// TestDeletedProductVariantLifecycleIntegration checks variant APIs against an already-migrated PostgreSQL test database.
func TestDeletedProductVariantLifecycleIntegration(t *testing.T) {
	db := openLifecycleTestDB(t)
	fixture := newLifecycleFixture(t, db)
	productRepository := products.NewPostgresRepository(db)
	inventoryRepository := NewPostgresRepository(db)
	ctx := context.Background()

	if err := productRepository.Delete(ctx, fixture.productID); err != nil {
		t.Fatalf("soft-delete unused product: %v", err)
	}
	if _, err := productRepository.GetVariantByID(ctx, fixture.productID, fixture.variantID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("get variant through deleted product error = %v, want not found", err)
	}
	if _, err := productRepository.GetManyVariantsByProduct(ctx, fixture.productID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("list variants through deleted product error = %v, want not found", err)
	}

	updated := products.ProductVariant{
		ID:        fixture.variantID,
		ProductID: fixture.productID,
		SKU:       "changed-sku",
		Name:      "Changed name",
		Price:     2,
		Unit:      "each",
		Status:    "active",
	}
	if err := productRepository.UpdateVariant(ctx, &updated); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("update variant through deleted product error = %v, want not found", err)
	}
	var sku, name string
	if err := db.GetContext(ctx, &sku, `SELECT sku FROM PRODUCT_VARIANTS WHERE id = $1`, fixture.variantID); err != nil || sku != "test-"+fixture.variantID {
		t.Fatalf("variant SKU after rejected update = %q, err=%v", sku, err)
	}
	if err := db.GetContext(ctx, &name, `SELECT name FROM PRODUCT_VARIANTS WHERE id = $1`, fixture.variantID); err != nil || name != "Integration variant" {
		t.Fatalf("variant name after rejected update = %q, err=%v", name, err)
	}
	if err := productRepository.DeleteVariant(ctx, fixture.productID, fixture.variantID); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("delete variant through deleted product error = %v, want not found", err)
	}

	newVariant := products.ProductVariant{
		ProductID: fixture.productID,
		SKU:       "unavailable-parent-" + fixture.variantID,
		Name:      "Must not be created",
		Price:     1,
		Unit:      "each",
		Status:    "active",
	}
	if err := productRepository.CreateVariant(ctx, &newVariant); !errors.Is(err, sql.ErrNoRows) {
		t.Fatalf("create variant under deleted product error = %v, want not found", err)
	}
	var created int
	if err := db.GetContext(ctx, &created, `SELECT COUNT(*) FROM PRODUCT_VARIANTS WHERE id = $1`, newVariant.ID); err != nil || created != 0 {
		t.Fatalf("variant rows created under deleted product = %d, err=%v", created, err)
	}

	if _, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
		WarehouseID: fixture.warehouseID,
		VariantID:   fixture.variantID,
		Type:        movementTypeAdjustmentIn,
		Quantity:    1,
	}); !errors.Is(err, errVariantNotFound) {
		t.Fatalf("adjust variant under deleted product error = %v, want %v", err, errVariantNotFound)
	}
}

func setWarehouseLifecycleStatus(
	repository interface {
		GetByID(context.Context, string) (*warehouses.Warehouse, error)
		Update(context.Context, *warehouses.Warehouse) error
	},
	ctx context.Context,
	id, status string,
) error {
	warehouse, err := repository.GetByID(ctx, id)
	if err != nil {
		return err
	}
	warehouse.Status = status
	return repository.Update(ctx, warehouse)
}

// TestInventoryAdjustmentLifecycleRaces exercises PostgreSQL row-lock ordering against an already-migrated test database.
func TestInventoryAdjustmentLifecycleRaces(t *testing.T) {
	db := openLifecycleTestDB(t)

	t.Run("warehouse", func(t *testing.T) {
		fixture := newLifecycleFixture(t, db)
		inventoryRepository := NewPostgresRepository(db)
		warehouseRepository := warehouses.NewPostgresRepository(db)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		start := make(chan struct{})
		adjustmentResult := make(chan error, 1)
		deletionResult := make(chan error, 1)
		go func() {
			<-start
			_, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
				WarehouseID: fixture.warehouseID,
				VariantID:   fixture.variantID,
				Type:        movementTypeAdjustmentIn,
				Quantity:    3,
			})
			adjustmentResult <- err
		}()
		go func() {
			<-start
			deletionResult <- warehouseRepository.Delete(ctx, fixture.warehouseID)
		}()
		close(start)

		adjustErr := <-adjustmentResult
		deleteErr := <-deletionResult
		switch {
		case adjustErr == nil:
			if deleteErr == nil || !strings.Contains(deleteErr.Error(), "warehouse has inventory history") {
				t.Fatalf("delete after adjustment error = %v, want warehouse_in_use", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{Stock: 3, Balance: 3, Movements: 1})
		case errors.Is(adjustErr, errWarehouseNotFound):
			if deleteErr != nil {
				t.Fatalf("delete after warehouse won race: %v", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{})
		default:
			t.Fatalf("unexpected adjustment result in warehouse race: %v", adjustErr)
		}
	})

	t.Run("variant", func(t *testing.T) {
		fixture := newLifecycleFixture(t, db)
		inventoryRepository := NewPostgresRepository(db)
		productRepository := products.NewPostgresRepository(db)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		start := make(chan struct{})
		adjustmentResult := make(chan error, 1)
		deletionResult := make(chan error, 1)
		go func() {
			<-start
			_, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
				WarehouseID: fixture.warehouseID,
				VariantID:   fixture.variantID,
				Type:        movementTypeAdjustmentIn,
				Quantity:    3,
			})
			adjustmentResult <- err
		}()
		go func() {
			<-start
			deletionResult <- productRepository.DeleteVariant(ctx, fixture.productID, fixture.variantID)
		}()
		close(start)

		adjustErr := <-adjustmentResult
		deleteErr := <-deletionResult
		switch {
		case adjustErr == nil:
			if deleteErr == nil || !strings.Contains(deleteErr.Error(), "product variant has inventory history") {
				t.Fatalf("delete after adjustment error = %v, want variant_in_use", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{Stock: 3, Balance: 3, Movements: 1})
		case errors.Is(adjustErr, errVariantNotFound):
			if deleteErr != nil {
				t.Fatalf("delete after variant won race: %v", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{})
		default:
			t.Fatalf("unexpected adjustment result in variant race: %v", adjustErr)
		}
	})

	t.Run("product", func(t *testing.T) {
		fixture := newLifecycleFixture(t, db)
		inventoryRepository := NewPostgresRepository(db)
		productRepository := products.NewPostgresRepository(db)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		start := make(chan struct{})
		adjustmentResult := make(chan error, 1)
		deletionResult := make(chan error, 1)
		go func() {
			<-start
			_, _, err := inventoryRepository.Adjust(ctx, &Adjustment{
				WarehouseID: fixture.warehouseID,
				VariantID:   fixture.variantID,
				Type:        movementTypeAdjustmentIn,
				Quantity:    3,
			})
			adjustmentResult <- err
		}()
		go func() {
			<-start
			deletionResult <- productRepository.Delete(ctx, fixture.productID)
		}()
		close(start)

		adjustErr := <-adjustmentResult
		deleteErr := <-deletionResult
		switch {
		case adjustErr == nil:
			if deleteErr == nil || !strings.Contains(deleteErr.Error(), "product contains variants with inventory history") {
				t.Fatalf("delete after adjustment error = %v, want product_in_use", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{Stock: 3, Balance: 3, Movements: 1})
		case errors.Is(adjustErr, errVariantNotFound):
			if deleteErr != nil {
				t.Fatalf("delete after product won race: %v", deleteErr)
			}
			assertLifecycleSnapshot(t, db, fixture.variantID, lifecycleSnapshot{})
		default:
			t.Fatalf("unexpected adjustment result in product race: %v", adjustErr)
		}
	})
}

func openLifecycleTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("PostgreSQL integration test skipped: set TEST_DATABASE_URL to an already-migrated test database; this test does not run migrations")
	}

	db, err := database.NewDB(context.Background(), databaseURL)
	if err != nil {
		t.Fatalf("connect to integration test database: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Errorf("close integration test database: %v", err)
		}
	})
	return db
}

type lifecycleFixture struct {
	warehouseID string
	productID   string
	variantID   string
}

func newLifecycleFixture(t *testing.T, db *sqlx.DB) lifecycleFixture {
	t.Helper()

	fixture := lifecycleFixture{
		warehouseID: database.GenerateID(),
		productID:   database.GenerateID(),
		variantID:   database.GenerateID(),
	}
	t.Cleanup(func() {
		cleanupQueries := []struct {
			query string
			id    string
		}{
			{`DELETE FROM STOCK_MOVEMENTS WHERE warehouse_id = $1`, fixture.warehouseID},
			{`DELETE FROM INVENTORY_BALANCES WHERE warehouse_id = $1`, fixture.warehouseID},
			{`DELETE FROM PRODUCT_VARIANTS WHERE id = $1`, fixture.variantID},
			{`DELETE FROM PRODUCTS WHERE id = $1`, fixture.productID},
			{`DELETE FROM WAREHOUSES WHERE id = $1`, fixture.warehouseID},
		}
		for _, cleanup := range cleanupQueries {
			if _, err := db.ExecContext(context.Background(), cleanup.query, cleanup.id); err != nil {
				t.Errorf("clean up integration fixture: %v", err)
			}
		}
	})

	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `
		INSERT INTO WAREHOUSES (id, code, name, status)
		VALUES ($1, $2, $3, 'active')`,
		fixture.warehouseID,
		"test-"+fixture.warehouseID,
		"Integration warehouse",
	); err != nil {
		t.Fatalf("insert test warehouse: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO PRODUCTS (id, name, slug, description, status)
		VALUES ($1, $2, $3, '', 'active')`,
		fixture.productID,
		"Integration product",
		"test-"+fixture.productID,
	); err != nil {
		t.Fatalf("insert test product: %v", err)
	}
	if _, err := db.ExecContext(ctx, `
		INSERT INTO PRODUCT_VARIANTS (id, product_id, sku, name, price, status)
		VALUES ($1, $2, $3, 'Integration variant', 1, 'active')`,
		fixture.variantID,
		fixture.productID,
		"test-"+fixture.variantID,
	); err != nil {
		t.Fatalf("insert test variant: %v", err)
	}

	return fixture
}

type lifecycleSnapshot struct {
	Stock     int64 `db:"stock"`
	Balance   int64 `db:"balance"`
	Movements int64 `db:"movements"`
}

func assertLifecycleSnapshot(t *testing.T, db *sqlx.DB, variantID string, want lifecycleSnapshot) {
	t.Helper()

	var got lifecycleSnapshot
	query := `
		SELECT
			variant.stock AS stock,
			COALESCE((
				SELECT SUM(quantity)
				FROM INVENTORY_BALANCES
				WHERE variant_id = variant.id
			), 0) AS balance,
			(
				SELECT COUNT(*)
				FROM STOCK_MOVEMENTS
				WHERE variant_id = variant.id
			) AS movements
		FROM PRODUCT_VARIANTS AS variant
		WHERE variant.id = $1`
	if err := db.GetContext(context.Background(), &got, query, variantID); err != nil {
		t.Fatalf("read inventory lifecycle snapshot: %v", err)
	}
	if got != want {
		t.Fatalf("inventory lifecycle snapshot = %+v, want %+v", got, want)
	}
}
