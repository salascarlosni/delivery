package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	Id           primitive.ObjectID `bson:"_id"`
	UserId       primitive.ObjectID `bson:"user_id"`
	RestaurantId string             `bson:"restaurant_id"`
	Status       string             `bson:"status"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	CanceledAt   time.Time          `bson:"canceled_at,omitempty"`
}
