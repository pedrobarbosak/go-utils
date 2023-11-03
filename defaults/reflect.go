package defaults

import (
	"log"
	"reflect"
	"strconv"
	"time"
)

func (s service) setDefaults(v interface{}) error {
	rv := reflect.ValueOf(v)

	//// Is a pointer aka can we set the value of it
	//if rv.Kind() != reflect.Ptr || rv.IsNil() {
	//	return ErrInvalidValue
	//}

	rv = rv.Elem()

	// Only supports struct for now. But everything inside a struct it works which is stupid.
	if rv.Kind() != reflect.Struct {
		return ErrInvalidValue
	}

	t := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		valueField := rv.Field(i)

		switch valueField.Kind() {
		case reflect.Struct:
			if !valueField.Addr().CanInterface() {
				continue
			}

			iface := valueField.Addr().Interface()
			err := s.setDefaults(iface)
			if err != nil {
				return err
			}

		case reflect.Slice:
			length := valueField.Len()
			for j := 0; j < length; j++ {
				sliceField := valueField.Index(j)
				if sliceField.Kind() != reflect.Pointer {
					sliceField = sliceField.Addr()
				}

				if err := s.SetDefaults(sliceField.Interface()); err != nil {
					return err
				}
			}
		}

		// Field already has a value set
		if !isFieldEmpty(valueField) {
			continue
		}

		typeField := t.Field(i)

		fTag := typeField.Tag.Get(s.tag)

		// Don't initialize field if 'ignoreTag' or it's missing the tag
		// Since we do 'fTag == ""' we end up not initializing slices and pointers - && valueField.Kind() != reflect.Slice ...
		if fTag == ignoreTag || (s.ignoreOnMissing && fTag == "") {
			continue
		}

		// It's a field that we can set
		if !valueField.CanSet() {
			return ErrUnexportedField
		}

		// Set the field with the corresponding value
		if err := set(typeField.Type, valueField, fTag); err != nil {
			return err
		}
	}

	return nil
}

func set(t reflect.Type, f reflect.Value, value string) error {
	switch t.Kind() {
	case reflect.Ptr:
		ptr := reflect.New(t.Elem())
		err := set(t.Elem(), ptr.Elem(), value)
		if err != nil {
			return err
		}
		f.Set(ptr)
	case reflect.String:
		f.SetString(value)
	case reflect.Bool:
		v, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		f.SetBool(v)
	case reflect.Float32:
		v, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return err
		}
		f.SetFloat(v)
	case reflect.Float64:
		v, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		f.SetFloat(v)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.PkgPath() == "time" && t.Name() == "Duration" {
			duration, err := time.ParseDuration(value)
			if err != nil {
				return err
			}

			f.Set(reflect.ValueOf(duration))
			break
		}

		v, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		f.SetInt(int64(v))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		f.SetUint(v)
	default:
		// @todo: initialize slices, maps & structs pointers
		log.Println("Kind:", t.Kind())
		return ErrUnsupportedType
	}

	return nil
}

func isFieldEmpty(field reflect.Value) bool {
	return reflect.DeepEqual(field.Interface(), reflect.Zero(field.Type()).Interface())
}
