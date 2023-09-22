// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AdminUser struct {
	ID             int64       `json:"id"`
	Account        string      `json:"account"`
	FullName       string      `json:"full_name"`
	HashedPassword string      `json:"hashed_password"`
	RoleID         pgtype.Int4 `json:"role_id"`
	// must be either 0 or 1
	Status            int32     `json:"status"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

type Cart struct {
	ID    int64       `json:"id"`
	Owner pgtype.Text `json:"owner"`
	// must be positive
	TotalPrice int32 `json:"total_price"`
	// must be positive
	FinalPrice int32     `json:"final_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type CartCoupon struct {
	CartID   pgtype.Int4 `json:"cart_id"`
	CouponID pgtype.Int4 `json:"coupon_id"`
}

type CartProduct struct {
	CartID    pgtype.Int4 `json:"cart_id"`
	ProductID pgtype.Int4 `json:"product_id"`
	Num       int32       `json:"num"`
}

type Coupon struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Code      string    `json:"code"`
	Percent   int32     `json:"percent"`
	CreatedBy string    `json:"created_by"`
	StartAt   time.Time `json:"start_at"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Gender struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Order struct {
	ID        int64       `json:"id"`
	IsPaid    pgtype.Bool `json:"is_paid"`
	StatusID  pgtype.Int4 `json:"status_id"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type OrderCoupon struct {
	OrderID  pgtype.Int4 `json:"order_id"`
	CouponID pgtype.Int4 `json:"coupon_id"`
}

type OrderProduct struct {
	OrderID   pgtype.Int4 `json:"order_id"`
	ProductID pgtype.Int4 `json:"product_id"`
	Num       int32       `json:"num"`
}

type OrderStatus struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type OrderUser struct {
	OrderID pgtype.Int4 `json:"order_id"`
	UserID  pgtype.Int4 `json:"user_id"`
}

type Permission struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	OriginPrice int32  `json:"origin_price"`
	Price       int32  `json:"price"`
	Unit        string `json:"unit"`
	Description string `json:"description"`
	Content     string `json:"content"`
	// must be either 0 or 1
	Status    int32     `json:"status"`
	ImageUrl  string    `json:"image_url"`
	ImagesUrl []string  `json:"images_url"`
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type RolePermission struct {
	RoleID       pgtype.Int4 `json:"role_id"`
	PermissionID pgtype.Int4 `json:"permission_id"`
}

type Session struct {
	ID           uuid.UUID `json:"id"`
	Account      string    `json:"account"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type User struct {
	ID              int64       `json:"id"`
	Account         string      `json:"account"`
	Email           string      `json:"email"`
	FullName        string      `json:"full_name"`
	GenderID        pgtype.Int4 `json:"gender_id"`
	Cellphone       string      `json:"cellphone"`
	Address         string      `json:"address"`
	ShippingAddress string      `json:"shipping_address"`
	PostCode        string      `json:"post_code"`
	HashedPassword  string      `json:"hashed_password"`
	// must be either 0 or 1
	Status            int32     `json:"status"`
	AvatarUrl         string    `json:"avatar_url"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}
