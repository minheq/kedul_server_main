package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minheq/kedul_server_main/auth"
	"github.com/minheq/kedul_server_main/errors"
)

// Business ...
type Business struct {
	ID             string
	UserID         string
	Name           string
	ProfileImageID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// BusinessService ...
type BusinessService struct {
	businessStore BusinessStore
	locationStore LocationStore
	employeeStore EmployeeStore
}

// NewBusinessService constructor for AuthService
func NewBusinessService(businessStore BusinessStore, locationStore LocationStore, employeeStore EmployeeStore) BusinessService {
	return BusinessService{businessStore: businessStore, locationStore: locationStore, employeeStore: employeeStore}
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

func (s *BusinessService) getBusinessesAsEmployeeByUserID(ctx context.Context, userID string) ([]*Business, error) {
	const op = "app/businessService.getBusinessesAsEmployeeByUserID"

	userAsEmployeeList, err := s.employeeStore.GetEmployeesByUserID(ctx, userID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get employees by user id")
	}

	locationIDs := []string{}

	for _, employee := range userAsEmployeeList {
		locationIDs = append(locationIDs, employee.LocationID)
	}

	locations, err := s.locationStore.GetLocationsByIDs(ctx, locationIDs)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get locations by location ids")
	}

	businessIDs := []string{}

	for _, location := range locations {
		businessIDs = append(businessIDs, location.BusinessID)
	}

	businesses, err := s.businessStore.GetBusinessesByIDs(ctx, businessIDs)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get businesses by businessIDs")
	}

	return businesses, nil
}

// GetBusinessesByUserID ...
func (s *BusinessService) GetBusinessesByUserID(ctx context.Context, userID string, currentUser *auth.User) ([]*Business, error) {
	const op = "app/businessService.GetBusinessesByUserID"

	if userID != currentUser.ID {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	businessesAsEmployee, err := s.getBusinessesAsEmployeeByUserID(ctx, userID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get businesses by based on user as employee")
	}

	businessesAsOwner, err := s.businessStore.GetBusinessesByUserID(ctx, userID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get businesses by based on user as owner")
	}

	businesses := []*Business{}

	for _, businessAsEmployee := range businessesAsEmployee {
		businesses = append(businesses, businessAsEmployee)
	}

	for _, business := range businessesAsOwner {
		alreadyExists := false
		for _, appendedBusiness := range businesses {
			if appendedBusiness.ID == business.ID {
				alreadyExists = true
			}
		}

		if alreadyExists == false {
			businesses = append(businesses, business)
		}
	}

	return businesses, nil
}

// CreateBusinessInput ...
type CreateBusinessInput struct {
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// CreateBusiness creates business
func (s *BusinessService) CreateBusiness(ctx context.Context, userID string, input *CreateBusinessInput) (*Business, error) {
	const op = "app/businessService.CreateBusiness"

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, strings.TrimSpace(input.Name))

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", strings.TrimSpace(input.Name)))
	}

	now := time.Now()

	business := &Business{
		ID:             uuid.Must(uuid.New(), nil).String(),
		UserID:         userID,
		Name:           strings.TrimSpace(input.Name),
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

	existingBusiness, err := s.businessStore.GetBusinessByName(ctx, strings.TrimSpace(input.Name))

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to get business by name")
	}

	if existingBusiness != nil {
		return nil, errors.Invalid(op, fmt.Sprintf("business with name %s already exists", strings.TrimSpace(input.Name)))
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
		business.Name = strings.TrimSpace(input.Name)
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
