package services

import (
	"context"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/dto"
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

func (c *fakeProductClient) GetProduct(ctx context.Context, skuId int64) (*dto.GetProductResponseDTO, error) {
	if c.err != nil {
		return nil, c.err
	}

	p, ok := c.products[skuId]
	if !ok {
		return nil, carterrors.ErrProductNotFound
	}

	return &dto.GetProductResponseDTO{
		Sku:   skuId,
		Name:  p.Name,
		Price: int32(p.Price),
	}, nil
}
