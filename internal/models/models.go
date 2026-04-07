package models

// Пользователь
type User struct {
	UserID     int    `json:"user_id"`
	FullName   string `json:"full_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Password   string `json:"password,omitempty"`
}

// Корзина
type Cart struct {
	CartID int `json:"cart_id"`
	UserID int `json:"user_id"`
}

// Товар
type Product struct {
	ProductID   int     `json:"product_id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
}

// Товар в корзине
type CartItem struct {
	CartItemID int `json:"cart_item_id"`
	CartID     int `json:"cart_id"`
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
}

// Расширенный товар в корзине (с данными товара)
type EnrichedCartItem struct {
	CartItemID  int     `json:"cart_item_id"`
	CartID      int     `json:"cart_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Category    string  `json:"category"`
	TotalPrice  float64 `json:"total_price"` // количество * цена
}

// Запросы
type AddToCartRequest struct {
	CartID    int `json:"cart_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	UserID    int `json:"userID"`
}

type UpdateCartCountRequest struct {
	Quantity int `json:"quantity"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupRequest struct {
	FullName   string `json:"full_name"`
	MiddleName string `json:"middle_name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
}

type LoginResponse struct {
	UserID     int    `json:"user_id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	CartID     int    `json:"cart_id"`
	MiddleName string `json:"middle_name"`
}
