package app

import (
	"context"
	"testing"
)

type mockLocationStore struct {
	locations []*Location
}

func (s *mockLocationStore) GetLocationByID(ctx context.Context, id string) (*Location, error) {
	for _, b := range s.locations {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, nil
}

func (s *mockLocationStore) StoreLocation(ctx context.Context, location *Location) error {
	s.locations = append(s.locations, location)

	return nil
}

func (s *mockLocationStore) UpdateLocation(ctx context.Context, location *Location) error {
	for i, b := range s.locations {
		if b.ID == location.ID {
			s.locations[i] = location
			break
		}
	}

	return nil
}

func (s *mockLocationStore) DeleteLocation(ctx context.Context, location *Location) error {
	for i, b := range s.locations {
		if b.ID == location.ID {
			s.locations = append(s.locations[:i], s.locations[i+1:]...)
			break
		}
	}

	return nil
}

func TestCreateLocationHappyPath(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(locationStore, employeeRoleStore)

	t.Run("should create location", func(t *testing.T) {
		_, err := locationService.CreateLocation(context.Background(), "1", "location1")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateLocationHappyPath(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(locationStore, employeeRoleStore)
	actor := &mockActor{}

	business := NewBusiness("", "business1")
	location := NewLocation(business.ID, "location2")

	err := locationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update location", func(t *testing.T) {
		_, err := locationService.UpdateLocation(context.Background(), location.ID, "location3", "", actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteLocationHappyPath(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(locationStore, employeeRoleStore)
	business := NewBusiness("", "business2")
	location := NewLocation(business.ID, "location4")

	err := locationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update location", func(t *testing.T) {
		_, err := locationService.DeleteLocation(context.Background(), location.ID)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
