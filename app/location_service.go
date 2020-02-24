package app

import (
	"context"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// LocationService ...
type LocationService struct {
	locationStore     LocationStore
	employeeRoleStore EmployeeRoleStore
}

// NewLocationService constructor for AuthService
func NewLocationService(locationStore LocationStore, employeeRoleStore EmployeeRoleStore) LocationService {
	return LocationService{locationStore: locationStore, employeeRoleStore: employeeRoleStore}
}

// GetLocationByID ...
func (s *LocationService) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	const op = "app/locationService.GetLocationByID"

	location, err := s.locationStore.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	return location, nil
}

// CreateLocation creates location
func (s *LocationService) CreateLocation(ctx context.Context, businessID string, name string) (*Location, error) {
	const op = "app/locationService.CreateLocation"

	location := NewLocation(businessID, name)

	err := s.locationStore.StoreLocation(ctx, location)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to locationStore location")
	}

	ownerRole := NewEmployeeRole(location.ID, name, defaultOwnerRolePermissions)
	adminRole := NewEmployeeRole(location.ID, name, defaultAdminRolePermissions)
	managerRole := NewEmployeeRole(location.ID, name, defaultManagerRolePermissions)
	receptionistRole := NewEmployeeRole(location.ID, name, defaultReceptionistRolePermissions)
	specialistRole := NewEmployeeRole(location.ID, name, defaultSpecialistRolePermissions)

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

	return location, nil
}

// UpdateLocation updates location
func (s *LocationService) UpdateLocation(ctx context.Context, id string, name string, profileImageID string, actor Actor) (*Location, error) {
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
	location.Name = name
	location.ProfileImageID = profileImageID

	err = s.locationStore.UpdateLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}

// DeleteLocation updates location
func (s *LocationService) DeleteLocation(ctx context.Context, id string) (*Location, error) {
	const op = "app/locationService.DeleteLocation"

	location, err := s.locationStore.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	if location == nil {
		return nil, errors.NotFound(op)
	}

	err = s.locationStore.DeleteLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}
