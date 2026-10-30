//go:build race

package repository

import (
	"math/rand/v2"
	"sync"
	"testing"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/models"
	"github.com/stretchr/testify/require"
)

const producers = 100

func TestCartRepositoryRace_AddProduct(t *testing.T) {
	t.Parallel()
	repo := NewCartRepository()

	var wg sync.WaitGroup
	allCartItems := make([]models.CartItem, 0)
	for i := 0; i < producers; i++ {
		wg.Add(1)
		cartItem := generateCartItem(int64(i))
		allCartItems = append(allCartItems, cartItem)
		go func() {
			defer wg.Done()
			err := repo.AddProduct(1, cartItem)
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	cart, cartErr := repo.GetUserCart(1)
	require.NoError(t, cartErr)
	require.Len(t, cart.CartItems, len(allCartItems))
	require.ElementsMatch(t, allCartItems, cart.CartItems)
}

func TestCartRepositoryRace_AddCount(t *testing.T) {
	t.Parallel()
	repo := NewCartRepository()

	var wg sync.WaitGroup
	for i := 0; i < producers-1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.AddProduct(1, models.CartItem{
				Sku:   1,
				Count: 1 + uint32(i),
			})
			if err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()

	cart, cartErr := repo.GetUserCart(1)
	require.NoError(t, cartErr)
	require.Len(t, cart.CartItems, 1)
	require.Equal(t, uint32(solveCount()), cart.CartItems[0].Count)
}

func generateCartItem(skuID int64) models.CartItem {
	return models.CartItem{
		Sku:   skuID,
		Count: rand.Uint32(),
	}
}

func solveCount() int {
	return producers * (producers - 1) / 2
}
