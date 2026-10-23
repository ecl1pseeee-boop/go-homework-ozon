package services

import (
	"cmp"
	"context"
	"fmt"
	"slices"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartService struct {
	repository    CartRepository
	productClient ProductGRPCClient
}

type ProductGRPCClient interface {
	GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error)
}

type CartRepository interface {
	AddProduct(userId models.UserID, cartItem models.CartItem) error
	GetUserCart(userId models.UserID) (models.Cart, error)
}

func NewCartService(repository CartRepository, productClient ProductGRPCClient) *CartService {
	return &CartService{
		repository:    repository,
		productClient: productClient,
	}
}

func (s *CartService) AddProduct(ctx context.Context, addProductDto dto.AddProductDTO) error {
	if addProductDto.UserId <= 0 || addProductDto.SkuId <= 0 || addProductDto.Count <= 0 {
		return carterrors.ErrInvalidData
	}
	getProductRequest := &pb.GetProductRequest{
		Sku: addProductDto.SkuId,
	}
	productResponse, err := s.productClient.GetProduct(ctx, getProductRequest)
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return carterrors.ErrProductUnavailable
		}

		switch st.Code() {
		case codes.NotFound:
			return carterrors.ErrProductNotFound
		case codes.InvalidArgument:
			return carterrors.ErrInvalidData
		default:
			return carterrors.ErrProductUnavailable
		}
	}

	cartItem := models.CartItem{
		Sku:   productResponse.Sku,
		Name:  productResponse.Name,
		Price: uint32(productResponse.Price),
		Count: addProductDto.Count,
	}

	addErr := s.repository.AddProduct(models.UserID(addProductDto.UserId), cartItem)
	if addErr != nil {
		return fmt.Errorf("ошибка добавления продукта: %w", addErr)
	}

	return nil
}

func (s *CartService) GetUserCart(ctx context.Context, userCartDto dto.GetUserCartDTO) (dto.UserCartResponseDTO, error) {
	if userCartDto.UserId <= 0 {
		return dto.UserCartResponseDTO{}, carterrors.ErrInvalidData
	}
	userCart, err := s.repository.GetUserCart(models.UserID(userCartDto.UserId))
	if err != nil {
		return dto.UserCartResponseDTO{}, fmt.Errorf("ошибка получения корзины пользователя: %w", err)
	}

	cartItems := slices.Clone(userCart.CartItems)

	slices.SortFunc(cartItems, func(a, b models.CartItem) int {
		return cmp.Compare(a.Sku, b.Sku)
	})

	totalPrice, err := s.solveTotalPrice(&userCart)
	if err != nil {
		return dto.UserCartResponseDTO{}, fmt.Errorf("ошибка вычисления суммы корзины пользователя: %w", err)
	}

	response := dto.UserCartResponseDTO{
		Items:      cartItems,
		TotalPrice: totalPrice,
	}
	return response, nil
}

func (s *CartService) solveTotalPrice(userCart *models.Cart) (uint32, error) {
	var price uint32
	for _, cart := range userCart.CartItems {
		price += cart.Price * cart.Count
	}
	return price, nil
}
