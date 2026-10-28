package main

import (
	"log"
	"net/http"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/clients/product"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/controller"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/repository"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/router"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/cart/services"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	cartRepository := repository.NewCartRepository()
	productClient := product.NewProductClient(pb.NewProductServiceClient(conn))
	cartService := services.NewCartService(cartRepository, productClient)
	cartHandler := controller.NewCartHandler(cartService)

	mux := router.SetupRouter(cartHandler)
	log.Println("Сервер запущен на порту :8080...")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
