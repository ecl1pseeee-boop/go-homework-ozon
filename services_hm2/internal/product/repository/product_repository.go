package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/models"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/producterrors"
)

type ProductRepository struct {
	products map[int64]models.Product
}

func NewProductRepository(filePath string) (*ProductRepository, error) {
	products, err := parseData(filePath)
	if err != nil {
		return nil, err
	}
	return &ProductRepository{
		products: products,
	}, nil
}

func parseData(filePath string) (map[int64]models.Product, error) {
	data, readErr := os.ReadFile(filePath)
	if readErr != nil {
		return nil, readErr
	}
	storage := make(map[int64]models.Product)
	err := json.Unmarshal(data, &storage)
	if err != nil {
		return nil, fmt.Errorf("ошибка парсинга данных: %w", err)
	}

	return storage, nil
}

func (repo *ProductRepository) GetProduct(sku int64) (models.Product, error) {
	product, ok := repo.products[sku]
	if !ok {
		return models.Product{}, producterrors.ErrProductNotFound
	}
	return product, nil
}
