package repository

import (
	"context"

	"github.com/i56/framework/db"
	"github.com/i56/modules/order/domain"
)

// PgOrderRepo implements OrderRepository backed by PostgreSQL.
type PgOrderRepo struct{}

func NewPgOrderRepo() *PgOrderRepo { return &PgOrderRepo{} }

func (r *PgOrderRepo) Create(ctx context.Context, o *domain.Order) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO orders (tenant_id, order_no, client_id, member_id, warehouse_id, route_id,
		 recipient_name, recipient_phone, recipient_address, total_actual_weight, total_chargeable_weight,
		 total_price, status, tracking_numbers, carrier_tracking_no, customs_number, remark)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17)
		 RETURNING id, created_at, updated_at`,
		o.TenantID, o.OrderNo, o.ClientID, o.MemberID, o.WarehouseID, o.RouteID,
		o.RecipientName, "", "", o.TotalActualWeight, o.TotalChargeableWeight,
		o.TotalPrice, string(o.Status), o.TrackingNumbers, o.CarrierTrackingNo, o.CustomsNumber, o.Remark,
	).Scan(&o.ID, &o.CreatedAt, &o.UpdatedAt)
}

func (r *PgOrderRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Order, error) {
	o := &domain.Order{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, order_no, client_id, member_id, warehouse_id, route_id,
		 recipient_name, total_actual_weight, total_chargeable_weight, total_price, status,
		 tracking_numbers, carrier_tracking_no, customs_number, remark, created_at, updated_at
		 FROM orders WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&o.ID, &o.TenantID, &o.OrderNo, &o.ClientID, &o.MemberID, &o.WarehouseID, &o.RouteID,
		&o.RecipientName, &o.TotalActualWeight, &o.TotalChargeableWeight, &o.TotalPrice, &status,
		&o.TrackingNumbers, &o.CarrierTrackingNo, &o.CustomsNumber, &o.Remark, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	o.Status = domain.OrderStatus(status)
	return o, nil
}

func (r *PgOrderRepo) GetByOrderNo(ctx context.Context, tenantID int64, orderNo string) (*domain.Order, error) {
	o := &domain.Order{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, order_no, client_id, member_id, warehouse_id, route_id,
		 recipient_name, total_actual_weight, total_chargeable_weight, total_price, status,
		 tracking_numbers, carrier_tracking_no, customs_number, remark, created_at, updated_at
		 FROM orders WHERE tenant_id=$1 AND order_no=$2`, tenantID, orderNo,
	).Scan(&o.ID, &o.TenantID, &o.OrderNo, &o.ClientID, &o.MemberID, &o.WarehouseID, &o.RouteID,
		&o.RecipientName, &o.TotalActualWeight, &o.TotalChargeableWeight, &o.TotalPrice, &status,
		&o.TrackingNumbers, &o.CarrierTrackingNo, &o.CustomsNumber, &o.Remark, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, nil
	}
	o.Status = domain.OrderStatus(status)
	return o, nil
}

func scanOrder(rows interface{ Scan(...interface{}) error }) (*domain.Order, error) {
	o := &domain.Order{}
	var status string
	err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNo, &o.ClientID, &o.MemberID, &o.WarehouseID, &o.RouteID,
		&o.RecipientName, &o.TotalActualWeight, &o.TotalChargeableWeight, &o.TotalPrice, &status,
		&o.TrackingNumbers, &o.CarrierTrackingNo, &o.CustomsNumber, &o.Remark, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	o.Status = domain.OrderStatus(status)
	return o, nil
}

func (r *PgOrderRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Order, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, order_no, client_id, member_id, warehouse_id, route_id,
		 recipient_name, total_actual_weight, total_chargeable_weight, total_price, status,
		 tracking_numbers, carrier_tracking_no, customs_number, remark, created_at, updated_at
		 FROM orders WHERE tenant_id=$1 ORDER BY id DESC LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Order
	for rows.Next() {
		var o domain.Order
		var status string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNo, &o.ClientID, &o.MemberID, &o.WarehouseID, &o.RouteID,
			&o.RecipientName, &o.TotalActualWeight, &o.TotalChargeableWeight, &o.TotalPrice, &status,
			&o.TrackingNumbers, &o.CarrierTrackingNo, &o.CustomsNumber, &o.Remark, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		o.Status = domain.OrderStatus(status)
		result = append(result, o)
	}
	return result, total, nil
}

func (r *PgOrderRepo) ListByClient(ctx context.Context, tenantID, clientID int64, offset, limit int) ([]domain.Order, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE tenant_id=$1 AND client_id=$2`, tenantID, clientID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, order_no, client_id, member_id, warehouse_id, route_id,
		 recipient_name, total_actual_weight, total_chargeable_weight, total_price, status,
		 tracking_numbers, carrier_tracking_no, customs_number, remark, created_at, updated_at
		 FROM orders WHERE tenant_id=$1 AND client_id=$2 ORDER BY id DESC LIMIT $3 OFFSET $4`, tenantID, clientID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Order
	for rows.Next() {
		var o domain.Order
		var status string
		if err := rows.Scan(&o.ID, &o.TenantID, &o.OrderNo, &o.ClientID, &o.MemberID, &o.WarehouseID, &o.RouteID,
			&o.RecipientName, &o.TotalActualWeight, &o.TotalChargeableWeight, &o.TotalPrice, &status,
			&o.TrackingNumbers, &o.CarrierTrackingNo, &o.CustomsNumber, &o.Remark, &o.CreatedAt, &o.UpdatedAt); err != nil {
			return nil, 0, err
		}
		o.Status = domain.OrderStatus(status)
		result = append(result, o)
	}
	return result, total, nil
}

func (r *PgOrderRepo) Update(ctx context.Context, o *domain.Order) error {
	_, err := db.Pool.Exec(ctx,
		`UPDATE orders SET order_no=$1, client_id=$2, member_id=$3, warehouse_id=$4, route_id=$5,
		 recipient_name=$6, total_actual_weight=$7, total_chargeable_weight=$8, total_price=$9, status=$10,
		 tracking_numbers=$11, carrier_tracking_no=$12, customs_number=$13, remark=$14, updated_at=NOW()
		 WHERE id=$15 AND tenant_id=$16`,
		o.OrderNo, o.ClientID, o.MemberID, o.WarehouseID, o.RouteID,
		o.RecipientName, o.TotalActualWeight, o.TotalChargeableWeight, o.TotalPrice, string(o.Status),
		o.TrackingNumbers, o.CarrierTrackingNo, o.CustomsNumber, o.Remark,
		o.ID, o.TenantID)
	return err
}

var _ OrderRepository = (*PgOrderRepo)(nil)
