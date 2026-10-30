package repository

import (
	"slices"
	"sync"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/carterrors"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"
)

type CartRepository struct {
	mu    sync.RWMutex
	carts map[models.UserID]*models.Cart
}

func NewCartRepository() *CartRepository {
	return &CartRepository{
		carts: make(map[models.UserID]*models.Cart),
	}
}

func (repo *CartRepository) AddProduct(userId models.UserID, cartItem models.CartItem) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	userCart, ok := repo.carts[userId]
	if !ok {
		userCart = repo.createUserCart(userId)
	}

	cartItems := userCart.CartItems
	for id, ci := range cartItems {
		if ci.Sku == cartItem.Sku {
			userCart.CartItems[id].Count += cartItem.Count
			return nil
		}
	}
	userCart.CartItems = append(userCart.CartItems, cartItem)

	return nil
}

func (repo *CartRepository) GetUserCart(userId models.UserID) (models.Cart, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	userCart, ok := repo.carts[userId]
	if !ok {
		return models.Cart{}, carterrors.ErrCartNotFound
	}
	answer := models.Cart{
		UserID:    userCart.UserID,
		CartItems: slices.Clone(userCart.CartItems),
	}
	return answer, nil
}

func (repo *CartRepository) createUserCart(userId models.UserID) *models.Cart {
	uc := models.Cart{
		UserID:    userId,
		CartItems: []models.CartItem{},
	}
	repo.carts[userId] = &uc
	return &uc
}
