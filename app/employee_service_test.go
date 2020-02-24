package app

import (
	"context"
	"testing"
)

type mockEmployeeStore struct {
	employees []*Employee
}

func (s *mockEmployeeStore) GetEmployeeByUserIDAndLocationID(ctx context.Context, userID string, locationID string) (*Employee, error) {
	for _, e := range s.employees {
		if e.UserID == userID && e.LocationID == locationID {
			return e, nil
		}
	}

	return nil, nil
}

func (s *mockEmployeeStore) GetEmployeeByID(ctx context.Context, id string) (*Employee, error) {
	for _, e := range s.employees {
		if e.ID == id {
			return e, nil
		}
	}

	return nil, nil
}

func (s *mockEmployeeStore) StoreEmployee(ctx context.Context, employee *Employee) error {
	s.employees = append(s.employees, employee)

	return nil
}

func (s *mockEmployeeStore) UpdateEmployee(ctx context.Context, employee *Employee) error {
	for i, e := range s.employees {
		if e.ID == employee.ID {
			s.employees[i] = employee
			break
		}
	}

	return nil
}

func (s *mockEmployeeStore) DeleteEmployee(ctx context.Context, employee *Employee) error {
	for i, e := range s.employees {
		if e.ID == employee.ID {
			s.employees = append(s.employees[:i], s.employees[i+1:]...)
			break
		}
	}

	return nil
}

func TestCreateEmployeeHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore)
	actor := &mockActor{}

	t.Run("should create employee", func(t *testing.T) {
		_, err := employeeService.CreateEmployee(context.Background(), "1", "employee1", "1", actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateEmployeeHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore)
	actor := &mockActor{}

	business := NewBusiness("", "business1")
	employee := NewEmployee(business.ID, "employee2", "1")

	err := employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employee", func(t *testing.T) {
		_, err := employeeService.UpdateEmployee(context.Background(), employee.ID, "employee3", "", actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteEmployeeHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore)
	business := NewBusiness("", "business2")
	employee := NewEmployee(business.ID, "employee4", "1")
	actor := &mockActor{}

	err := employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employee", func(t *testing.T) {
		_, err := employeeService.DeleteEmployee(context.Background(), employee.ID, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
