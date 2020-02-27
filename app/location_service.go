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

var (
	defaultOwnerRolePermissionIDs = []string{
		permManageLocation.ID,
		permManageEmployeeRole.ID,
		permManageEmployee.ID,
	}
	defaultAdminRolePermissionIDs        = []string{}
	defaultManagerRolePermissionIDs      = []string{}
	defaultReceptionistRolePermissionIDs = []string{}
	defaultSpecialistRolePermissionIDs   = []string{}
)

// Location ...
type Location struct {
	ID             string    `json:"id"`
	BusinessID     string    `json:"business_id"`
	Name           string    `json:"name"`
	ProfileImageID string    `json:"profile_image_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// LocationService ...
type LocationService struct {
	businessStore     BusinessStore
	locationStore     LocationStore
	employeeStore     EmployeeStore
	employeeRoleStore EmployeeRoleStore
}

// NewLocationService constructor for AuthService
func NewLocationService(businessStore BusinessStore, locationStore LocationStore, employeeStore EmployeeStore, employeeRoleStore EmployeeRoleStore) LocationService {
	return LocationService{businessStore: businessStore, locationStore: locationStore, employeeStore: employeeStore, employeeRoleStore: employeeRoleStore}
}

// GetLocationByID ...
func (s *LocationService) GetLocationByID(ctx context.Context, id string, actor Actor) (*Location, error) {
	const op = "app/locationService.GetLocationByID"

	location, err := s.locationStore.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	return location, nil
}

// GetLocationsByUserIDAndBusinessID ...
func (s *LocationService) GetLocationsByUserIDAndBusinessID(ctx context.Context, userID string, businessID string, currentUser *auth.User) ([]*Location, error) {
	const op = "app/locationService.GetLocationsByUserID"

	if userID != currentUser.ID {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

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

	businessLocations := []*Location{}

	for _, location := range locations {
		if location.BusinessID == businessID {
			businessLocations = append(businessLocations, location)
		}
	}

	return businessLocations, nil
}

// CreateLocationInput ...
type CreateLocationInput struct {
	BusinessID     string `json:"business_id"`
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// CreateLocation creates location
func (s *LocationService) CreateLocation(ctx context.Context, input *CreateLocationInput, currentUser *auth.User) (*Location, error) {
	const op = "app/locationService.CreateLocation"

	businesses, err := s.businessStore.GetBusinessesByUserID(ctx, currentUser.ID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to to get businesses by user id")
	}

	isOwner := false

	for _, business := range businesses {
		if business.ID == input.BusinessID {
			isOwner = true
		}
	}

	if isOwner == false {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	now := time.Now()

	location := &Location{
		ID:         uuid.Must(uuid.New(), nil).String(),
		BusinessID: input.BusinessID,
		Name:       input.Name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	err = s.locationStore.StoreLocation(ctx, location)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to locationStore location")
	}

	ownerRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    location.ID,
		Name:          "admin",
		PermissionIDs: defaultOwnerRolePermissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	adminRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    location.ID,
		Name:          "admin",
		PermissionIDs: defaultAdminRolePermissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	managerRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    location.ID,
		Name:          "manager",
		PermissionIDs: defaultManagerRolePermissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	receptionistRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    location.ID,
		Name:          "receptionist",
		PermissionIDs: defaultReceptionistRolePermissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	specialistRole := &EmployeeRole{
		ID:            uuid.Must(uuid.New(), nil).String(),
		LocationID:    location.ID,
		Name:          "specialist",
		PermissionIDs: defaultSpecialistRolePermissionIDs,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, ownerRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default owner role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, adminRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default admin role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, managerRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default manager role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, receptionistRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default receptionist role")
	}

	err = s.employeeRoleStore.StoreEmployeeRole(ctx, specialistRole)
	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create default employee role")
	}

	owner := &Employee{
		ID:             uuid.Must(uuid.New(), nil).String(),
		LocationID:     location.ID,
		Name:           currentUser.FullName,
		EmployeeRoleID: ownerRole.ID,
		UserID:         currentUser.ID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err = s.employeeStore.StoreEmployee(ctx, owner)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to create location owner")
	}

	return location, nil
}

// UpdateLocationInput ...
type UpdateLocationInput struct {
	Name           string `json:"name"`
	ProfileImageID string `json:"profile_image_id"`
}

// UpdateLocation updates location
func (s *LocationService) UpdateLocation(ctx context.Context, id string, input *UpdateLocationInput, actor Actor) (*Location, error) {
	const op = "app/locationService.UpdateLocation"

	err := actor.can(ctx, opUpdateLocation)

	if err != nil {
		return nil, errors.Unauthorized(op, err)
	}

	location, err := s.locationStore.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	if location == nil {
		return nil, errors.NotFound(op)
	}

	location.UpdatedAt = time.Now()
	if input.Name != "" {
		location.Name = strings.TrimSpace(input.Name)
	}
	if input.ProfileImageID != "" {
		location.ProfileImageID = input.ProfileImageID
	}

	err = s.locationStore.UpdateLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}

// DeleteLocation updates location
func (s *LocationService) DeleteLocation(ctx context.Context, id string, currentUser *auth.User) (*Location, error) {
	const op = "app/locationService.DeleteLocation"

	location, err := s.locationStore.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	if location == nil {
		return nil, errors.NotFound(op)
	}

	businesses, err := s.businessStore.GetBusinessesByUserID(ctx, currentUser.ID)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to to get businesses by user id")
	}

	isOwner := false

	for _, business := range businesses {
		if business.ID == location.BusinessID {
			isOwner = true
		}
	}

	if isOwner == false {
		return nil, errors.Unauthorized(op, fmt.Errorf("current user not owner"))
	}

	err = s.locationStore.DeleteLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}
