// Package service implements application services for the customer module.
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/i56/framework/core/errors"
	"github.com/i56/modules/customer/domain"
	"github.com/i56/modules/customer/dto"
	"github.com/i56/modules/customer/repository"
)

// ClientService handles client business logic.
type ClientService struct {
	repo repository.ClientRepository
}

// NewClientService creates a ClientService.
func NewClientService(repo repository.ClientRepository) *ClientService {
	return &ClientService{repo: repo}
}

// Create creates a new client.
func (s *ClientService) Create(ctx context.Context, tenantID int64, input dto.CreateClientRequest) (*dto.ClientResponse, error) {
	// Validate client type
	valid := false
	for _, t := range domain.ValidClientTypes() {
		if t == input.ClientType {
			valid = true
			break
		}
	}
	if !valid {
		return nil, errors.NewValidation("invalid client_type: " + input.ClientType)
	}

	// Generate client code
	code := fmt.Sprintf("C-%d", time.Now().Unix()%100000)

	client := &domain.Client{
		TenantID:     tenantID,
		Name:         input.Name,
		Code:         code,
		ClientType:   domain.ClientType(input.ClientType),
		ContactName:  input.ContactName,
		ContactPhone: input.ContactPhone,
		ContactEmail: input.ContactEmail,
		Balance:      0,
		IsActive:     true,
		Remark:       input.Remark,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.Create(ctx, tenantID, client); err != nil {
		return nil, fmt.Errorf("client create: %w", err)
	}

	return toClientResponse(client), nil
}

// GetByID retrieves a client by ID.
func (s *ClientService) GetByID(ctx context.Context, tenantID, id int64) (*dto.ClientResponse, error) {
	client, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.NewNotFound("Client")
	}
	return toClientResponse(client), nil
}

// List retrieves paginated clients.
func (s *ClientService) List(ctx context.Context, tenantID int64, offset, limit int) ([]dto.ClientResponse, int64, error) {
	clients, total, err := s.repo.List(ctx, tenantID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.ClientResponse, len(clients))
	for i, c := range clients {
		result[i] = *toClientResponse(&c)
	}
	return result, total, nil
}

// Update updates a client.
func (s *ClientService) Update(ctx context.Context, tenantID, id int64, input dto.UpdateClientRequest) (*dto.ClientResponse, error) {
	client, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, errors.NewNotFound("Client")
	}

	if input.Name != nil {
		client.Name = *input.Name
	}
	if input.ContactName != nil {
		client.ContactName = *input.ContactName
	}
	if input.ContactPhone != nil {
		client.ContactPhone = *input.ContactPhone
	}
	if input.ContactEmail != nil {
		client.ContactEmail = *input.ContactEmail
	}
	if input.Remark != nil {
		client.Remark = *input.Remark
	}
	if input.IsActive != nil {
		client.IsActive = *input.IsActive
	}
	client.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, tenantID, id, client); err != nil {
		return nil, fmt.Errorf("client update: %w", err)
	}

	return toClientResponse(client), nil
}

// Delete soft-deletes a client.
func (s *ClientService) Delete(ctx context.Context, tenantID, id int64) error {
	client, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if client == nil {
		return errors.NewNotFound("Client")
	}
	return s.repo.Delete(ctx, tenantID, id)
}

func toClientResponse(c *domain.Client) *dto.ClientResponse {
	return &dto.ClientResponse{
		ID:           c.ID,
		TenantID:     c.TenantID,
		Name:         c.Name,
		Code:         c.Code,
		ClientType:   string(c.ClientType),
		ContactName:  c.ContactName,
		ContactPhone: c.ContactPhone,
		ContactEmail: c.ContactEmail,
		Balance:      c.Balance,
		IsActive:     c.IsActive,
		Remark:       c.Remark,
		CreatedAt:    c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
	}
}
