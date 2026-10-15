package server

import (
	"context"
	"errors"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/models"
	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/producterrors"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	pb.UnimplementedProductServiceServer
	productRepo ProductRepository
}

type ProductRepository interface {
	GetProduct(sku int64) (models.Product, error)
}

func NewServer(productRepo ProductRepository) *Server {
	return &Server{productRepo: productRepo}
}

func (s *Server) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	validateErr := req.Validate()
	if validateErr != nil {
		return nil, status.Error(codes.InvalidArgument, "Идентификатор товара должен быть натуральным числом (больше нуля)")
	}

	product, serviceErr := s.productRepo.GetProduct(req.Sku)
	if serviceErr != nil {
		if errors.Is(serviceErr, producterrors.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "Информацию можно получить только по существующему товару")
		}
		return nil, status.Error(codes.Internal, "Проблемы из-за неисправностей в системе")
	}

	return &pb.GetProductResponse{
		Sku:   req.Sku,
		Name:  product.Name,
		Price: int32(product.Price),
	}, nil

}
