package app

import (
	"context"
	"testing"

	"github.com/minheq/kedul_server_main/auth"
)

type mockBusinessStore struct {
	businesses []*Business
}

func (s *mockBusinessStore) GetBusinessesByIDs(ctx context.Context, ids []string) ([]*Business, error) {
	businesses := []*Business{}

	for _, l := range s.businesses {
		for _, id := range ids {
			if l.ID == id {
				businesses = append(businesses, l)
			}
		}
	}

	return businesses, nil
}

func (s *mockBusinessStore) GetBusinessesByUserID(ctx context.Context, userID string) ([]*Business, error) {
	businesses := make([]*Business, 0)

	for _, e := range s.businesses {
		if e.UserID == userID {
			businesses = append(businesses, e)
		}
	}

	return businesses, nil
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
	locationStore := &mockLocationStore{}
	employeeStore := &mockEmployeeStore{}
	businessService := NewBusinessService(businessStore, locationStore, employeeStore)

	t.Run("should create business", func(t *testing.T) {
		input := &CreateBusinessInput{
			Name: "business1",
		}

		_, err := businessService.CreateBusiness(context.Background(), "1", input)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateBusinessHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	locationStore := &mockLocationStore{}
	employeeStore := &mockEmployeeStore{}
	businessService := NewBusinessService(businessStore, locationStore, employeeStore)

	currentUser := &auth.User{
		ID: "1",
	}
	business := &Business{
		ID:     "1",
		UserID: currentUser.ID,
		Name:   "business2",
	}

	err := businessStore.StoreBusiness(context.Background(), business)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update business", func(t *testing.T) {
		input := &UpdateBusinessInput{
			Name: "new business2",
		}

		_, err := businessService.UpdateBusiness(context.Background(), business.ID, input, currentUser)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteBusinessHappyPath(t *testing.T) {
	businessStore := &mockBusinessStore{}
	locationStore := &mockLocationStore{}
	employeeStore := &mockEmployeeStore{}
	businessService := NewBusinessService(businessStore, locationStore, employeeStore)
	currentUser := &auth.User{
		ID: "2",
	}
	business := &Business{
		ID:     "2",
		UserID: currentUser.ID,
		Name:   "business4",
	}

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
