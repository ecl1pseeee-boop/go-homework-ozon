package services

import (
	"cmp"
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"
	"golang.org/x/sync/errgroup"
)

const timeout = time.Second

type CartService struct {
	repository    CartRepository
	productClient ProductClient
}

type CartRepository interface {
	AddProduct(userId models.UserID, cartItem models.CartItem) error
	GetUserCart(userId models.UserID) (models.Cart, error)
}

type ProductClient interface {
	GetProduct(ctx context.Context, skuId int64) (*dto.GetProductResponseDTO, error)
}

func NewCartService(repository CartRepository, productClient ProductClient) *CartService {
	return &CartService{
		repository:    repository,
		productClient: productClient,
	}
}

func (s *CartService) AddProduct(ctx context.Context, addProductDto dto.AddProductRequestDTO) error {
	if addProductDto.UserId <= 0 || addProductDto.SkuId <= 0 || addProductDto.Count <= 0 {
		return carterrors.ErrInvalidData
	}

	newCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	productResponse, err := s.productClient.GetProduct(newCtx, addProductDto.SkuId)
	if err != nil {
		return err
	}

	cartItem := models.CartItem{
		Sku:   productResponse.Sku,
		Count: addProductDto.Count,
	}

	addErr := s.repository.AddProduct(models.UserID(addProductDto.UserId), cartItem)
	if addErr != nil {
		return fmt.Errorf("ошибка добавления продукта: %w", addErr)
	}

	return nil
}

func (s *CartService) GetUserCart(ctx context.Context, userCartDto dto.GetUserCartRequestDTO) (dto.UserCartResponseDTO, error) {
	if userCartDto.UserId <= 0 {
		return dto.UserCartResponseDTO{}, carterrors.ErrInvalidData
	}

	userCart, err := s.repository.GetUserCart(models.UserID(userCartDto.UserId))
	if err != nil {
		return dto.UserCartResponseDTO{}, fmt.Errorf("ошибка получения корзины пользователя: %w", err)
	}

	cartItems := userCart.CartItems
	if len(cartItems) == 0 {
		return dto.UserCartResponseDTO{}, carterrors.ErrCartEmpty
	}

	results := make([]dto.CartItemInfo, len(cartItems))

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	g, groupCtx := errgroup.WithContext(timeoutCtx)

	for i, cartItem := range cartItems {
		g.Go(func() error {
			productResponse, respErr := s.productClient.GetProduct(groupCtx, cartItem.Sku)
			if respErr == nil {
				results[i] = dto.CartItemInfo{
					Sku:   cartItem.Sku,
					Price: uint32(productResponse.Price),
					Name:  productResponse.Name,
					Count: cartItem.Count,
				}
			}
			return respErr
		})
	}

	if gerr := g.Wait(); gerr != nil {
		return dto.UserCartResponseDTO{}, gerr
	}

	totalPrice, solveErr := s.solveTotalPrice(results)
	if solveErr != nil {
		return dto.UserCartResponseDTO{}, fmt.Errorf("ошибка вычисления суммы корзины пользователя: %w", solveErr)
	}

	slices.SortFunc(results, func(a, b dto.CartItemInfo) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	return dto.UserCartResponseDTO{
		Items:      results,
		TotalPrice: totalPrice,
	}, nil
}

func (s *CartService) solveTotalPrice(products []dto.CartItemInfo) (uint32, error) {
	var price uint32
	for _, cart := range products {
		price += cart.Price * cart.Count
	}
	return price, nil
}
