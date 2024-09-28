package payload

import "fmt"

var _ Loader = &EnvfileLoader[any]{}
var _ Loader = &EnvVarsLoader[string]{}

type Loader interface {
	Load() (*map[string]any, error)
}

type Payload[K comparable] struct {
	Raw  Loader
	data map[string]any
}

func New[K comparable](raw Loader) (*Payload[K], error) {
	c := &Payload[K]{Raw: raw}
	data, err := raw.Load()
	if err != nil {
		return &Payload[K]{}, err
	}

	c.data = *data
	return c, nil
}

func (c Payload[T]) Get(key T) (any, bool) {
	_key := stringify(key)
	v, ok := c.data[_key]
	if !ok {
		return nil, false
	}

	return v, true
}

func (c Payload[T]) GetBool(key T) (bool, bool) {
	v, ok := c.Get(key)
	if !ok {
		return false, false
	}

	vv, ok := v.(bool)
	return vv, ok
}

func (c Payload[T]) GetInt(key T) (int, bool) {
	v, ok := c.Get(key)
	if !ok {
		return 0, false
	}

	vv, ok := v.(int)
	return vv, ok
}

func (c Payload[T]) GetString(key T) (string, bool) {
	v, ok := c.Get(key)
	if !ok {
		return "", false
	}

	vv, ok := v.(string)
	return vv, ok
}

func (c *Payload[T]) Set(key T, value any) {
	_key := stringify(key)
	c.data[_key] = value
}

func stringify(v any) string {
	switch vv := v.(type) {
	case string:
		return vv
	case fmt.Stringer:
		return vv.String()
	case fmt.GoStringer:
		return vv.GoString()
	default:
		return fmt.Sprintf("%s", v)
	}
}
