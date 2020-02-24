package app

import (
	"context"
	"fmt"
	"testing"
)

type mockEmployeeRoleStore struct {
	employeeRoles []*EmployeeRole
}

func (s *mockEmployeeRoleStore) GetEmployeeRoleByID(ctx context.Context, id string) (*EmployeeRole, error) {
	for _, employeeRole := range s.employeeRoles {
		if employeeRole.ID == id {
			permissions := []Permission{}

			for _, permissionID := range employeeRole.PermissionIDs {
				found := false

				for _, permission := range permissionList {
					if permission.ID == id {
						permissions = append(permissions, permission)
						found = true
					}
				}

				if found == false {
					return nil, fmt.Errorf("permission for permissionID=%s not found", permissionID)
				}
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

type mockEmployeeRoleActor struct{}

func (m *mockEmployeeRoleActor) can(ctx context.Context, operation Operation) error {
	return nil
}

var (
	testEmployeeRoleStore   = &mockEmployeeRoleStore{}
	testEmployeeRoleService = NewEmployeeRoleService(testEmployeeRoleStore)
	testEmployeeRoleActor   = &mockEmployeeRoleActor{}
)

func TestCreateEmployeeRoleHappyPath(t *testing.T) {
	t.Run("should create employee", func(t *testing.T) {
		_, err := testEmployeeRoleService.CreateEmployeeRole(context.Background(), "location1", "role_name1", []Permission{}, testEmployeeRoleActor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestUpdateEmployeeRoleHappyPath(t *testing.T) {
	location := NewLocation("", "location1")
	employeeRole := NewEmployeeRole(location.ID, "role_name2", []Permission{})

	err := testEmployeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employeeRole", func(t *testing.T) {
		_, err := testEmployeeRoleService.UpdateEmployeeRole(context.Background(), employeeRole.ID, "role_name3", []Permission{}, testEmployeeRoleActor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}

func TestDeleteEmployeeRoleHappyPath(t *testing.T) {
	location := NewLocation("", "location2")
	employeeRole := NewEmployeeRole(location.ID, "employeeRole4", []Permission{})

	err := testEmployeeRoleStore.StoreEmployeeRole(context.Background(), employeeRole)

	if err != nil {
		t.Error(err)
		return
	}

	t.Run("should update employeeRole", func(t *testing.T) {
		_, err := testEmployeeRoleService.DeleteEmployeeRole(context.Background(), employeeRole.ID, testEmployeeRoleActor)

		if err != nil {
			t.Error(err)
			return
		}
	})
}
