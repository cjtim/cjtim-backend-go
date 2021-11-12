package repository

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Collection string

const (
	Binance Collection = "binance"
	Files   Collection = "files"
	Urls    Collection = "urls"
	Users   Collection = "users"
)

var (
	Client = &mongo.Client{}
	DB     = &mongo.Database{}
)

func MongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(config.Config.MongoURI))
	if err != nil {
		return nil, err
	}
	err = client.Connect(context.TODO())
	if err != nil {
		return nil, err
	}
	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}
	zap.L().Info("DB connected!")
	return client, nil
}

func GetCollection(col Collection) *mongo.Collection {
	return DB.Collection(string(col))
}
