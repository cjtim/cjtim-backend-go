package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserScheama struct {
	ID                     primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	LineUid                string             `json:"lineUid,omitempty" bson:"lineUid,omitempty"`
	FirebaseMessagingToken string             `json:"fbMsgToken,omitempty" bson:"fbMsgToken,omitempty"`
}
