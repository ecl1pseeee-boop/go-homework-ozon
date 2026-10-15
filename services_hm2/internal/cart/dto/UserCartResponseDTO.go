package dto

import "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"

type UserCartResponseDTO struct {
	Items      []models.CartItem `json:"items,omitempty"`
	TotalPrice uint32            `json:"total_price,omitempty"`
}
