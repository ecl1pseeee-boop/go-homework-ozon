package product

import (
	"context"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Client struct {
	productClient GRPCClient
}

type GRPCClient interface {
	GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error)
}

func NewProductClient(productClient GRPCClient) *Client {
	return &Client{
		productClient: productClient,
	}
}

func (c *Client) GetProduct(ctx context.Context, skuId int64) (*dto.GetProductResponseDTO, error) {
	productResponse, err := c.productClient.GetProduct(ctx, &pb.GetProductRequest{
		Sku: skuId,
	})

	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return nil, carterrors.ErrProductUnavailable
		}

		switch st.Code() {
		case codes.NotFound:
			return nil, carterrors.ErrProductNotFound
		case codes.InvalidArgument:
			return nil, carterrors.ErrInvalidData
		default:
			return nil, carterrors.ErrProductUnavailable
		}
	}

	return &dto.GetProductResponseDTO{
		Sku:   productResponse.GetSku(),
		Price: productResponse.GetPrice(),
		Name:  productResponse.GetName(),
	}, nil
}
