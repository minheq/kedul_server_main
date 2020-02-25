package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/minheq/kedul_server_main/errors"
)

type mockEmployeeRoleStore struct {
	employeeRoles []*EmployeeRole
}

func (s *mockEmployeeRoleStore) GetEmployeeRoleByID(ctx context.Context, id string) (*EmployeeRole, error) {
	for _, employeeRole := range s.employeeRoles {
		if employeeRole.ID == id {
			permissions, err := getPermissionsByPermissionIDs(employeeRole.PermissionIDs)

			if err != nil {
				return nil, fmt.Errorf("failed to get permissions")
			}

			employeeRole.Permissions = permissions

			return employeeRole, nil
		}
	}

	return nil, nil
}

func (s *mockEmployeeRoleStore) StoreEmployeeRole(ctx context.Context, employee *EmployeeRole) error {
	s.employeeRoles = append(s.employeeRoles, employee)

	return nil
}

func (s *mockEmployeeRoleStore) UpdateEmployeeRole(ctx context.Context, employee *EmployeeRole) error {
	for i, b := range s.employeeRoles {
		if b.ID == employee.ID {
			s.employeeRoles[i] = employee
			break
		}
	}

	return nil
}

func (s *mockEmployeeRoleStore) DeleteEmployeeRole(ctx context.Context, employee *EmployeeRole) error {
	for i, b := range s.employeeRoles {
		if b.ID == employee.ID {
			s.employeeRoles = append(s.employeeRoles[:i], s.employeeRoles[i+1:]...)
			break
		}
	}

	return nil
}

type mockActor struct{}

func (m *mockActor) can(ctx context.Context, operation Operation) error {
	return nil
}

func TestCreateEmployeeRoleHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeRoleService := NewEmployeeRoleService(employeeStore, employeeRoleStore)
	actor := &mockActor{}

	t.Run("should create employee", func(t *testing.T) {
		input := &CreateEmployeeRoleInput{
			Name:          "role_name1",
			PermissionIDs: []string{},
		}
		_, err := employeeRoleService.CreateEmployeeRole(context.Background(), input, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateEmployeeRoleHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeRoleService := NewEmployeeRoleService(employeeStore, employeeRoleStore)
	actor := &mockActor{}
	location := &Location{
		ID:   "1",
		Name: "location1",
	}
	employeeRole := &EmployeeRole{
		LocationID:    location.ID,
		Name:          "role_name2",
		PermissionIDs: []string{},
	}

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employeeRole", func(t *testing.T) {
		input := &UpdateEmployeeRoleInput{
			Name: "role_name3",
		}
		_, err := employeeRoleService.UpdateEmployeeRole(context.Background(), employeeRole.ID, input, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteEmployeeRoleHappyPath(t *testing.T) {
	employeeStore := &mockEmployeeStore{}
	employeeRoleStore := &mockEmployeeRoleStore{}
	employeeRoleService := NewEmployeeRoleService(employeeStore, employeeRoleStore)
	actor := &mockActor{}

	location := &Location{
		ID:   "2",
		Name: "location2",
	}
	employeeRole := &EmployeeRole{
		LocationID:    location.ID,
		Name:          "role_name4",
		PermissionIDs: []string{},
	}

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should delete employeeRole", func(t *testing.T) {
		_, err := employeeRoleService.DeleteEmployeeRole(context.Background(), employeeRole.ID, actor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestEmployeeRolePermissions(t *testing.T) {
	employeeRoleStore := &mockEmployeeRoleStore{}
	businessStore := &mockBusinessStore{}
	employeeStore := &mockEmployeeStore{}
	locationStore := &mockLocationStore{}
	permissionService := NewPermissionService(employeeRoleStore, employeeStore)
	locationService := NewLocationService(businessStore, locationStore, employeeStore, employeeRoleStore)
	employeeRoleService := NewEmployeeRoleService(employeeStore, employeeRoleStore)

	location := &Location{
		ID:   "3",
		Name: "location3",
	}
	employeeRole := &EmployeeRole{
		LocationID:    location.ID,
		Name:          "role_name5",
		PermissionIDs: []string{permManageLocation.ID},
	}

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	employee := &Employee{
		LocationID:     location.ID,
		Name:           "employee1",
		EmployeeRoleID: employeeRole.ID,
	}

	err = employeeStore.StoreEmployee(context.Background(), employee)

	if err != nil {
		t.Error(err)
		return
	}

	err = locationStore.StoreLocation(context.Background(), location)

	if err != nil {
		t.Error(err)
		return
	}

	actor, err := permissionService.GetEmployeeActor(context.Background(), employee.UserID, employee.LocationID)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should not be able to updateEmployeeRole", func(t *testing.T) {
		_, err = employeeRoleService.DeleteEmployeeRole(context.Background(), employeeRole.ID, actor)

		if errors.Is(errors.KindUnauthorized, err) == false {
			t.Errorf("deleting employee role should fail due insufficient permissions")
			return
		}
	})

	t.Run("should be able to updateLocation", func(t *testing.T) {
		input := &UpdateLocationInput{
			Name: "new name",
		}
		_, err := locationService.UpdateLocation(context.Background(), location.ID, input, actor)

		if err != nil {
			t.Errorf("updating location should pass")
			return
		}
	})
}
