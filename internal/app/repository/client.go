package repository

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/configs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type IClient interface {
	Connect() error
	Disconnect() error
	GetDatabase() *mongo.Database
	GetClient() *mongo.Client
}

type ClientImpl struct {
	mongoClient *mongo.Client
	db          *mongo.Database
}

var (
	Client IClient = &ClientImpl{}
)

func (r *ClientImpl) Connect() error {
	c, err := mongo.NewClient(options.Client().ApplyURI(configs.Config.MongoURI))
	if err != nil {
		return err
	}

	if err := c.Connect(context.TODO()); err != nil {
		return err
	}

	if err := c.Ping(context.TODO(), nil); err != nil {
		return err
	}
	zap.L().Info("DB connected!")

	r.mongoClient = c
	r.db = c.Database(configs.Config.MongoDB)

	BinanceRepo = &RepoImpl{col: r.db.Collection(string(Binance))}
	FileRepo = &RepoImpl{col: r.db.Collection(string(Files))}
	UrlRepo = &RepoImpl{col: r.db.Collection(string(Urls))}
	UserRepo = &RepoImpl{col: r.db.Collection(string(Users))}
	return nil
}

func (r *ClientImpl) Disconnect() error {
	zap.L().Info("DB disconnecting...")
	err := r.mongoClient.Disconnect(context.TODO())
	if err != nil {
		zap.L().Error("DB disconect error!", zap.Error(err))
		return err
	}
	zap.L().Info("DB disconected!")
	return nil
}

func (r *ClientImpl) GetDatabase() *mongo.Database {
	return r.db
}

func (r *ClientImpl) GetClient() *mongo.Client {
	return r.mongoClient
}
