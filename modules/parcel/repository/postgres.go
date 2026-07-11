package repository

import (
	"context"

	"github.com/i56/framework/db"
	"github.com/i56/modules/parcel/domain"
)

// PgParcelRepo implements ParcelRepository backed by PostgreSQL.
type PgParcelRepo struct{}

func NewPgParcelRepo() *PgParcelRepo { return &PgParcelRepo{} }

func (r *PgParcelRepo) Create(ctx context.Context, p *domain.Parcel) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO parcels (tenant_id, warehouse_id, client_id, tracking_number, courier_code, cargo_type, product_name, length, width, height, actual_weight, status, remark)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id, created_at, updated_at`,
		p.TenantID, p.WarehouseID, p.ClientID, p.TrackingNumber, p.CourierCode, p.CargoType,
		p.ProductName, p.Length, p.Width, p.Height, p.ActualWeight, string(p.Status), "",
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PgParcelRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Parcel, error) {
	p := &domain.Parcel{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, warehouse_id, client_id, tracking_number, courier_code, cargo_type, product_name, length, width, height, actual_weight, status, created_at, updated_at
		 FROM parcels WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&p.ID, &p.TenantID, &p.WarehouseID, &p.ClientID, &p.TrackingNumber, &p.CourierCode,
		&p.CargoType, &p.ProductName, &p.Length, &p.Width, &p.Height, &p.ActualWeight, &status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	p.Status = domain.ParcelStatus(status)
	return p, nil
}

func (r *PgParcelRepo) GetByTrackingNo(ctx context.Context, tenantID int64, trackingNo string) (*domain.Parcel, error) {
	p := &domain.Parcel{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, warehouse_id, client_id, tracking_number, courier_code, cargo_type, product_name, length, width, height, actual_weight, status, created_at, updated_at
		 FROM parcels WHERE tenant_id=$1 AND tracking_number=$2`, tenantID, trackingNo,
	).Scan(&p.ID, &p.TenantID, &p.WarehouseID, &p.ClientID, &p.TrackingNumber, &p.CourierCode,
		&p.CargoType, &p.ProductName, &p.Length, &p.Width, &p.Height, &p.ActualWeight, &status, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	p.Status = domain.ParcelStatus(status)
	return p, nil
}

func (r *PgParcelRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM parcels WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, warehouse_id, client_id, tracking_number, courier_code, cargo_type, product_name, length, width, height, actual_weight, status, created_at, updated_at
		 FROM parcels WHERE tenant_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Parcel
	for rows.Next() {
		var p domain.Parcel
		var status string
		if err := rows.Scan(&p.ID, &p.TenantID, &p.WarehouseID, &p.ClientID, &p.TrackingNumber,
			&p.CourierCode, &p.CargoType, &p.ProductName, &p.Length, &p.Width, &p.Height,
			&p.ActualWeight, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		p.Status = domain.ParcelStatus(status)
		result = append(result, p)
	}
	return result, total, nil
}

func (r *PgParcelRepo) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Parcel, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM parcels WHERE tenant_id=$1 AND client_id=$2`, tenantID, clientID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, warehouse_id, client_id, tracking_number, courier_code, cargo_type, product_name, length, width, height, actual_weight, status, created_at, updated_at
		 FROM parcels WHERE tenant_id=$1 AND client_id=$2 ORDER BY id DESC LIMIT $3 OFFSET $4`, tenantID, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Parcel
	for rows.Next() {
		var p domain.Parcel
		var status string
		if err := rows.Scan(&p.ID, &p.TenantID, &p.WarehouseID, &p.ClientID, &p.TrackingNumber,
			&p.CourierCode, &p.CargoType, &p.ProductName, &p.Length, &p.Width, &p.Height,
			&p.ActualWeight, &status, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		p.Status = domain.ParcelStatus(status)
		result = append(result, p)
	}
	return result, total, nil
}

func (r *PgParcelRepo) Update(ctx context.Context, p *domain.Parcel) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE parcels SET warehouse_id=$1, client_id=$2, tracking_number=$3, courier_code=$4, cargo_type=$5,
		 product_name=$6, length=$7, width=$8, height=$9, actual_weight=$10, status=$11, updated_at=NOW()
		 WHERE id=$12 AND tenant_id=$13`,
		p.WarehouseID, p.ClientID, p.TrackingNumber, p.CourierCode, p.CargoType,
		p.ProductName, p.Length, p.Width, p.Height, p.ActualWeight, string(p.Status),
		p.ID, p.TenantID)
	return err
}

var _ ParcelRepository = (*PgParcelRepo)(nil)
