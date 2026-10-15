package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
)

type CartHandler struct {
	cartService CartService
}

type CartService interface {
	AddProduct(ctx context.Context, dto dto.AddProductDTO) error
	GetUserCart(ctx context.Context, dto dto.GetUserCartDTO) (dto.UserCartResponseDTO, error)
}

func NewCartHandler(cartService CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

func (h *CartHandler) GetCartHandler(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil {
		http.Error(w, "Неправильный формат user_id", http.StatusBadRequest)
		return
	}

	userCartDTO := dto.GetUserCartDTO{
		UserId: userID,
	}
	userCart, serviceErr := h.cartService.GetUserCart(r.Context(), userCartDTO)

	if serviceErr != nil {
		switch {
		case errors.Is(serviceErr, carterrors.ErrCartNotFound):
			http.Error(w, serviceErr.Error(), http.StatusNotFound)
		case errors.Is(serviceErr, carterrors.ErrInvalidData):
			http.Error(w, serviceErr.Error(), http.StatusBadRequest)
		default:
			http.Error(w, "Ошибка получения доступа к сервису", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	encodeErr := json.NewEncoder(w).Encode(userCart)
	if encodeErr != nil {
		http.Error(w, encodeErr.Error(), http.StatusInternalServerError)
		return
	}
}

type AddProductRequest struct {
	Count uint32 `json:"count"`
}

func (h *CartHandler) AddProductHandler(w http.ResponseWriter, r *http.Request) {
	addProductDTO, err := h.parseAddProductRequest(r)
	if err != nil {
		http.Error(w, "Неправильный формат данных", http.StatusBadRequest)
		return
	}

	serviceErr := h.cartService.AddProduct(r.Context(), addProductDTO)
	if serviceErr != nil {
		switch {
		case errors.Is(serviceErr, carterrors.ErrInvalidData):
			http.Error(w, serviceErr.Error(), http.StatusBadRequest)
		case errors.Is(serviceErr, carterrors.ErrProductNotFound):
			http.Error(w, serviceErr.Error(), http.StatusPreconditionFailed)
		default:
			http.Error(w, "Ошибка получения доступа к сервису", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *CartHandler) parseAddProductRequest(r *http.Request) (dto.AddProductDTO, error) {
	userID, userParseErr := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if userParseErr != nil {
		return dto.AddProductDTO{}, userParseErr
	}

	skuID, skuParseErr := strconv.ParseInt(r.PathValue("sku_id"), 10, 64)
	if skuParseErr != nil {
		return dto.AddProductDTO{}, skuParseErr
	}

	var req AddProductRequest
	decodeErr := json.NewDecoder(r.Body).Decode(&req)
	if decodeErr != nil {
		return dto.AddProductDTO{}, decodeErr
	}

	return dto.AddProductDTO{
		UserId: userID,
		SkuId:  skuID,
		Count:  req.Count,
	}, nil
}
