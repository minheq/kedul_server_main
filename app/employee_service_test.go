package app

import (
	"context"
	"testing"
)

type mockEmployeeStore struct {
	employees []*Employee
}

func (s *mockEmployeeStore) GetEmployeeByID(ctx context.Context, id string) (*Employee, error) {
	for _, b := range s.employees {
		if b.ID == id {
			return b, nil
		}
	}

	return nil, nil
}

func (s *mockEmployeeStore) StoreEmployee(ctx context.Context, employee *Employee) error {
	s.employees = append(s.employees, employee)

	return nil
}

func (s *mockEmployeeStore) UpdateEmployee(ctx context.Context, employee *Employee) error {
	for i, b := range s.employees {
		if b.ID == employee.ID {
			s.employees[i] = employee
			break
		}
	}

	return nil
}

func (s *mockEmployeeStore) DeleteEmployee(ctx context.Context, employee *Employee) error {
	for i, b := range s.employees {
		if b.ID == employee.ID {
			s.employees = append(s.employees[:i], s.employees[i+1:]...)
			break
		}
	}

	return nil
}

func TestCreateEmployeeHappyPath(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)

	t.Run("should create employee", func(t *testing.T) {
		_, err := employeeService.CreateEmployee(context.Background(), "1", "employee1", "1")

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateEmployeeHappyPath(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)
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
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeStore := &mockEmployeeStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)
	business := NewBusiness("", "business2")
	employee := NewEmployee(business.ID, "employee4", "1")

	err := employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employee", func(t *testing.T) {
		_, err := employeeService.DeleteEmployee(context.Background(), employee.ID)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
