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
	CreateCart(ctx context.Context, arg CreateCartParams) (Cart, error)
	CreateCartCoupon(ctx context.Context, arg CreateCartCouponParams) (CartCoupon, error)
	CreateCartProduct(ctx context.Context, arg CreateCartProductParams) (CartProduct, error)
	CreateCoupon(ctx context.Context, arg CreateCouponParams) (Coupon, error)
	CreateGender(ctx context.Context, name string) (Gender, error)
	CreateOrder(ctx context.Context, arg CreateOrderParams) (Order, error)
	CreateOrderCoupon(ctx context.Context, arg CreateOrderCouponParams) (OrderCoupon, error)
	CreateOrderProduct(ctx context.Context, arg CreateOrderProductParams) (OrderProduct, error)
	CreateOrderStatus(ctx context.Context, arg CreateOrderStatusParams) (OrderStatus, error)
	CreateOrderUser(ctx context.Context, arg CreateOrderUserParams) (OrderUser, error)
	CreatePayMethod(ctx context.Context, name string) (PayMethod, error)
	CreatePermission(ctx context.Context, name string) (Permission, error)
	CreateProduct(ctx context.Context, arg CreateProductParams) (Product, error)
	CreateRole(ctx context.Context, name string) (Role, error)
	CreateRolePermission(ctx context.Context, arg CreateRolePermissionParams) (RolePermission, error)
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	DeleteAdminUser(ctx context.Context, id int64) error
	DeleteCart(ctx context.Context, id int64) error
	DeleteCartCouponByCartId(ctx context.Context, cartID pgtype.Int4) error
	DeleteCartCouponByCouponId(ctx context.Context, couponID pgtype.Int4) error
	DeleteCartProductByCartId(ctx context.Context, cartID pgtype.Int4) error
	DeleteCartProductByProductId(ctx context.Context, productID pgtype.Int4) error
	DeleteCoupon(ctx context.Context, id int64) error
	DeleteOrder(ctx context.Context, id int64) error
	DeleteOrderCouponByCouponId(ctx context.Context, couponID pgtype.Int4) error
	DeleteOrderCouponByOrderId(ctx context.Context, orderID pgtype.Int4) error
	DeleteOrderProductByOrderId(ctx context.Context, orderID pgtype.Int4) error
	DeleteOrderProductByProductId(ctx context.Context, productID pgtype.Int4) error
	DeleteOrderUserByOrderId(ctx context.Context, orderID pgtype.Int4) error
	DeleteOrderUserByUserId(ctx context.Context, userID pgtype.Int4) error
	DeletePermission(ctx context.Context, id int64) error
	DeleteProduct(ctx context.Context, id int64) error
	DeleteRole(ctx context.Context, id int64) error
	DeleteRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) error
	DeleteRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) error
	DeleteUser(ctx context.Context, id int64) error
	GetAdminUser(ctx context.Context, id int64) (AdminUser, error)
	GetAdminUserByAccount(ctx context.Context, account string) (AdminUser, error)
	GetAdminUsersCount(ctx context.Context) (int64, error)
	GetCart(ctx context.Context, id int64) (Cart, error)
	GetCartByOwner(ctx context.Context, owner pgtype.Text) (Cart, error)
	GetCartCoupon(ctx context.Context, arg GetCartCouponParams) (CartCoupon, error)
	GetCartProduct(ctx context.Context, arg GetCartProductParams) (CartProduct, error)
	GetCoupon(ctx context.Context, id int64) (Coupon, error)
	GetCouponByCode(ctx context.Context, code string) (Coupon, error)
	GetCouponsCount(ctx context.Context) (int64, error)
	GetGender(ctx context.Context, id int64) (Gender, error)
	GetOrder(ctx context.Context, id int64) (Order, error)
	GetOrderCoupon(ctx context.Context, arg GetOrderCouponParams) (OrderCoupon, error)
	GetOrderProduct(ctx context.Context, arg GetOrderProductParams) (OrderProduct, error)
	GetOrderStatus(ctx context.Context, id int64) (OrderStatus, error)
	GetOrderStatusCount(ctx context.Context) (int64, error)
	GetOrderUser(ctx context.Context, arg GetOrderUserParams) (OrderUser, error)
	GetOrdersCount(ctx context.Context) (int64, error)
	GetPayMethod(ctx context.Context, id int64) (PayMethod, error)
	GetPayMethodCount(ctx context.Context) (int64, error)
	GetPermission(ctx context.Context, id int64) (Permission, error)
	GetProduct(ctx context.Context, id int64) (Product, error)
	GetProductsCount(ctx context.Context) (int64, error)
	GetRole(ctx context.Context, id int64) (Role, error)
	GetRolePermission(ctx context.Context, arg GetRolePermissionParams) (RolePermission, error)
	GetRolesCount(ctx context.Context) (int64, error)
	GetSesstion(ctx context.Context, id uuid.UUID) (Session, error)
	GetUser(ctx context.Context, id int64) (User, error)
	GetUserByAccount(ctx context.Context, account string) (User, error)
	GetUsersCount(ctx context.Context) (int64, error)
	ListAdminUsers(ctx context.Context, arg ListAdminUsersParams) ([]AdminUser, error)
	ListCartCouponByCartId(ctx context.Context, cartID pgtype.Int4) ([]CartCoupon, error)
	ListCartCouponByCouponId(ctx context.Context, couponID pgtype.Int4) ([]CartCoupon, error)
	ListCartProductByCartId(ctx context.Context, cartID pgtype.Int4) ([]CartProduct, error)
	ListCartProductByProductId(ctx context.Context, productID pgtype.Int4) ([]CartProduct, error)
	ListCoupons(ctx context.Context, arg ListCouponsParams) ([]Coupon, error)
	ListGenders(ctx context.Context) ([]Gender, error)
	ListOrderCouponByCouponId(ctx context.Context, couponID pgtype.Int4) ([]OrderCoupon, error)
	ListOrderCouponByOrderId(ctx context.Context, orderID pgtype.Int4) ([]OrderCoupon, error)
	ListOrderProductByOrderId(ctx context.Context, orderID pgtype.Int4) ([]OrderProduct, error)
	ListOrderProductByProductId(ctx context.Context, productID pgtype.Int4) ([]OrderProduct, error)
	ListOrderStatus(ctx context.Context) ([]OrderStatus, error)
	ListOrderUserByOrderId(ctx context.Context, orderID pgtype.Int4) ([]OrderUser, error)
	ListOrderUserByUserId(ctx context.Context, userID pgtype.Int4) ([]OrderUser, error)
	// -- name: GetOrderByCode :one
	// SELECT * FROM orders
	// WHERE code = $1 LIMIT 1;
	ListOrders(ctx context.Context, arg ListOrdersParams) ([]Order, error)
	ListPayMethod(ctx context.Context) ([]PayMethod, error)
	ListPermissions(ctx context.Context, arg ListPermissionsParams) ([]Permission, error)
	ListPermissionsForAdminUser(ctx context.Context, id int64) ([]Permission, error)
	ListPermissionsForRole(ctx context.Context, id int64) ([]Permission, error)
	ListPermissionsIDByAccount(ctx context.Context, account string) ([]int64, error)
	ListProducts(ctx context.Context, arg ListProductsParams) ([]Product, error)
	ListRolePermissionByPermissionId(ctx context.Context, permissionID pgtype.Int4) ([]RolePermission, error)
	ListRolePermissionByRoleId(ctx context.Context, roleID pgtype.Int4) ([]RolePermission, error)
	ListRoles(ctx context.Context, arg ListRolesParams) ([]Role, error)
	ListRolesOption(ctx context.Context) ([]Role, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	UpdateAdminUser(ctx context.Context, arg UpdateAdminUserParams) (AdminUser, error)
	// -- name: ListProductsForCart :many
	// SELECT DISTINCT product.id, product.name, product.created_at
	// FROM carts as c
	// JOIN cart_products AS cp ON c.id = cp.role_id
	// JOIN products AS product ON cp.product_id = product.id
	// WHERE c.id = $1;
	UpdateCart(ctx context.Context, arg UpdateCartParams) (Cart, error)
	UpdateCartProduct(ctx context.Context, arg UpdateCartProductParams) (CartProduct, error)
	UpdateCoupon(ctx context.Context, arg UpdateCouponParams) (Coupon, error)
	UpdateOrder(ctx context.Context, arg UpdateOrderParams) (Order, error)
	UpdateOrderProduct(ctx context.Context, arg UpdateOrderProductParams) (OrderProduct, error)
	UpdatePermission(ctx context.Context, arg UpdatePermissionParams) (Permission, error)
	UpdateProduct(ctx context.Context, arg UpdateProductParams) (Product, error)
	UpdateRole(ctx context.Context, arg UpdateRoleParams) (Role, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
