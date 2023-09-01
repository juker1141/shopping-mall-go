// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Querier interface {
	CreateAdminUser(ctx context.Context, arg CreateAdminUserParams) (AdminUser, error)
	CreateAdminUserRole(ctx context.Context, arg CreateAdminUserRoleParams) (AdminUserRole, error)
	CreateCart(ctx context.Context, arg CreateCartParams) (Cart, error)
	CreateCategory(ctx context.Context, name string) (Category, error)
	CreateCoupon(ctx context.Context, arg CreateCouponParams) (Coupon, error)
	CreatePermission(ctx context.Context, name string) (Permission, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateRole(ctx context.Context, name string) (Role, error)
	CreateRolePermission(ctx context.Context, arg CreateRolePermissionParams) (RolePermission, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAdminUser(ctx context.Context, id int64) error
	DeleteAdminUserRoleByAdminUserId(ctx context.Context, adminUserID pgtype.Int4) error
	DeleteAdminUserRoleByRoleId(ctx context.Context, roleID pgtype.Int4) error
	DeleteCart(ctx context.Context, id int64) error
	DeleteCategory(ctx context.Context, id int64) error
	DeleteCoupon(ctx context.Context, id int64) error
	DeletePermission(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteRole(ctx context.Context, id int64) error
	DeleteRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) error
	DeleteRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) error
	DeleteUser(ctx context.Context, id int64) error
	GetAdminUser(ctx context.Context, id int64) (AdminUser, error)
	GetAdminUserByAccount(ctx context.Context, account string) (AdminUser, error)
	GetAdminUserRole(ctx context.Context, arg GetAdminUserRoleParams) (AdminUserRole, error)
	GetAdminUsersCount(ctx context.Context) (int64, error)
	GetCart(ctx context.Context, id int64) (Cart, error)
	GetCartByOwner(ctx context.Context, owner pgtype.Text) (Cart, error)
	GetCategoriesCount(ctx context.Context) (int64, error)
	GetCategory(ctx context.Context, id int64) (Category, error)
	GetCoupon(ctx context.Context, id int64) (Coupon, error)
	GetPermission(ctx context.Context, id int64) (Permission, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetRolePermission(ctx context.Context, arg GetRolePermissionParams) (RolePermission, error)
	GetRolesCount(ctx context.Context) (int64, error)
	GetSesstion(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUsersCount(ctx context.Context) (int64, error)
	ListAdminUserRoleByAdminUserId(ctx context.Context, adminUserID pgtype.Int4) ([]AdminUserRole, error)
	ListAdminUserRoleByRoleId(ctx context.Context, roleID pgtype.Int4) ([]AdminUserRole, error)
	ListAdminUsers(ctx context.Context, arg ListAdminUsersParams) ([]AdminUser, error)
	ListCategories(ctx context.Context) ([]Category, error)
	ListCoupons(ctx context.Context, arg ListCouponsParams) ([]Coupon, error)
	ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]Permission, error)
	ListPermissionsForAdminUser(ctx context.Context, id int64) ([]Permission, error)
	ListPermissionsForRole(ctx context.Context, id int64) ([]Permission, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) ([]RolePermission, error)
	ListRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) ([]RolePermission, error)
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListRolesForAdminUser(ctx context.Context, id int64) ([]Role, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAdminUser(ctx context.Context, arg UpdateAdminUserParams) (AdminUser, error)
	// -- name: ListProductsForCart :many
	// SELECT DISTINCT product.id, product.name, product.created_at
	// FROM carts as c
	// JOIN cart_products AS cp ON c.id = cp.role_id
	// JOIN products AS product ON cp.product_id = product.id
	// WHERE c.id = $1;
	UpdateCart(ctx context.Context, arg UpdateCartParams) (Cart, error)
	UpdateCategory(ctx context.Context, arg UpdateCategoryParams) (Category, error)
	UpdateCoupon(ctx context.Context, arg UpdateCouponParams) (Coupon, error)
	UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
