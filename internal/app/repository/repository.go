package repository

import (
	"context"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection string

const (
	Binance Collection = "binance"
	Files   Collection = "files"
	Urls    Collection = "urls"
	Users   Collection = "users"
)

type Repository interface {
	FindOne(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	Find(data interface{}, filter interface{}, opts ...*options.FindOptions) error
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	Health() error
}

type RepoImpl struct {
	col *mongo.Collection
}

var (
	BinanceRepo Repository = &RepoImpl{}
	FileRepo    Repository = &RepoImpl{}
	UrlRepo     Repository = &RepoImpl{}
	UserRepo    Repository = &RepoImpl{}
)

func RestoreRepoMock() {
	BinanceRepo = &RepoImpl{}
	FileRepo = &RepoImpl{}
	UrlRepo = &RepoImpl{}
	UserRepo = &RepoImpl{}
}

func Health() error {
	var wg sync.WaitGroup
	errs := []error{}
	all := []func() error{
		BinanceRepo.Health,
		FileRepo.Health,
		UrlRepo.Health,
		UserRepo.Health,
	}
	for _, check := range all {
		wg.Add(1)
		go func(c func() error) {
			defer wg.Done()
			errs = append(errs, c())
		}(check)
	}
	go func() {
		wg.Wait()
	}()
	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RepoImpl) Health() error {
	_, err := r.col.CountDocuments(context.Background(), bson.M{})
	return err
}

func (r *RepoImpl) FindOne(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	result := r.col.FindOne(context.TODO(), filter, opts...)
	return result.Decode(data)
}

func (r *RepoImpl) Find(data interface{}, filter interface{}, opts ...*options.FindOptions) error {
	cur, err := r.col.Find(context.TODO(), filter, opts...)
	if err != nil {
		return err
	}
	err = cur.All(context.TODO(), data)
	if err != nil {
		return err
	}
	return err
}

func (r *RepoImpl) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	result := r.col.FindOneAndReplace(ctx, filter, replacement, opts...)
	return result.Decode(replacement)
}

func (r *RepoImpl) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	result, err := r.col.InsertOne(ctx, document, opts...)
	id := result.InsertedID.(primitive.ObjectID)
	return id, err
}

func (r *RepoImpl) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	result, err := r.col.DeleteMany(ctx, filter, opts...)
	return result.DeletedCount, err
}

func (r *RepoImpl) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.col.CountDocuments(ctx, filter, opts...)
}

func (r *RepoImpl) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	result, err := r.col.UpdateOne(ctx, filter, update, opts...)
	return result.MatchedCount, err
}

func (r *RepoImpl) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	result, err := r.col.DeleteOne(ctx, filter, opts...)
	return result.DeletedCount, err
}
