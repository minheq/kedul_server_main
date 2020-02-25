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
		_, err := employeeRoleService.CreateEmployeeRole(context.Background(), "location1", "role_name1", []Permission{}, actor)

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
	location := NewLocation("", "location1")
	employeeRole := NewEmployeeRole(location.ID, "role_name2", []Permission{})

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employeeRole", func(t *testing.T) {
		_, err := employeeRoleService.UpdateEmployeeRole(context.Background(), employeeRole.ID, "role_name3", []Permission{}, actor)

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

	location := NewLocation("", "location2")
	employeeRole := NewEmployeeRole(location.ID, "employeeRole4", []Permission{})

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
	location := NewLocation("", "location3")
	employeeRoleStore := &mockEmployeeRoleStore{}
	businessStore := &mockBusinessStore{}
	employeeStore := &mockEmployeeStore{}
	locationStore := &mockLocationStore{}
	permissionService := NewPermissionService(employeeRoleStore, employeeStore)
	locationService := NewLocationService(businessStore, locationStore, employeeRoleStore)
	employeeRoleService := NewEmployeeRoleService(employeeStore, employeeRoleStore)

	permissions := []Permission{permManageLocation}
	employeeRole := NewEmployeeRole(location.ID, "employeeRole4", permissions)

	err := employeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	employee := NewEmployee(location.ID, "employee1", employeeRole.ID)

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
		_, err := locationService.UpdateLocation(context.Background(), location.ID, "name", "profile_image_id", actor)

		if err != nil {
			t.Errorf("updating location should pass")
			return
		}
	})
}
