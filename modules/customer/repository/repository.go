// Package repository defines persistence interfaces for the customer module.
package repository

import (
	"context"
	"github.com/i56/modules/customer/domain"
)

// ClientRepository defines client persistence operations.
type ClientRepository interface {
	Create(ctx context.Context, tenantID int64, client *domain.Client) error
	GetByID(ctx context.Context, tenantID, id int64) (*domain.Client, error)
	GetByCode(ctx context.Context, tenantID int64, code string) (*domain.Client, error)
	List(ctx context.Context, tenantID int64, offset, limit int) ([]domain.Client, int64, error)
	Update(ctx context.Context, tenantID, id int64, client *domain.Client) error
	Delete(ctx context.Context, tenantID, id int64) error
}

// DeclarantRepository defines declarant persistence operations.
type DeclarantRepository interface {
	Create(ctx context.Context, clientID int64, d *domain.Declarant) error
	GetByID(ctx context.Context, clientID, id int64) (*domain.Declarant, error)
	List(ctx context.Context, clientID int64, offset, limit int) ([]domain.Declarant, int64, error)
	Update(ctx context.Context, clientID, id int64, d *domain.Declarant) error
}

// MemberRepository defines member persistence operations.
type MemberRepository interface {
	Create(ctx context.Context, clientID int64, m *domain.ClientMember) error
	GetByID(ctx context.Context, clientID, id int64) (*domain.ClientMember, error)
	List(ctx context.Context, clientID int64, offset, limit int) ([]domain.ClientMember, int64, error)
}

// AddressRepository defines address persistence operations.
type AddressRepository interface {
	Create(ctx context.Context, memberID int64, a *domain.MemberAddress) error
	List(ctx context.Context, memberID int64) ([]domain.MemberAddress, error)
	SetDefault(ctx context.Context, memberID, addressID int64) error
}

// ClientUserRepository defines client user persistence operations.
type ClientUserRepository interface {
	Create(ctx context.Context, clientID int64, u *domain.ClientUser) error
	GetByID(ctx context.Context, id int64) (*domain.ClientUser, error)
	GetByUsername(ctx context.Context, clientID int64, username string) (*domain.ClientUser, error)
	ListByClient(ctx context.Context, clientID int64) ([]domain.ClientUser, error)
	Update(ctx context.Context, id int64, u *domain.ClientUser) error
	Delete(ctx context.Context, id int64) error
	RecordLogin(ctx context.Context, id int64) error
}
