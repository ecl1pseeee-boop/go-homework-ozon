package models

type Cart struct {
	UserID    UserID
	CartItems []CartItem
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type CartID int64
