package repository

import (
	"context"

	"github.com/i56/framework/db"
	"github.com/i56/modules/transport/domain"
)

// PgRouteRepo implements RouteRepository backed by PostgreSQL.
type PgRouteRepo struct{}

func NewPgRouteRepo() *PgRouteRepo { return &PgRouteRepo{} }

func (r *PgRouteRepo) Create(ctx context.Context, route *domain.Route) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO routes (tenant_id, name, warehouse_id, transport_mode, base_weight_price, status)
		 VALUES ($1,$2,$3,$4,$5,$6) RETURNING id, created_at, updated_at`,
		route.TenantID, route.Name, route.WarehouseID, route.TransportType, route.BaseWeightPrice, "active",
	).Scan(&route.ID, &route.CreatedAt, &route.UpdatedAt)
}

func (r *PgRouteRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Route, error) {
	rt := &domain.Route{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, warehouse_id, transport_mode, base_weight_price, status, created_at, updated_at
		 FROM routes WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&rt.ID, &rt.TenantID, &rt.Name, &rt.WarehouseID, &rt.TransportType, &rt.BaseWeightPrice,
		&status, &rt.CreatedAt, &rt.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	rt.IsActive = (status == "active")
	return rt, nil
}

func (r *PgRouteRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Route, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM routes WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, name, warehouse_id, transport_mode, base_weight_price, status, created_at, updated_at
		 FROM routes WHERE tenant_id=$1 ORDER BY id LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Route
	for rows.Next() {
		var rt domain.Route
		var status string
		if err := rows.Scan(&rt.ID, &rt.TenantID, &rt.Name, &rt.WarehouseID, &rt.TransportType, &rt.BaseWeightPrice,
			&status, &rt.CreatedAt, &rt.UpdatedAt); err != nil {
			return nil, 0, err
		}
		rt.IsActive = (status == "active")
		result = append(result, rt)
	}
	return result, total, nil
}

func (r *PgRouteRepo) Update(ctx context.Context, route *domain.Route) error {
	status := "inactive"
	if route.IsActive {
		status = "active"
	}
	_, err := db.Pool.Exec(ctx,
		`UPDATE routes SET name=$1, warehouse_id=$2, transport_mode=$3, base_weight_price=$4, status=$5, updated_at=NOW()
		 WHERE id=$6 AND tenant_id=$7`,
		route.Name, route.WarehouseID, route.TransportType, route.BaseWeightPrice, status,
		route.ID, route.TenantID)
	return err
}

var _ RouteRepository = (*PgRouteRepo)(nil)

// PgCarrierRepo implements CourierRepository backed by PostgreSQL.
type PgCarrierRepo struct{}

func NewPgCarrierRepo() *PgCarrierRepo { return &PgCarrierRepo{} }

func (r *PgCarrierRepo) Create(ctx context.Context, c *domain.Courier) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO carriers (tenant_id, name, code, type, status) VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		c.ID, c.Name, c.Code, "", "active",
	).Scan(&c.ID)
}

func (r *PgCarrierRepo) List(ctx context.Context) ([]domain.Courier, error) {
	rows, err := db.Pool.Query(ctx,
		`SELECT id, name, code, type FROM carriers ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []domain.Courier
	for rows.Next() {
		var c domain.Courier
		var carrierType string
		if err := rows.Scan(&c.ID, &c.Name, &c.Code, &carrierType); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, nil
}

func (r *PgCarrierRepo) DetectByTrackingNo(trackingNo string) *domain.Courier {
	if len(trackingNo) < 2 {
		return nil
	}
	prefix := trackingNo[:2]
	c := &domain.Courier{}
	err := db.Pool.QueryRow(context.Background(),
		`SELECT id, name, code FROM carriers WHERE code LIKE $1 LIMIT 1`, prefix+"%",
	).Scan(&c.ID, &c.Name, &c.Code)
	if err != nil {
		return nil
	}
	return c
}

var _ CourierRepository = (*PgCarrierRepo)(nil)
