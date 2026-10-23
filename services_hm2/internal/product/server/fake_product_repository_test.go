package server

import (
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/models"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/producterrors"
)

type fakeProductRepository struct {
	products map[int64]models.Product
	err      error
}

func newFakeProductRepository() *fakeProductRepository {
	return &fakeProductRepository{
		products: map[int64]models.Product{},
	}
}

func (repo *fakeProductRepository) GetProduct(sku int64) (models.Product, error) {
	if repo.err != nil {
		return models.Product{}, repo.err
	}
	product, ok := repo.products[sku]
	if !ok {
		return models.Product{}, producterrors.ErrProductNotFound
	}
	return product, nil
}
