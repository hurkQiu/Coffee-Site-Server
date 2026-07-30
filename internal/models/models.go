package models

import "time"

type Permission string

const (
	PermissionMember Permission = "member"
	PermissionAdmin  Permission = "admin"
)

type User struct {
	ID           int64      `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	Permission   Permission `json:"permission"`
	Verified     bool       `json:"verified"`
	CreatedAt    time.Time  `json:"createdAt"`
}

type VerificationPurpose string

const (
	PurposeRegister VerificationPurpose = "register"
	PurposeReset    VerificationPurpose = "reset"
)

type VerificationCode struct {
	ID        int64
	Email     string
	Code      string
	Purpose   VerificationPurpose
	ExpiresAt time.Time
	Consumed  bool
	CreatedAt time.Time
}

type NamedOption struct {
	Name   string `json:"name"`
	Hidden bool   `json:"hidden"`
}

type CoffeeBean struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	Image   string `json:"image"`
	Stock   int    `json:"stock"`
	Roast   string `json:"roast"`
	Process string `json:"process"`
}

type Utensil struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Image    string `json:"image"`
	Stock    int    `json:"stock"`
	Category string `json:"category"`
}

type CartItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Image    string `json:"image"`
	Quantity int    `json:"quantity"`
}

type Order struct {
	ID          int64       `json:"id"`
	OrderNumber string      `json:"orderNumber"`
	UserID      int64       `json:"-"`
	TotalPrice  int         `json:"totalPrice"`
	CreatedAt   time.Time   `json:"createdAt"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	ID       int64  `json:"-"`
	OrderID  int64  `json:"-"`
	Name     string `json:"name"`
	Image    string `json:"image"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}

type ContactMessage struct {
	ID        int64     `json:"id"`
	UserID    *int64    `json:"userId,omitempty"`
	OrderID   string    `json:"orderId,omitempty"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
