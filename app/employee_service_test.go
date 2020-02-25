package app

import (
	"context"
	"testing"
)

type mockEmployeeStore struct {
	employees []*Employee
}

func (s *mockEmployeeStore) GetEmployeesByEmployeeRoleID(ctx context.Context, employeeRoleID string) ([]*Employee, error) {
	employees := make([]*Employee, 0)

	for _, e := range s.employees {
		if e.EmployeeRoleID == employeeRoleID {
			employees = append(employees, e)
		}
	}

	return employees, nil
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
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)
	actor := &mockActor{}

	t.Run("should create employee", func(t *testing.T) {
		input := &CreateEmployeeInput{
			Name: "employee1",
		}
		_, err := employeeService.CreateEmployee(context.Background(), input, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateEmployeeHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)
	actor := &mockActor{}

	location := &Location{
		ID:   "1",
		Name: "location1",
	}
	employee := &Employee{
		ID:         "1",
		LocationID: location.ID,
		Name:       "employee2",
	}

	err := employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employee", func(t *testing.T) {
		input := &UpdateEmployeeInput{
			Name: "employee3",
		}
		_, err := employeeService.UpdateEmployee(context.Background(), employee.ID, input, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteEmployeeHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeService := NewEmployeeService(employeeStore, employeeRoleStore)
	location := &Location{
		ID:   "2",
		Name: "location2",
	}

	employeeRole := &EmployeeRole{
		ID:            "1",
		LocationID:    location.ID,
		Name:          "employee_role1",
		PermissionIDs: []string{},
	}

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	employee := &Employee{
		ID:             "4",
		LocationID:     location.ID,
		Name:           "employee4",
		EmployeeRoleID: employeeRole.ID,
	}
	actor := &mockActor{}

	err = employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should delete employee", func(t *testing.T) {
		_, err := employeeService.DeleteEmployee(context.Background(), employee.ID, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
