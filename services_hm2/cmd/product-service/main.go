package main

import (
	"log"
	"net"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/repository"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/server"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc"
)

func main() {
	repo, err := repository.NewProductRepository("data.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Сервер запущен на порту :50051...")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, server.NewServer(repo))
	grpcError := grpcServer.Serve(listener)
	if grpcError != nil {
		log.Fatal(grpcError)
	}
}
