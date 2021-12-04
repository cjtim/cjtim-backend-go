package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type URLScheama struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RebrandlyID string             `json:"id,omitempty" bson:"id,omitempty"`
	ShortURL    string             `json:"shortUrl,omitempty" bson:"shortUrl,omitempty"`
	Destination string             `json:"destination,omitempty" bson:"destination,omitempty"`
	LineUID     string             `json:"lineUid,omitempty" bson:"lineUid,omitempty"`
}
