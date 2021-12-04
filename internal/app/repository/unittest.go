package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UnitTestSchema struct {
	ID   primitive.ObjectID `json:"_id,omitempty"`
	Test string             `json:"test,omitempty"`
}
