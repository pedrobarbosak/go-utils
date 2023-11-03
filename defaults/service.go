package defaults

import "reflect"

type Defaulter interface {
	SetDefaults(s interface{}) error
}

type service struct {
	ignoreOnMissing bool
	tag             string
}

func (s service) SetDefaults(v interface{}) error {
	rv := reflect.ValueOf(v)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrInvalidValue
	}

	return s.setDefaults(v)
}

func New(cfg ...Config) Defaulter {
	if len(cfg) != 0 {
		return service{ignoreOnMissing: cfg[0].IgnoreOnMissing, tag: cfg[0].Tag}
	}

	c := DefaultConfig()
	return service{ignoreOnMissing: c.IgnoreOnMissing, tag: c.Tag}
}

func SetDefaults(v interface{}) error {
	return New().SetDefaults(v)
}
