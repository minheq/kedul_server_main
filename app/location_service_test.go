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

var (
	testLocationStore   = &mockLocationStore{}
	testLocationService = NewLocationService(testLocationStore)
)

func TestCreateLocationHappyPath(t *testing.T) {
	t.Run("should create location", func(t *testing.T) {
		_, err := testLocationService.CreateLocation(context.Background(), "1", "location1")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateLocationHappyPath(t *testing.T) {
	business := NewBusiness("", "business1")
	location := NewLocation(business.ID, "location2")

	err := testLocationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update location", func(t *testing.T) {
		_, err := testLocationService.UpdateLocation(context.Background(), location.ID, "location3", "")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteLocationHappyPath(t *testing.T) {
	business := NewBusiness("", "business2")
	location := NewLocation(business.ID, "location4")

	err := testLocationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update location", func(t *testing.T) {
		_, err := testLocationService.DeleteLocation(context.Background(), location.ID)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
