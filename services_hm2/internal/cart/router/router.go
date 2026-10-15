package router

import (
	"net/http"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/controller"
)

func SetupRouter(handler *controller.CartHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", handler.AddProductHandler)
	mux.HandleFunc("GET /user/{user_id}/cart", handler.GetCartHandler)

	return mux
}
