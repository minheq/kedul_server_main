package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

// Business ...
type Business struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Name           string    `json:"name"`
	ProfileImageID string    `json:"profile_image_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

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

// CreateBusinessInput ...
type CreateBusinessInput struct {
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// CreateBusiness creates business
func (s *BusinessService) CreateBusiness(ctx context.Context, userID string, input *CreateBusinessInput) (*Business, error) {
	const op = "app/businessService.CreateBusiness"

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, input.Name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", input.Name))
	}

	now := time.Now()

	business := &Business{
		ID:             uuid.Must(uuid.New(), nil).String(),
		UserID:         userID,
		Name:           input.Name,
		ProfileImageID: input.ProfileImageID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = s.businessStore.StoreBusiness(ctx, business)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to businessStore business")
	}

	return business, nil
}

// UpdateBusinessInput ...
type UpdateBusinessInput struct {
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// UpdateBusiness updates business
func (s *BusinessService) UpdateBusiness(ctx context.Context, id string, input *UpdateBusinessInput, currentUser *auth.User) (*Business, error) {
	const op = "app/businessService.UpdateBusiness"

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, input.Name)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", input.Name))
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

	if input.Name != "" {
		business.Name = input.Name
	}
	if input.ProfileImageID != "" {
		business.ProfileImageID = input.ProfileImageID
	}

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
