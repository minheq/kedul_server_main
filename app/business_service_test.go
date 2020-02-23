package app

import (
	"context"
	"testing"

	"github.com/minheq/kedul_server_main/auth"
)

type mockBusinessStore struct {
	businesses []*Business
}

func (s *mockBusinessStore) GetBusinessByID(ctx context.Context, id string) (*Business, error) {
	for _, b := range s.businesses {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, nil
}

func (s *mockBusinessStore) GetBusinessByName(ctx context.Context, name string) (*Business, error) {
	for _, b := range s.businesses {
		if b.Name == name {
			return b, nil
		}
	}

	return nil, nil
}

func (s *mockBusinessStore) StoreBusiness(ctx context.Context, business *Business) error {
	s.businesses = append(s.businesses, business)

	return nil
}

func (s *mockBusinessStore) UpdateBusiness(ctx context.Context, business *Business) error {
	for i, b := range s.businesses {
		if b.ID == business.ID {
			s.businesses[i] = business
			break
		}
	}

	return nil
}

func (s *mockBusinessStore) DeleteBusiness(ctx context.Context, business *Business) error {
	for i, b := range s.businesses {
		if b.ID == business.ID {
			s.businesses = append(s.businesses[:i], s.businesses[i+1:]...)
			break
		}
	}

	return nil
}

var (
	testBusinessStore   = &mockBusinessStore{}
	testBusinessService = NewBusinessService(testBusinessStore)
)

func TestCreateBusinessHappyPath(t *testing.T) {
	t.Run("should create business", func(t *testing.T) {
		_, err := testBusinessService.CreateBusiness(context.Background(), "1", "business1")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateBusinessHappyPath(t *testing.T) {
	currentUser := auth.NewUser("", "")
	business := NewBusiness(currentUser.ID, "business2")

	err := testBusinessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update business", func(t *testing.T) {
		_, err := testBusinessService.UpdateBusiness(context.Background(), business.ID, "business3", "", currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteBusinessHappyPath(t *testing.T) {
	currentUser := auth.NewUser("", "")
	business := NewBusiness(currentUser.ID, "business4")

	err := testBusinessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update business", func(t *testing.T) {
		_, err := testBusinessService.DeleteBusiness(context.Background(), business.ID, currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
