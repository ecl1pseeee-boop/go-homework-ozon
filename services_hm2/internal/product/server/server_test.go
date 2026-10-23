package server

import (
	"context"
	"testing"

	"github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/internal/product/repository"
	pb "github.com/ecl1pseeee-boop/go-homework-ozon/services_hm2/pkg/api/product"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const filePath = "testdata/fake_data.json"

func TestServer_Fail(t *testing.T) {
	t.Parallel()
	repo := newFakeProductRepository()
	repo.err = status.Error(codes.DataLoss, "data loss")
	server := NewServer(repo)

	response, serverErr := server.GetProduct(context.Background(), &pb.GetProductRequest{
		Sku: 1,
	})
	st, ok := status.FromError(serverErr)
	require.True(t, ok)

	require.Equal(t, codes.Internal, st.Code())
	require.Nil(t, response)

}

func TestServer_GetProduct_Success(t *testing.T) {
	t.Parallel()

	repo, err := repository.NewProductRepository(filePath)
	require.NoError(t, err)
	server := NewServer(repo)

	response, serverErr := server.GetProduct(context.Background(), &pb.GetProductRequest{
		Sku: 1,
	})

	require.NoError(t, serverErr)
	require.Equal(t, int64(1), response.Sku)
	require.Equal(t, int32(1), response.Price)
	require.Equal(t, "test1", response.Name)
}

func TestServer_GetProduct_Fail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		testName       string
		ctx            context.Context
		productRequest *pb.GetProductRequest
		wantStatusCode codes.Code
	}{
		{
			testName: "Не валидный идентификатор товара",
			ctx:      context.Background(),
			productRequest: &pb.GetProductRequest{
				Sku: 0,
			},
			wantStatusCode: codes.InvalidArgument,
		},
		{
			testName: "Не существующий товар.",
			ctx:      context.Background(),
			productRequest: &pb.GetProductRequest{
				Sku: 5,
			},
			wantStatusCode: codes.NotFound,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testName, func(t *testing.T) {
			t.Parallel()
			repo, err := repository.NewProductRepository(filePath)
			require.NoError(t, err)
			server := NewServer(repo)

			response, serverErr := server.GetProduct(tc.ctx, tc.productRequest)
			st, ok := status.FromError(serverErr)
			require.True(t, ok)

			require.Equal(t, tc.wantStatusCode, st.Code())
			require.Nil(t, response)
		})
	}

}
