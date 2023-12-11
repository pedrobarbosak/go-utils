package mongo

import (
	"context"
)

const embedTag = "m-embed"

type Repository interface {
	Create(ctx context.Context, object StorableObject) error
	Update(ctx context.Context, objectID string, object StorableObject) error
	GetBy(ctx context.Context, object StorableObject, filters ...Filter) error
	GetByID(ctx context.Context, objectID string, object StorableObject) error
	Fetch(ctx context.Context, object StorableObject, out interface{}, filters ...Filter) error

	WithTransaction(ctx context.Context, fn func(sc context.Context) error) error
	Aggregate(ctx context.Context, object StorableObject, query string, out interface{}) error
	Count(ctx context.Context, object StorableObject, filter interface{}) (int64, error)

	UpdateOne(ctx context.Context, object StorableObject, filter interface{}, update interface{}) (int64, error)

	CreateMany(ctx context.Context, obj StorableObject, data []interface{}) error
	DeleteAll(ctx context.Context, object StorableObject) error

	CreateUniqueIndexes(ctx context.Context, obj StorableObject, values []map[string]int) error

	Preload(ctx context.Context, object any) error
	Disconnect(ctx context.Context) error
}

type Object interface {
	GetID() string
	SetID(id string)
}

type StorableObject interface {
	Object
	GetCollection() string
}

type Filter struct {
	Key   string
	Value interface{}
}
