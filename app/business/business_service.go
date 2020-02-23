package business

import (
	"context"
	"fmt"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// Service ...
type Service struct {
	store Store
}

// NewService constructor for AuthService
func NewService(store Store) Service {
	return Service{store: store}
}

// CreateBusiness creates business
func (s *Service) CreateBusiness(ctx context.Context, userID string, name string) (*Business, error) {
	const op = "business/service.CreateBusiness"

	existingBusiness, err := s.store.GetBusinessByName(ctx, name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", name))
	}

	business := NewBusiness(userID, name)

	err = s.store.StoreBusiness(ctx, business)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to store business")
	}

	return business, nil
}

// UpdateBusiness updates business
func (s *Service) UpdateBusiness(ctx context.Context, id string, name string, profileImageID string) (*Business, error) {
	const op = "business/service.UpdateBusiness"

	existingBusiness, err := s.store.GetBusinessByName(ctx, name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", name))
	}

	business, err := s.store.GetBusinessByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get business by id")
	}

	if business == nil {
		return nil, errors.NotFound(op)
	}

	business.UpdatedAt = time.Now()
	business.Name = name
	business.ProfileImageID = profileImageID

	err = s.store.UpdateBusiness(ctx, business)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update business")
	}

	return business, nil
}

// DeleteBusiness updates business
func (s *Service) DeleteBusiness(ctx context.Context, id string) (*Business, error) {
	const op = "business/service.DeleteBusiness"

	business, err := s.store.GetBusinessByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get business by id")
	}

	if business == nil {
		return nil, errors.NotFound(op)
	}

	err = s.store.DeleteBusiness(ctx, business)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update business")
	}

	return business, nil
}
