package repository

import (
	"context"

	"github.com/cjtim/cjtim-backend-go/configs"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type Repository struct {
	FindOne           func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	Find              func(data interface{}, filter interface{}, opts ...*options.FindOptions) error
	FindOneAndReplace func(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error
	InsertOne         func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error)
	DeleteMany        func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	CountDocuments    func(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	UpdateOne         func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	DeleteOne         func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
}

var (
	Client      = &mongo.Client{}
	DB          = &mongo.Database{}
	BinanceRepo = &Repository{}
	FileRepo    = &Repository{}
	UrlRepo     = &Repository{}
	UserRepo    = &Repository{}
)

func MongoClient() (*mongo.Client, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(configs.Config.MongoURI))
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

	Client = client
	DB = client.Database(configs.Config.MongoDB)
	BinanceRepo = getCollection(Binance)
	FileRepo = getCollection(Files)
	UrlRepo = getCollection(Urls)
	UserRepo = getCollection(Users)
	return client, nil
}

func RestoreRepoMock() {
	BinanceRepo = &Repository{}
	FileRepo = &Repository{}
	UrlRepo = &Repository{}
	UserRepo = &Repository{}
}

func newRepository(col *mongo.Collection) *Repository {
	return &Repository{
		FindOne: func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
			result := col.FindOne(context.TODO(), filter, opts...)
			return result.Decode(data)
		},
		Find: func(data interface{}, filter interface{}, opts ...*options.FindOptions) error {
			cur, err := col.Find(context.TODO(), filter, opts...)
			if err != nil {
				return err
			}
			err = cur.All(context.TODO(), data)
			if err != nil {
				return err
			}
			return err
		},
		FindOneAndReplace: func(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
			result := col.FindOneAndReplace(ctx, filter, replacement, opts...)
			return result.Decode(replacement)
		},
		InsertOne: func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
			result, err := col.InsertOne(ctx, document, opts...)
			id := result.InsertedID.(primitive.ObjectID)
			return id, err
		},
		DeleteMany: func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
			result, err := col.DeleteMany(ctx, filter, opts...)
			return result.DeletedCount, err
		},
		CountDocuments: col.CountDocuments,
		UpdateOne: func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
			result, err := col.UpdateOne(ctx, filter, update, opts...)
			return result.MatchedCount, err
		},
		DeleteOne: func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
			result, err := col.DeleteOne(ctx, filter, opts...)
			return result.DeletedCount, err
		},
	}
}

func getCollection(col Collection) *Repository {
	colInstance := DB.Collection(string(col))
	return newRepository(colInstance)
}
