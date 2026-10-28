package dto

type UserCartResponseDTO struct {
	Items      []CartItemInfo `json:"items,omitempty"`
	TotalPrice uint32         `json:"total_price,omitempty"`
}

type CartItemInfo struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Price uint32 `json:"price"`
	Count uint32 `json:"count"`
}
