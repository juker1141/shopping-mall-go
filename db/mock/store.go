// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juker1141/shopping-mall-go/db/sqlc (interfaces: Store)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	db "github.com/juker1141/shopping-mall-go/db/sqlc"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// CreateAdminUser mocks base method.
func (m *MockStore) CreateAdminUser(arg0 context.Context, arg1 db.CreateAdminUserParams) (db.AdminUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminUser", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdminUser indicates an expected call of CreateAdminUser.
func (mr *MockStoreMockRecorder) CreateAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminUser", reflect.TypeOf((*MockStore)(nil).CreateAdminUser), arg0, arg1)
}

// CreateAdminUserRole mocks base method.
func (m *MockStore) CreateAdminUserRole(arg0 context.Context, arg1 db.CreateAdminUserRoleParams) (db.AdminUserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminUserRole", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdminUserRole indicates an expected call of CreateAdminUserRole.
func (mr *MockStoreMockRecorder) CreateAdminUserRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminUserRole", reflect.TypeOf((*MockStore)(nil).CreateAdminUserRole), arg0, arg1)
}

// CreateAdminUserTx mocks base method.
func (m *MockStore) CreateAdminUserTx(arg0 context.Context, arg1 db.CreateAdminUserTxParams) (db.AdminUserTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAdminUserTx", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUserTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAdminUserTx indicates an expected call of CreateAdminUserTx.
func (mr *MockStoreMockRecorder) CreateAdminUserTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAdminUserTx", reflect.TypeOf((*MockStore)(nil).CreateAdminUserTx), arg0, arg1)
}

// CreatePermission mocks base method.
func (m *MockStore) CreatePermission(arg0 context.Context, arg1 string) (db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePermission", arg0, arg1)
	ret0, _ := ret[0].(db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePermission indicates an expected call of CreatePermission.
func (mr *MockStoreMockRecorder) CreatePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePermission", reflect.TypeOf((*MockStore)(nil).CreatePermission), arg0, arg1)
}

// CreateRole mocks base method.
func (m *MockStore) CreateRole(arg0 context.Context, arg1 string) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockStoreMockRecorder) CreateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockStore)(nil).CreateRole), arg0, arg1)
}

// CreateRolePermission mocks base method.
func (m *MockStore) CreateRolePermission(arg0 context.Context, arg1 db.CreateRolePermissionParams) (db.RolePermission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRolePermission", arg0, arg1)
	ret0, _ := ret[0].(db.RolePermission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRolePermission indicates an expected call of CreateRolePermission.
func (mr *MockStoreMockRecorder) CreateRolePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRolePermission", reflect.TypeOf((*MockStore)(nil).CreateRolePermission), arg0, arg1)
}

// CreateRoleTx mocks base method.
func (m *MockStore) CreateRoleTx(arg0 context.Context, arg1 db.CreateRoleTxParams) (db.RoleTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRoleTx", arg0, arg1)
	ret0, _ := ret[0].(db.RoleTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRoleTx indicates an expected call of CreateRoleTx.
func (mr *MockStoreMockRecorder) CreateRoleTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRoleTx", reflect.TypeOf((*MockStore)(nil).CreateRoleTx), arg0, arg1)
}

// DeleteAdminUser mocks base method.
func (m *MockStore) DeleteAdminUser(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdminUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdminUser indicates an expected call of DeleteAdminUser.
func (mr *MockStoreMockRecorder) DeleteAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminUser", reflect.TypeOf((*MockStore)(nil).DeleteAdminUser), arg0, arg1)
}

// DeleteAdminUserRoleByAdminUserId mocks base method.
func (m *MockStore) DeleteAdminUserRoleByAdminUserId(arg0 context.Context, arg1 pgtype.Int4) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdminUserRoleByAdminUserId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdminUserRoleByAdminUserId indicates an expected call of DeleteAdminUserRoleByAdminUserId.
func (mr *MockStoreMockRecorder) DeleteAdminUserRoleByAdminUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminUserRoleByAdminUserId", reflect.TypeOf((*MockStore)(nil).DeleteAdminUserRoleByAdminUserId), arg0, arg1)
}

// DeleteAdminUserRoleByRoleId mocks base method.
func (m *MockStore) DeleteAdminUserRoleByRoleId(arg0 context.Context, arg1 pgtype.Int4) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdminUserRoleByRoleId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAdminUserRoleByRoleId indicates an expected call of DeleteAdminUserRoleByRoleId.
func (mr *MockStoreMockRecorder) DeleteAdminUserRoleByRoleId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminUserRoleByRoleId", reflect.TypeOf((*MockStore)(nil).DeleteAdminUserRoleByRoleId), arg0, arg1)
}

