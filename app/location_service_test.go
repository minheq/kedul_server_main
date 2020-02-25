package app

import (
	"context"
	"testing"

	"github.com/minheq/kedul_server_main/auth"
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
	businessStore := &mockBusinessStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(businessStore, locationStore, employeeRoleStore)

	currentUser := auth.NewUser("", "")
	business := NewBusiness(currentUser.ID, "business1")

	err := businessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should create location", func(t *testing.T) {
		_, err := locationService.CreateLocation(context.Background(), business.ID, "location1", currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateLocationHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(businessStore, locationStore, employeeRoleStore)
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
	businessStore := &mockBusinessStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	locationStore := &mockLocationStore{}
	locationService := NewLocationService(businessStore, locationStore, employeeRoleStore)
	currentUser := auth.NewUser("", "")

	business := NewBusiness(currentUser.ID, "business2")
	err := businessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}
	location := NewLocation(business.ID, "location4")

	err = locationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update location", func(t *testing.T) {
		_, err := locationService.DeleteLocation(context.Background(), location.ID, currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
