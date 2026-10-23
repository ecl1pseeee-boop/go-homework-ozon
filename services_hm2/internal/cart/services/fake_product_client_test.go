package services

import (
	"context"

	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type fakeProduct struct {
	Name  string
	Price int64
}

type fakeProductClient struct {
	products map[int64]fakeProduct
	err      error
}

func newFakeProductClient(products map[int64]fakeProduct) *fakeProductClient {
	return &fakeProductClient{
		products: products,
	}
}

func (c *fakeProductClient) GetProduct(ctx context.Context, in *pb.GetProductRequest, opts ...grpc.CallOption) (*pb.GetProductResponse, error) {
	if c.err != nil {
		return nil, c.err
	}

	skuId := in.Sku
	p, ok := c.products[skuId]
	if !ok {
		return nil, status.Error(codes.NotFound, "fakeProduct not found")
	}

	return &pb.GetProductResponse{
		Sku:   skuId,
		Name:  p.Name,
		Price: int32(p.Price),
	}, nil
}
