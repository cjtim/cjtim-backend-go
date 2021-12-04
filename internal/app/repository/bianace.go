package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BinanceScheama struct {
	ID               primitive.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	LineNotifyToken  string                 `json:"lineNotifyToken,omitempty" bson:"lineNotifyToken,omitempty"`
	BinanceApiKey    string                 `json:"binanceApiKey,omitempty" bson:"binanceApiKey,omitempty"`
	BinanceSecretKey string                 `json:"binanceSecretKey,omitempty" bson:"binanceSecretKey,omitempty"`
	LineUID          string                 `json:"lineUid,omitempty" bson:"lineUid,omitempty"`
	Prices           map[string]interface{} `json:"prices,omitempty" bson:"prices,omitempty"`
	LineNotifyTime   int64                  `json:"lineNotifyTime,omitempty" bson:"lineNotifyTime,omitempty"`
}
