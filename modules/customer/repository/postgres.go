package repository

import (
	"context"

	"github.com/i56/framework/core/errors"
	"github.com/i56/framework/db"
	"github.com/i56/modules/customer/domain"
)

// PgClientRepo implements ClientRepository backed by PostgreSQL.
type PgClientRepo struct{}

func NewPgClientRepo() *PgClientRepo { return &PgClientRepo{} }

func (r *PgClientRepo) Create(ctx context.Context, tenantID int64, c *domain.Client) error {
	return db.Pool.QueryRow(ctx,
		`INSERT INTO clients (tenant_id, name, code, contact, phone, email, status, balance)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id, created_at`,
		tenantID, c.Name, c.Code, c.ContactName, c.ContactPhone, c.ContactEmail, "active", c.Balance,
	).Scan(&c.ID, &c.CreatedAt)
}

func (r *PgClientRepo) GetByID(ctx context.Context, tenantID, id int64) (*domain.Client, error) {
	c := &domain.Client{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, code, contact, phone, email, status, balance, created_at
		 FROM clients WHERE id=$1 AND tenant_id=$2`, id, tenantID,
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.Code, &c.ContactName, &c.ContactPhone, &c.ContactEmail,
		&status, &c.Balance, &c.CreatedAt)
	if err != nil {
		return nil, nil
	}
	c.IsActive = (status == "active")
	return c, nil
}

func (r *PgClientRepo) GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Client, error) {
	c := &domain.Client{}
	var status string
	err := db.Pool.QueryRow(ctx,
		`SELECT id, tenant_id, name, code, contact, phone, email, status, balance, created_at
		 FROM clients WHERE tenant_id=$1 AND code=$2`, tenantID, code,
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.Code, &c.ContactName, &c.ContactPhone, &c.ContactEmail,
		&status, &c.Balance, &c.CreatedAt)
	if err != nil {
		return nil, nil
	}
	c.IsActive = (status == "active")
	return c, nil
}

func (r *PgClientRepo) List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Client, int64, error) {
	var total int64
	db.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM clients WHERE tenant_id=$1`, tenantID).Scan(&total)
	rows, err := db.Pool.Query(ctx,
		`SELECT id, tenant_id, name, code, contact, phone, email, status, balance, created_at
		 FROM clients WHERE tenant_id=$1 ORDER BY id LIMIT $2 OFFSET $3`, tenantID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var result []domain.Client
	for rows.Next() {
		var c domain.Client
		var status string
		if err := rows.Scan(&c.ID, &c.TenantID, &c.Name, &c.Code, &c.ContactName, &c.ContactPhone, &c.ContactEmail,
			&status, &c.Balance, &c.CreatedAt); err != nil {
			return nil, 0, err
		}
		c.IsActive = (status == "active")
		result = append(result, c)
	}
	return result, total, nil
}

func (r *PgClientRepo) Update(ctx context.Context, tenantID, id int64, c *domain.Client) error {
	status := "inactive"
	if c.IsActive {
		status = "active"
	}
	tag, err := db.Pool.Exec(ctx,
		`UPDATE clients SET name=$1, code=$2, contact=$3, phone=$4, email=$5, status=$6, balance=$7
		 WHERE id=$8 AND tenant_id=$9`,
		c.Name, c.Code, c.ContactName, c.ContactPhone, c.ContactEmail, status, c.Balance, id, tenantID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errors.NewNotFound("Client")
	}
	return nil
}

func (r *PgClientRepo) Delete(ctx context.Context, tenantID, id int64) error {
	_, err := db.Pool.Exec(ctx, `DELETE FROM clients WHERE id=$1 AND tenant_id=$2`, id, tenantID)
	return err
}

var _ ClientRepository = (*PgClientRepo)(nil)
