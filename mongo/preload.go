package mongo

import (
	"context"
	"reflect"

	"go.mongodb.org/mongo-driver/mongo"
)

func (repo *repository) Preload(ctx context.Context, obj any) error {
	return repo.search(ctx, obj)
}

func (repo *repository) search(ctx context.Context, obj interface{}) error {
	values := reflect.ValueOf(obj)
	fields := reflect.TypeOf(obj)

	if values.Kind() == reflect.Ptr {
		values = values.Elem()
		fields = fields.Elem()
	}

	if values.Kind() != reflect.Struct {
		if values.Kind() == reflect.Slice {
			return repo.searchSlice(ctx, values)
		}

		return nil
	}

	num := values.NumField()
	for i := 0; i < num; i++ {
		value := values.Field(i)
		field := fields.Field(i)

		collection, exists := field.Tag.Lookup(embedTag)
		if exists {
			if err := repo.fetch(ctx, value, collection); err != nil {
				return err
			}
			continue
		}

		v := value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Slice:
			if err := repo.searchSlice(ctx, value); err != nil {
				return err
			}

		case reflect.Struct:
			if err := repo.search(ctx, value.Interface()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *repository) searchSlice(ctx context.Context, obj reflect.Value) error {
	length := obj.Len()
	for i := 0; i < length; i++ {
		if err := repo.search(ctx, obj.Index(i).Interface()); err != nil {
			return err
		}
	}
	return nil
}

func (repo *repository) fetch(ctx context.Context, value reflect.Value, collection string) error {
	if value.Kind() == reflect.Slice {
		length := value.Len()
		for i := 0; i < length; i++ {
			if err := repo.fetchObject(ctx, value.Index(i), collection); err != nil {
				return err
			}
		}

		return nil
	}

	return repo.fetchObject(ctx, value, collection)
}

func (repo *repository) fetchObject(ctx context.Context, value reflect.Value, collection string) error {
	if value.IsNil() {
		return nil
	}

	if !value.Type().Implements(reflect.TypeOf((*Object)(nil)).Elem()) {
		return nil
	}

	storableObject, ok := value.Interface().(StorableObject)
	if ok {
		id := storableObject.GetID()
		if id == "" {
			return nil
		}

		return repo.GetByID(ctx, id, storableObject)
	}

	if collection == "" {
		return nil
	}

	object, ok := value.Interface().(Object)
	if !ok {
		return nil
	}

	id := object.GetID()
	if id == "" {
		return nil
	}

	return repo.getObjectByID(ctx, id, object, collection)
}

func (repo *repository) getObjectByID(ctx context.Context, objectID string, object Object, collection string) error {
	filter, err := repo.getIDFilter(objectID)
	if err != nil {
		return err
	}

	result := repo.database.Collection(collection).FindOne(ctx, filter)
	if err = result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return ErrNoResults
		}
		return err
	}

	if err = result.Decode(object); err != nil {
		return err
	}

	if repo.config.AutoPreload {
		return repo.Preload(ctx, object)
	}

	return nil
}