// DeleteAdminUserTx mocks base method.
func (m *MockStore) DeleteAdminUserTx(arg0 context.Context, arg1 db.DeleteAdminUserTxParams) (db.DeleteRoleTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAdminUserTx", arg0, arg1)
	ret0, _ := ret[0].(db.DeleteRoleTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteAdminUserTx indicates an expected call of DeleteAdminUserTx.
func (mr *MockStoreMockRecorder) DeleteAdminUserTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAdminUserTx", reflect.TypeOf((*MockStore)(nil).DeleteAdminUserTx), arg0, arg1)
}

// DeletePermission mocks base method.
func (m *MockStore) DeletePermission(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePermission indicates an expected call of DeletePermission.
func (mr *MockStoreMockRecorder) DeletePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePermission", reflect.TypeOf((*MockStore)(nil).DeletePermission), arg0, arg1)
}

// DeleteRole mocks base method.
func (m *MockStore) DeleteRole(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockStoreMockRecorder) DeleteRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockStore)(nil).DeleteRole), arg0, arg1)
}

// DeleteRolePermissionByPermissionId mocks base method.
func (m *MockStore) DeleteRolePermissionByPermissionId(arg0 context.Context, arg1 pgtype.Int4) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRolePermissionByPermissionId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRolePermissionByPermissionId indicates an expected call of DeleteRolePermissionByPermissionId.
func (mr *MockStoreMockRecorder) DeleteRolePermissionByPermissionId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRolePermissionByPermissionId", reflect.TypeOf((*MockStore)(nil).DeleteRolePermissionByPermissionId), arg0, arg1)
}

// DeleteRolePermissionByRoleId mocks base method.
func (m *MockStore) DeleteRolePermissionByRoleId(arg0 context.Context, arg1 pgtype.Int4) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRolePermissionByRoleId", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRolePermissionByRoleId indicates an expected call of DeleteRolePermissionByRoleId.
func (mr *MockStoreMockRecorder) DeleteRolePermissionByRoleId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRolePermissionByRoleId", reflect.TypeOf((*MockStore)(nil).DeleteRolePermissionByRoleId), arg0, arg1)
}

// DeleteRoleTx mocks base method.
func (m *MockStore) DeleteRoleTx(arg0 context.Context, arg1 db.DeleteRoleTxParams) (db.DeleteRoleTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRoleTx", arg0, arg1)
	ret0, _ := ret[0].(db.DeleteRoleTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteRoleTx indicates an expected call of DeleteRoleTx.
func (mr *MockStoreMockRecorder) DeleteRoleTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRoleTx", reflect.TypeOf((*MockStore)(nil).DeleteRoleTx), arg0, arg1)
}

// GetAdminUser mocks base method.
func (m *MockStore) GetAdminUser(arg0 context.Context, arg1 int64) (db.AdminUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminUser", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminUser indicates an expected call of GetAdminUser.
func (mr *MockStoreMockRecorder) GetAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminUser", reflect.TypeOf((*MockStore)(nil).GetAdminUser), arg0, arg1)
}

// GetAdminUserByAccount mocks base method.
func (m *MockStore) GetAdminUserByAccount(arg0 context.Context, arg1 string) (db.AdminUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminUserByAccount", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminUserByAccount indicates an expected call of GetAdminUserByAccount.
func (mr *MockStoreMockRecorder) GetAdminUserByAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminUserByAccount", reflect.TypeOf((*MockStore)(nil).GetAdminUserByAccount), arg0, arg1)
}

// GetAdminUserRole mocks base method.
func (m *MockStore) GetAdminUserRole(arg0 context.Context, arg1 db.GetAdminUserRoleParams) (db.AdminUserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminUserRole", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminUserRole indicates an expected call of GetAdminUserRole.
func (mr *MockStoreMockRecorder) GetAdminUserRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminUserRole", reflect.TypeOf((*MockStore)(nil).GetAdminUserRole), arg0, arg1)
}

// GetAdminUsersCount mocks base method.
func (m *MockStore) GetAdminUsersCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminUsersCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminUsersCount indicates an expected call of GetAdminUsersCount.
func (mr *MockStoreMockRecorder) GetAdminUsersCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminUsersCount", reflect.TypeOf((*MockStore)(nil).GetAdminUsersCount), arg0)
}

