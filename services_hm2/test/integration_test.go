//go:build api

package test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCProduct_Integration(t *testing.T) {
	t.Parallel()
	productClient := NewProductClient(t)
	t.Run("Тестирование получения существующего продукта.", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(func() {
			cancel()
		})

		resp, productErr := productClient.GetProduct(ctx, &pb.GetProductRequest{
			Sku: 1148162,
		})

		require.Equal(t, codes.OK, status.Code(productErr))
		require.NotZero(t, resp.Price)
		require.NotZero(t, resp.Sku)
		require.NotEmpty(t, resp.Name)

	})

	t.Run("Тестирование запроса на не существующий продукт.", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(func() {
			cancel()
		})
		_, productErr := productClient.GetProduct(ctx, &pb.GetProductRequest{
			Sku: 1,
		})

		require.Error(t, productErr)
		require.Equal(t, codes.NotFound, status.Code(productErr))
	})

	t.Run("Тестирование запроса с не валидным id.", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(func() {
			cancel()
		})
		_, productErr := productClient.GetProduct(ctx, &pb.GetProductRequest{
			Sku: -5,
		})

		require.Error(t, productErr)
		require.Equal(t, codes.InvalidArgument, status.Code(productErr))
	})
}

func TestHTTPCart_Integration(t *testing.T) {
	t.Parallel()

	cartClient := NewCartClient(t)
	t.Run("Добавление продукта с существующем SKU.", func(t *testing.T) {
		t.Parallel()

		jsonData := []byte(`{"count": 2}`)
		userId := time.Now().UnixNano()
		skuId := int64(1148162)
		resp, err := cartClient.AddProductToCart(t, jsonData, userId, skuId)

		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Добавление продукта с не существующем SKU.", func(t *testing.T) {
		t.Parallel()

		jsonData := []byte(`{"count": 2}`)
		userId := time.Now().UnixNano()
		skuId := int64(9999999999999)

		resp, err := cartClient.AddProductToCart(t, jsonData, userId, skuId)
		require.NoError(t, err)
		require.Equal(t, http.StatusPreconditionFailed, resp.StatusCode)
	})

	t.Run("Добавление продукта для невалидного пользователя.", func(t *testing.T) {
		t.Parallel()

		jsonData := []byte(`{"count": 2}`)
		userId := int64(-5)
		skuId := int64(1148162)
		resp, err := cartClient.AddProductToCart(t, jsonData, userId, skuId)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Получение корзины пользователя.", func(t *testing.T) {
		t.Parallel()

		jsonData := []byte(`{"count": 2}`)
		userId := time.Now().UnixNano()
		skuId := int64(1148162)

		addResp, addErr := cartClient.AddProductToCart(t, jsonData, userId, skuId)
		require.NoError(t, addErr)
		require.Equal(t, http.StatusOK, addResp.StatusCode)

		cartResp, srvErr := cartClient.GetUserCart(t, userId)
		require.NoError(t, srvErr)
		require.Equal(t, http.StatusOK, cartResp.StatusCode)

		parsedResponse, parseErr := parseCartResponse(t, cartResp)
		require.NoError(t, parseErr)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		t.Cleanup(func() {
			cancel()
		})
		product, err := NewProductClient(t).GetProduct(ctx, &pb.GetProductRequest{
			Sku: skuId,
		})

		require.NoError(t, err)
		require.Equal(t, 1, len(parsedResponse.Items))
		require.Equal(t, product.Sku, parsedResponse.Items[0].Sku)
		require.Equal(t, product.Name, parsedResponse.Items[0].Name)
		require.Equal(t, uint32(product.Price), parsedResponse.Items[0].Price)
		require.Equal(t, uint32(2*product.Price), parsedResponse.TotalPrice)
	})

	t.Run("Получение корзины пользователя, у которого ничего нет.", func(t *testing.T) {
		t.Parallel()

		userId := time.Now().UnixNano()

		resp, err := cartClient.GetUserCart(t, userId)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Получение корзины невалидного пользователя", func(t *testing.T) {
		t.Parallel()

		userId := int64(-5)

		resp, err := cartClient.GetUserCart(t, userId)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

type CartItem struct {
	Sku   int64  `json:"sku"`
	Name  string `json:"name"`
	Count uint32 `json:"count"`
	Price uint32 `json:"price"`
}

type UserCartResponse struct {
	Items      []CartItem `json:"items,omitempty"`
	TotalPrice uint32     `json:"total_price,omitempty"`
}

func parseCartResponse(t *testing.T, response *http.Response) (*UserCartResponse, error) {
	t.Helper()
	var resp UserCartResponse

	err := json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
