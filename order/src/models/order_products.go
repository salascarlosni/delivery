package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderProduct struct {
	Id        primitive.ObjectID `bson:"_id"`
	OrderId   primitive.ObjectID `bson:"order_id"`
	ProductId primitive.ObjectID `bson:"product_id"`
	Quantity  int64              `bson:"quantity"`
}
