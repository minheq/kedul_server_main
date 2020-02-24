package app

import (
	"context"
	"fmt"
	"time"

	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

// BusinessService ...
type BusinessService struct {
	businessStore BusinessStore
}

// NewBusinessService constructor for AuthService
func NewBusinessService(businessStore BusinessStore) BusinessService {
	return BusinessService{businessStore: businessStore}
}

// GetBusinessByID ...
func (s *BusinessService) GetBusinessByID(ctx context.Context, id string) (*Business, error) {
	const op = "app/businessService.GetBusinessByID"

	business, err := s.businessStore.GetBusinessByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get business by id")
	}

	return business, nil
}

// CreateBusiness creates business
func (s *BusinessService) CreateBusiness(ctx context.Context, userID string, name string) (*Business, error) {
	const op = "app/businessService.CreateBusiness"

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", name))
	}

	business := NewBusiness(userID, name)

	err = s.businessStore.StoreBusiness(ctx, business)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to businessStore business")
	}

	return business, nil
}

// UpdateBusiness updates business
func (s *BusinessService) UpdateBusiness(ctx context.Context, id string, name string, profileImageID string, currentUser *auth.User) (*Business, error) {
	const op = "app/businessService.UpdateBusiness"

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", name))
	}

	business, err := s.businessStore.GetBusinessByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get business by id")
	}

	if business == nil {
		return nil, errors.NotFound(op)
	}

	if business.UserID != currentUser.ID {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	business.UpdatedAt = time.Now()
	business.Name = name
	business.ProfileImageID = profileImageID

	err = s.businessStore.UpdateBusiness(ctx, business)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update business")
	}

	return business, nil
}

// DeleteBusiness updates business
func (s *BusinessService) DeleteBusiness(ctx context.Context, id string, currentUser *auth.User) (*Business, error) {
	const op = "app/businessService.DeleteBusiness"

	business, err := s.businessStore.GetBusinessByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get business by id")
	}

	if business == nil {
		return nil, errors.NotFound(op)
	}

	if business.UserID != currentUser.ID {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	err = s.businessStore.DeleteBusiness(ctx, business)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update business")
	}

	return business, nil
}
