package dto

type AddProductRequestDTO struct {
	UserId int64  `json:"user_id"`
	SkuId  int64  `json:"sku_id"`
	Count  uint32 `json:"count"`
}
