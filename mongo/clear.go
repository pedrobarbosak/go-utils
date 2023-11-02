package mongo

import "reflect"

func (repo *repository) clear(object Object) error {
	fields := reflect.TypeOf(object).Elem()
	values := reflect.ValueOf(object).Elem()

	num := values.NumField()
	for i := 0; i < num; i++ {
		field := fields.Field(i)
		value := values.Field(i)

		if _, exists := field.Tag.Lookup(embedTag); !exists {
			continue
		}

		v := value
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		switch v.Kind() {
		case reflect.Slice:
			if err := repo.clearSlice(value); err != nil {
				return err
			}

		case reflect.Struct:
			if err := repo.clearStruct(value); err != nil {
				return err
			}
		}
	}

	return nil
}

func (repo *repository) clearSlice(obj reflect.Value) error {
	length := obj.Len()
	for i := 0; i < length; i++ {
		value := obj.Index(i)
		if err := repo.clearStruct(value); err != nil {
			return err
		}
	}

	return nil
}

func (repo *repository) clearStruct(object reflect.Value) error {
	if !object.Type().Implements(reflect.TypeOf((*Object)(nil)).Elem()) {
		return nil
	}

	obj, ok := object.Interface().(Object)
	if !ok {
		return nil
	}

	if object.Kind() == reflect.Pointer {
		object = object.Elem()
	}

	id := obj.GetID()

	num := object.NumField()
	for i := 0; i < num; i++ {
		value := object.Field(i)
		value.Set(reflect.New(value.Type()).Elem())
	}

	obj.SetID(id)
	return nil
}
