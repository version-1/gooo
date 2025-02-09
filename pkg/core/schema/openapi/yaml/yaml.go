package yaml

import (
	yaml "gopkg.in/yaml.v3"
)

func Unmarshal(b []byte, d any) error {
	return yaml.Unmarshal(b, d)
}

type OrderedMap[T any] struct {
	keys   []string
	Values map[string]T
}

func (o *OrderedMap[T]) Set(key string, value T) {
	o.keys = append(o.keys, key)
	o.Values[key] = value
}

func (o OrderedMap[T]) Get(key string) T {
	return o.Values[key]
}

func (o OrderedMap[T]) Each(cb func(key string, v T) error) error {
	for _, key := range o.keys {
		err := cb(key, o.Values[key])
		if err != nil {
			return err
		}
	}

	return nil
}

func (o OrderedMap[T]) Index(i int) (string, T) {
	key := o.keys[i]
	return key, o.Values[key]
}

func (o OrderedMap[T]) Len() int {
	return len(o.keys)
}

func (o OrderedMap[T]) Keys() []string {
	return o.keys
}

func (o *OrderedMap[T]) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode {
		return nil
	}

	o.Values = make(map[string]T)
	for i := 0; i < len(node.Content); i += 2 {
		key := node.Content[i].Value
		value := node.Content[i+1]
		var v T
		if err := value.Decode(&v); err != nil {
			return err
		}
		o.Set(key, v)
	}

	return nil
}
