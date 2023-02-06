package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restaurant struct {
	Id        primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Image     string             `bson:"image"`
	OpenHour  time.Time          `bson:"open_hour"`
	CloseHour time.Time          `bson:"close_hour"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	DeletedAt time.Time          `bson:"deleted_at,omitempty"`
}
