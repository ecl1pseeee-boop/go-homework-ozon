//go:build api

package test

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"time"

	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	GRPCAddress = "localhost:50051"
	CartAddress = "http://localhost:8080"
)

func NewProductClient(t *testing.T) pb.ProductServiceClient {
	t.Helper()
	conn, err := grpc.NewClient(GRPCAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Ошибка при создании GRPC соединения: %v", err)
	}
	t.Cleanup(func() {
		conn.Close()
	})
	return pb.NewProductServiceClient(conn)
}

type CartClient struct {
	Client  *http.Client
	Address string
}

func NewCartClient(t *testing.T) CartClient {
	t.Helper()
	cartClient := CartClient{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Address: CartAddress,
	}

	t.Cleanup(func() {
		cartClient.Client.CloseIdleConnections()
	})
	return cartClient
}

func (c *CartClient) GetUserCart(t *testing.T, userID int64) (*http.Response, error) {
	t.Helper()

	getAddress := c.Address + fmt.Sprintf("/user/%d/cart", userID)
	resp, err := c.Client.Get(getAddress)
	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		closeErr := resp.Body.Close()
		if closeErr != nil {
			require.NoError(t, closeErr)
		}
	})

	return resp, nil
}

func (c *CartClient) AddProductToCart(t *testing.T, jsonData []byte, userID int64, skuID int64) (*http.Response, error) {
	t.Helper()

	postAddress := c.Address + fmt.Sprintf("/user/%d/cart/%d", userID, skuID)
	resp, err := c.Client.Post(postAddress, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		return nil, err
	}

	t.Cleanup(func() {
		resp.Body.Close()
	})

	return resp, nil
}
