package config

import "os"

type Config struct {
	values map[string]any
	keys   map[string]struct{}
}

func (c Config) Keys() []string {
	res := []string{}
	for key := range c.keys {
		res = append(res, key)
	}

	return res
}

func New(keys []string) *Config {
	c := &Config{}

	for _, key := range keys {
		c.values[key] = os.Getenv(key)
		c.keys[key] = struct{}{}
	}

	return c
}

func (c *Config) Set(key string, value any) {
	c.values[key] = value
}

func (c *Config) Get(key string) (any, bool) {
	v, ok := c.values[key]
	return v, ok
}

func (c *Config) GetInt(key string) (int, bool) {
	v, ok := c.values[key]
	if !ok {
		return 0, false
	}

	n, ok := v.(int)
	return n, ok
}

func (c *Config) GetString(key string) (string, bool) {
	v, ok := c.values[key]
	if !ok {
		return "", false
	}

	s, ok := v.(string)
	return s, ok
}
