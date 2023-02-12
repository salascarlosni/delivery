package servers

import (
	"context"
	"products/src/proto"

	"go.mongodb.org/mongo-driver/mongo"
)

type IProductServer interface {
}

type ProductServer struct {
	proto.UnimplementedProductServiceServer
	coll *mongo.Collection
}

func NewRProductServer(coll *mongo.Collection) *ProductServer {
	return &ProductServer{coll: coll}
}

func (srv *ProductServer) CreateProduct(
	ctx context.Context,
	args *proto.CreateProductRequest,
) (*proto.CreateProductResponse, error) {
	return nil, nil
}
