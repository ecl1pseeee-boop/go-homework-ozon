package services

import (
	"context"
	"testing"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/repository"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCartService_AddProduct_Success(t *testing.T) {
	t.Parallel()
	defaultProducts := map[int64]fakeProduct{
		1: {
			Name:  "fakeProduct",
			Price: 54,
		},
	}

	tests := []struct {
		testName         string
		ctx              context.Context
		initProducts     map[int64]fakeProduct
		newProduct       dto.AddProductDTO
		expectedUserCart []models.CartItem
	}{
		{
			testName:     "Базовое добавление продукта",
			ctx:          context.Background(),
			initProducts: defaultProducts,
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  1,
				Count:  3,
			},
			expectedUserCart: []models.CartItem{
				{
					Sku:   1,
					Count: 3,
					Name:  "fakeProduct",
					Price: 54,
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewCartRepository()
			svc := NewCartService(
				repo,
				newFakeProductClient(tc.initProducts),
			)

			addProductErr := svc.AddProduct(tc.ctx, tc.newProduct)
			require.NoError(t, addProductErr)

			cart, err := repo.GetUserCart(models.UserID(tc.newProduct.UserId))
			require.NoError(t, err)
			require.Equal(t, tc.expectedUserCart, cart.CartItems)

		})
	}
}

func TestCartService_AddProduct_Fail(t *testing.T) {
	t.Parallel()

	defaultProducts := map[int64]fakeProduct{
		1: {
			Name:  "fakeProduct",
			Price: 54,
		},
	}
	tests := []struct {
		testName    string
		ctx         context.Context
		clientErr   error
		initProduct map[int64]fakeProduct
		newProduct  dto.AddProductDTO
		wantErr     error
	}{
		{
			testName:  "Product service не доступен.",
			ctx:       context.Background(),
			clientErr: status.Error(codes.Unavailable, "down"),
			wantErr:   carterrors.ErrProductUnavailable,
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  1,
				Count:  3,
			},
		},
		{
			testName:    "Не существующий product.",
			ctx:         context.Background(),
			initProduct: defaultProducts,
			wantErr:     carterrors.ErrProductNotFound,
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  2,
				Count:  4,
			},
		},
		{
			testName: "Неправильный формат данных (UserID меньше чем 0).",
			ctx:      context.Background(),
			newProduct: dto.AddProductDTO{
				UserId: -1,
				SkuId:  1,
				Count:  35,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName: "Неправильный формат данных (UserID = 0).",
			ctx:      context.Background(),
			newProduct: dto.AddProductDTO{
				UserId: 0,
				SkuId:  1,
				Count:  35,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName: "Неправильный формат данных (SkuID меньше чем 0).",
			ctx:      context.Background(),
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  -1,
				Count:  35,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName: "Неправильный формат данных (SkuID = 0).",
			ctx:      context.Background(),
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  0,
				Count:  35,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName: "Неправильный формат данных. (Count = 0).",
			ctx:      context.Background(),
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  1,
				Count:  0,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName:    "Не валидные данные внутри product service.",
			ctx:         context.Background(),
			clientErr:   status.Error(codes.InvalidArgument, "invalid argument"),
			initProduct: defaultProducts,
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  1,
				Count:  35,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName:  "Не правильный адрес Product service.",
			ctx:       context.Background(),
			clientErr: status.Error(codes.Unknown, "down"),
			wantErr:   carterrors.ErrProductUnavailable,
			newProduct: dto.AddProductDTO{
				UserId: 1,
				SkuId:  1,
				Count:  3,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			repo := repository.NewCartRepository()
			fake := newFakeProductClient(tc.initProduct)
			fake.err = tc.clientErr
			svc := NewCartService(repo, fake)

			addProductErr := svc.AddProduct(tc.ctx, tc.newProduct)

			require.ErrorIs(t, addProductErr, tc.wantErr)
		})
	}
}

