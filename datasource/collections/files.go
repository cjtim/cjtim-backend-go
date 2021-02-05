package collections

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type FileScheama struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	URL     string             `json:"url,omitempty" bson:"url,omitempty"`
	LineUID string             `json:"lineUid,omitempty" bson:"lineUid,omitempty"`
}
