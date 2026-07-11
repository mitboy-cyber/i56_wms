package repository

import (
	"context"

	"github.com/i56/framework/db"
	"github.com/i56/modules/warehouse/domain"
	"github.com/i56/framework/core/errors"
)

// PgWarehouseRepo implements WarehouseRepository backed by PostgreSQL.
type PgWarehouseRepo struct{}

func NewPgWarehouseRepo() *PgWarehouseRepo { return &PgWarehouseRepo{} }

func (r *PgWarehouseRepo) Create(ctx context.Context, tenantID int64, w *domain.Warehouse) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO warehouses (tenant_id, name, code, contact, phone, address)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at`,
		tenantID, w.Name, w.Code, w.Contact, w.Phone, w.Address,
	).Scan(&w.ID, &w.CreatedAt)
}

func (r *PgWarehouseRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Warehouse, error) {
	w := &domain.Warehouse{}
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, code, contact, phone, address, created_at
		 FROM warehouses WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&w.ID, &w.TenantID, &w.Name, &w.Code, &w.Contact, &w.Phone, &w.Address, &w.CreatedAt)
	if err != nil {
		return nil, nil
	}
	w.IsActive = true
	return w, nil
}

func (r *PgWarehouseRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Warehouse, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM warehouses WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, name, code, contact, phone, address, created_at
		 FROM warehouses WHERE tenant_id=$1 ORDER BY id LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Warehouse
	for rows.Next() {
		var w domain.Warehouse
		if err := rows.Scan(&w.ID, &w.TenantID, &w.Name, &w.Code, &w.Contact, &w.Phone, &w.Address, &w.CreatedAt); err != nil {
			return nil, 0, err
		}
		w.IsActive = true
		result = append(result, w)
	}
	return result, total, nil
}

func (r *PgWarehouseRepo) Update(ctx context.Context, tenantID, id int64, w *domain.Warehouse) error {
	tag, err := db.Pool.Exec(ctx,
		`UPDATE warehouses SET name=$1, code=$2, contact=$3, phone=$4, address=$5
		 WHERE id=$6 AND tenant_id=$7`,
		w.Name, w.Code, w.Contact, w.Phone, w.Address, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NewNotFound("Warehouse")
	}
	return nil
}

var _ WarehouseRepository = (*PgWarehouseRepo)(nil)
