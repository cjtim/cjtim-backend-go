package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mock_Repository struct {
	M_FindOne           func(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error
	M_Find              func(data interface{}, filter interface{}, opts ...*options.FindOptions) error
	M_FindOneAndReplace func(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error
	M_InsertOne         func(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error)
	M_DeleteMany        func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
	M_CountDocuments    func(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
	M_UpdateOne         func(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error)
	M_DeleteOne         func(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error)
}

func (r *Mock_Repository) FindOne(data interface{}, filter interface{}, opts ...*options.FindOneOptions) error {
	return r.M_FindOne(data, filter, opts...)
}

func (r *Mock_Repository) Find(data interface{}, filter interface{}, opts ...*options.FindOptions) error {
	return r.M_Find(data, filter, opts...)
}

func (r *Mock_Repository) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) error {
	return r.M_FindOneAndReplace(ctx, filter, replacement, opts...)
}

func (r *Mock_Repository) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (primitive.ObjectID, error) {
	return r.M_InsertOne(ctx, document, opts...)
}

func (r *Mock_Repository) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	return r.M_DeleteMany(ctx, filter, opts...)
}

func (r *Mock_Repository) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return r.M_CountDocuments(ctx, filter, opts...)
}

func (r *Mock_Repository) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	return r.M_UpdateOne(ctx, filter, update, opts...)
}

func (r *Mock_Repository) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	return r.M_DeleteOne(ctx, filter, opts...)
}
