package models

type Cart struct {
	UserID    UserID
	CartItems []CartItem
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Count uint32 `json:"count"`
}