func TestCartService_AddProduct_Again_Success(t *testing.T) {
	t.Parallel()
	repo := repository.NewCartRepository()
	svc := NewCartService(
		repo,
		newFakeProductClient(map[int64]fakeProduct{
			1: {
				Name:  "fakeProduct",
				Price: 54,
			},
		}),
	)

	repoAddErr := repo.AddProduct(1, models.CartItem{
		Sku:   1,
		Count: 1,
		Name:  "fakeProduct",
		Price: 54,
	})
	require.NoError(t, repoAddErr)

	svcAddErr := svc.AddProduct(context.Background(), dto.AddProductDTO{
		UserId: 1,
		SkuId:  1,
		Count:  1,
	})
	require.NoError(t, svcAddErr)

	updatedUserCartResponse, repoErr := repo.GetUserCart(1)
	require.NoError(t, repoErr)
	require.Equal(t, 1, len(updatedUserCartResponse.CartItems))
	require.Equal(t, uint32(2), updatedUserCartResponse.CartItems[0].Count)
}

func TestCartService_GetUserCart_Success(t *testing.T) {
	t.Parallel()

	products := map[int64]fakeProduct{
		2: {
			Name:  "fakeProduct2",
			Price: 54,
		},
		1: {
			Name:  "fakeProduct1",
			Price: 54,
		},
	}

	repo := repository.NewCartRepository()
	svc := NewCartService(
		repo,
		newFakeProductClient(products),
	)

	userCartItem2 := models.CartItem{
		Sku:   2,
		Count: 3,
		Name:  "fakeProduct2",
		Price: uint32(54),
	}
	addProductErr := repo.AddProduct(1, userCartItem2)
	require.NoError(t, addProductErr)

	userCartItem1 := models.CartItem{
		Sku:   1,
		Count: 3,
		Name:  "fakeProduct1",
		Price: uint32(54),
	}
	addProductErr2 := repo.AddProduct(1, userCartItem1)
	require.NoError(t, addProductErr2)

	userCart, cartErr := svc.GetUserCart(context.Background(), dto.GetUserCartDTO{
		UserId: 1,
	})

	expectedUserCartItems := []models.CartItem{
		{
			Sku:   1,
			Count: 3,
			Name:  "fakeProduct1",
			Price: 54,
		},
		{
			Sku:   2,
			Count: 3,
			Name:  "fakeProduct2",
			Price: 54,
		},
	}

	require.NoError(t, cartErr)
	require.Equal(t, expectedUserCartItems, userCart.Items)

}

func TestCartService_GetUserCart_Fail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName    string
		wantErr     error
		getUserCart dto.GetUserCartDTO
	}{
		{
			testName: "Не валидный id пользователя.",
			getUserCart: dto.GetUserCartDTO{
				UserId: -1,
			},
			wantErr: carterrors.ErrInvalidData,
		},
		{
			testName: "Пользовательская корзина не найдена.",
			getUserCart: dto.GetUserCartDTO{
				UserId: 4,
			},
			wantErr: carterrors.ErrCartNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()

			repo := repository.NewCartRepository()
			svc := NewCartService(
				repo,
				newFakeProductClient(make(map[int64]fakeProduct)),
			)

			userCartResponse, cartErr := svc.GetUserCart(context.Background(), tc.getUserCart)
			require.ErrorIs(t, cartErr, tc.wantErr)
			require.Nil(t, userCartResponse.Items)
			require.Zero(t, userCartResponse.TotalPrice)
		})
	}

}

func TestCartService_SolveTotalPrice(t *testing.T) {
	t.Parallel()

	userCart := models.Cart{
		UserID: 1,
		CartItems: []models.CartItem{
			{
				Sku:   1,
				Name:  "fakeProduct",
				Count: 2,
				Price: uint32(54),
			},
			{
				Sku:   2,
				Count: 2,
				Name:  "fakeProduct2",
				Price: uint32(50),
			},
			{
				Sku:   3,
				Count: 2,
				Name:  "fakeProduct3",
				Price: uint32(100),
			},
		},
	}

	expectedTotalPrice := uint32(408)

	repo := repository.NewCartRepository()
	svc := NewCartService(
		repo,
		newFakeProductClient(make(map[int64]fakeProduct)),
	)

	for _, cartItem := range userCart.CartItems {
		addProductErr := repo.AddProduct(1, cartItem)
		require.NoError(t, addProductErr)
	}

	userCartResponse, cartErr := svc.GetUserCart(context.Background(), dto.GetUserCartDTO{
		UserId: 1,
	})

	require.NoError(t, cartErr)
	require.Equal(t, expectedTotalPrice, userCartResponse.TotalPrice)
}
