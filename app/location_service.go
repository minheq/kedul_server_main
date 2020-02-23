package app

import (
	"context"
	"time"

	"github.com/minheq/kedul_server_main/errors"
)

// LocationService ...
type LocationService struct {
	store LocationStore
}

// NewLocationService constructor for AuthService
func NewLocationService(store LocationStore) LocationService {
	return LocationService{store: store}
}

// GetLocationByID ...
func (s *LocationService) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	const op = "location/service.CreateLocation"

	location, err := s.store.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	return location, nil
}

// CreateLocation creates location
func (s *LocationService) CreateLocation(ctx context.Context, businessID string, name string) (*Location, error) {
	const op = "location/service.CreateLocation"

	location := NewLocation(businessID, name)

	err := s.store.StoreLocation(ctx, location)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to store location")
	}

	return location, nil
}

// UpdateLocation updates location
func (s *LocationService) UpdateLocation(ctx context.Context, id string, name string, profileImageID string) (*Location, error) {
	const op = "location/service.UpdateLocation"

	location, err := s.store.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	if location == nil {
		return nil, errors.NotFound(op)
	}

	location.UpdatedAt = time.Now()
	location.Name = name
	location.ProfileImageID = profileImageID

	err = s.store.UpdateLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}

// DeleteLocation updates location
func (s *LocationService) DeleteLocation(ctx context.Context, id string) (*Location, error) {
	const op = "location/service.DeleteLocation"

	location, err := s.store.GetLocationByID(ctx, id)

	if err != nil {
		return nil, errors.Wrap(op, err, "failed to get location by id")
	}

	if location == nil {
		return nil, errors.NotFound(op)
	}

	err = s.store.DeleteLocation(ctx, location)

	if err != nil {
		return nil, errors.Unexpected(op, err, "failed to update location")
	}

	return location, nil
}