// GetPermission mocks base method.
func (m *MockStore) GetPermission(arg0 context.Context, arg1 int64) (db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPermission", arg0, arg1)
	ret0, _ := ret[0].(db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPermission indicates an expected call of GetPermission.
func (mr *MockStoreMockRecorder) GetPermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPermission", reflect.TypeOf((*MockStore)(nil).GetPermission), arg0, arg1)
}

// GetRole mocks base method.
func (m *MockStore) GetRole(arg0 context.Context, arg1 int64) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRole indicates an expected call of GetRole.
func (mr *MockStoreMockRecorder) GetRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRole", reflect.TypeOf((*MockStore)(nil).GetRole), arg0, arg1)
}

// GetRolePermission mocks base method.
func (m *MockStore) GetRolePermission(arg0 context.Context, arg1 db.GetRolePermissionParams) (db.RolePermission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRolePermission", arg0, arg1)
	ret0, _ := ret[0].(db.RolePermission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRolePermission indicates an expected call of GetRolePermission.
func (mr *MockStoreMockRecorder) GetRolePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolePermission", reflect.TypeOf((*MockStore)(nil).GetRolePermission), arg0, arg1)
}

// GetRolesCount mocks base method.
func (m *MockStore) GetRolesCount(arg0 context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRolesCount", arg0)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRolesCount indicates an expected call of GetRolesCount.
func (mr *MockStoreMockRecorder) GetRolesCount(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRolesCount", reflect.TypeOf((*MockStore)(nil).GetRolesCount), arg0)
}

// ListAdminUserRoleByAdminUserId mocks base method.
func (m *MockStore) ListAdminUserRoleByAdminUserId(arg0 context.Context, arg1 pgtype.Int4) ([]db.AdminUserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAdminUserRoleByAdminUserId", arg0, arg1)
	ret0, _ := ret[0].([]db.AdminUserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAdminUserRoleByAdminUserId indicates an expected call of ListAdminUserRoleByAdminUserId.
func (mr *MockStoreMockRecorder) ListAdminUserRoleByAdminUserId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAdminUserRoleByAdminUserId", reflect.TypeOf((*MockStore)(nil).ListAdminUserRoleByAdminUserId), arg0, arg1)
}

// ListAdminUserRoleByRoleId mocks base method.
func (m *MockStore) ListAdminUserRoleByRoleId(arg0 context.Context, arg1 pgtype.Int4) ([]db.AdminUserRole, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAdminUserRoleByRoleId", arg0, arg1)
	ret0, _ := ret[0].([]db.AdminUserRole)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAdminUserRoleByRoleId indicates an expected call of ListAdminUserRoleByRoleId.
func (mr *MockStoreMockRecorder) ListAdminUserRoleByRoleId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAdminUserRoleByRoleId", reflect.TypeOf((*MockStore)(nil).ListAdminUserRoleByRoleId), arg0, arg1)
}

// ListAdminUsers mocks base method.
func (m *MockStore) ListAdminUsers(arg0 context.Context, arg1 db.ListAdminUsersParams) ([]db.AdminUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListAdminUsers", arg0, arg1)
	ret0, _ := ret[0].([]db.AdminUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListAdminUsers indicates an expected call of ListAdminUsers.
func (mr *MockStoreMockRecorder) ListAdminUsers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListAdminUsers", reflect.TypeOf((*MockStore)(nil).ListAdminUsers), arg0, arg1)
}

// ListPermissions mocks base method.
func (m *MockStore) ListPermissions(arg0 context.Context, arg1 db.ListPermissionsParams) ([]db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPermissions", arg0, arg1)
	ret0, _ := ret[0].([]db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPermissions indicates an expected call of ListPermissions.
func (mr *MockStoreMockRecorder) ListPermissions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPermissions", reflect.TypeOf((*MockStore)(nil).ListPermissions), arg0, arg1)
}

// ListPermissionsForAdminUser mocks base method.
func (m *MockStore) ListPermissionsForAdminUser(arg0 context.Context, arg1 int64) ([]db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPermissionsForAdminUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPermissionsForAdminUser indicates an expected call of ListPermissionsForAdminUser.
func (mr *MockStoreMockRecorder) ListPermissionsForAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPermissionsForAdminUser", reflect.TypeOf((*MockStore)(nil).ListPermissionsForAdminUser), arg0, arg1)
}

// ListPermissionsForRole mocks base method.
func (m *MockStore) ListPermissionsForRole(arg0 context.Context, arg1 int64) ([]db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPermissionsForRole", arg0, arg1)
	ret0, _ := ret[0].([]db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPermissionsForRole indicates an expected call of ListPermissionsForRole.
func (mr *MockStoreMockRecorder) ListPermissionsForRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPermissionsForRole", reflect.TypeOf((*MockStore)(nil).ListPermissionsForRole), arg0, arg1)
}

// ListRolePermissionByPermissionId mocks base method.
func (m *MockStore) ListRolePermissionByPermissionId(arg0 context.Context, arg1 pgtype.Int4) ([]db.RolePermission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRolePermissionByPermissionId", arg0, arg1)
	ret0, _ := ret[0].([]db.RolePermission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRolePermissionByPermissionId indicates an expected call of ListRolePermissionByPermissionId.
func (mr *MockStoreMockRecorder) ListRolePermissionByPermissionId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRolePermissionByPermissionId", reflect.TypeOf((*MockStore)(nil).ListRolePermissionByPermissionId), arg0, arg1)
}

// ListRolePermissionByRoleId mocks base method.
func (m *MockStore) ListRolePermissionByRoleId(arg0 context.Context, arg1 pgtype.Int4) ([]db.RolePermission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRolePermissionByRoleId", arg0, arg1)
	ret0, _ := ret[0].([]db.RolePermission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRolePermissionByRoleId indicates an expected call of ListRolePermissionByRoleId.
func (mr *MockStoreMockRecorder) ListRolePermissionByRoleId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRolePermissionByRoleId", reflect.TypeOf((*MockStore)(nil).ListRolePermissionByRoleId), arg0, arg1)
}

// ListRoles mocks base method.
func (m *MockStore) ListRoles(arg0 context.Context, arg1 db.ListRolesParams) ([]db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRoles", arg0, arg1)
	ret0, _ := ret[0].([]db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRoles indicates an expected call of ListRoles.
func (mr *MockStoreMockRecorder) ListRoles(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRoles", reflect.TypeOf((*MockStore)(nil).ListRoles), arg0, arg1)
}

// ListRolesForAdminUser mocks base method.
func (m *MockStore) ListRolesForAdminUser(arg0 context.Context, arg1 int64) ([]db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRolesForAdminUser", arg0, arg1)
	ret0, _ := ret[0].([]db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRolesForAdminUser indicates an expected call of ListRolesForAdminUser.
func (mr *MockStoreMockRecorder) ListRolesForAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRolesForAdminUser", reflect.TypeOf((*MockStore)(nil).ListRolesForAdminUser), arg0, arg1)
}

// UpdateAdminUser mocks base method.
func (m *MockStore) UpdateAdminUser(arg0 context.Context, arg1 db.UpdateAdminUserParams) (db.AdminUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdminUser", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAdminUser indicates an expected call of UpdateAdminUser.
func (mr *MockStoreMockRecorder) UpdateAdminUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdminUser", reflect.TypeOf((*MockStore)(nil).UpdateAdminUser), arg0, arg1)
}

// UpdateAdminUserTx mocks base method.
func (m *MockStore) UpdateAdminUserTx(arg0 context.Context, arg1 db.UpdateAdminUserTxParams) (db.AdminUserTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAdminUserTx", arg0, arg1)
	ret0, _ := ret[0].(db.AdminUserTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAdminUserTx indicates an expected call of UpdateAdminUserTx.
func (mr *MockStoreMockRecorder) UpdateAdminUserTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAdminUserTx", reflect.TypeOf((*MockStore)(nil).UpdateAdminUserTx), arg0, arg1)
}

// UpdatePermission mocks base method.
func (m *MockStore) UpdatePermission(arg0 context.Context, arg1 db.UpdatePermissionParams) (db.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePermission", arg0, arg1)
	ret0, _ := ret[0].(db.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdatePermission indicates an expected call of UpdatePermission.
func (mr *MockStoreMockRecorder) UpdatePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePermission", reflect.TypeOf((*MockStore)(nil).UpdatePermission), arg0, arg1)
}

// UpdateRole mocks base method.
func (m *MockStore) UpdateRole(arg0 context.Context, arg1 db.UpdateRoleParams) (db.Role, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRole", arg0, arg1)
	ret0, _ := ret[0].(db.Role)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRole indicates an expected call of UpdateRole.
func (mr *MockStoreMockRecorder) UpdateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRole", reflect.TypeOf((*MockStore)(nil).UpdateRole), arg0, arg1)
}

// UpdateRoleTx mocks base method.
func (m *MockStore) UpdateRoleTx(arg0 context.Context, arg1 db.UpdateRoleTxParams) (db.RoleTxResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRoleTx", arg0, arg1)
	ret0, _ := ret[0].(db.RoleTxResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRoleTx indicates an expected call of UpdateRoleTx.
func (mr *MockStoreMockRecorder) UpdateRoleTx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRoleTx", reflect.TypeOf((*MockStore)(nil).UpdateRoleTx), arg0, arg1)
}
