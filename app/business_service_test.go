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

func TestCreateBusinessHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	businessService := NewBusinessService(businessStore)

	t.Run("should create business", func(t *testing.T) {
		_, err := businessService.CreateBusiness(context.Background(), "1", "business1")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateBusinessHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	businessService := NewBusinessService(businessStore)

	currentUser := auth.NewUser("", "")
	business := NewBusiness(currentUser.ID, "business2")

	err := businessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update business", func(t *testing.T) {
		_, err := businessService.UpdateBusiness(context.Background(), business.ID, "business3", "", currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteBusinessHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	businessService := NewBusinessService(businessStore)
	currentUser := auth.NewUser("", "")
	business := NewBusiness(currentUser.ID, "business4")

	err := businessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update business", func(t *testing.T) {
		_, err := businessService.DeleteBusiness(context.Background(), business.ID, currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
