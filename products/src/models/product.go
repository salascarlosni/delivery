package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	Id           primitive.ObjectID `bson:"_id"`
	Name         string             `bson:"name"`
	Image        string             `bson:"image"`
	RestaurantId primitive.ObjectID `bson:"restaurant_id"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	DeletedAt    time.Time          `bson:"deleted_at,omitempty"`
}
