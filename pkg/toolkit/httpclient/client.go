package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const pkgName = "pkg/http/client"

type Config struct {
	BaseURI string
	Token   string
	Logger  Logger
	Headers map[string]string
}

type Logger interface {
	Infof(string, ...any)
	Debugf(string, ...any)
}

type Client struct {
	client *http.Client
	config *Config
}

type Request struct {
	client   *Client
	Path     string
	method   string
	header   map[string]string
	response *http.Response
}

func New(conf *Config) *Client {
	return &Client{
		client: client(),
		config: conf,
	}
}

func client() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
	}
}

func (c *Client) SetToken(token string) {
	c.config.Token = token
}

func (c *Client) SetHeader(key, value string) {
	headers := c.config.Headers
	if headers == nil {
		headers = map[string]string{}
	}

	c.config.Headers = headers
}

func (c Client) Prepare(path string, method string) *Request {
	return &Request{
		client: &c,
		Path:   path,
		method: method,
		header: c.config.Headers,
	}
}

func (r Request) RequestURI() string {
	return fmt.Sprintf("%s%s", r.client.config.BaseURI, r.Path)
}

func (r Request) Token() string {
	return r.client.config.Token
}

func (r Request) Method() string {
	return r.method
}

func (r Request) Logger() Logger {
	return r.client.config.Logger
}

func (r Request) Header() map[string]string {
	return r.header
}

func (r *Request) SetHeader(m map[string]string) {
	r.header = m
}

func (r Request) Do(req *http.Request) (*http.Response, error) {
	return r.client.client.Do(req)
}

func (r *Request) UnmarshalBody(res *http.Response, v any) error {
	r.response = res

	defer r.response.Body.Close()
	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	r.Logger().Debugf("%s: response body: %s", pkgName, string(buf))

	if err = json.Unmarshal(buf, &v); err != nil {
		return err
	}

	return nil
}

type Doer interface {
	Do(*http.Request) (*http.Response, error)
	RequestURI() string
	Method() string
	Token() string
	Logger() Logger
	Header() map[string]string
	UnmarshalBody(*http.Response, any) error
}

func Do[K, V any](ctx context.Context, d Doer, body *K, response *V) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	d.Logger().Debugf("%s: request uri: %s, method: %s", pkgName, d.RequestURI(), d.Method())
	d.Logger().Debugf("%s: request body: %s", pkgName, string(b))

	var reqBody io.Reader = nil
	if body != nil && (d.Method() == http.MethodPost || d.Method() == http.MethodPut || d.Method() == http.MethodPatch) {
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, d.Method(), d.RequestURI(), reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if d.Token() != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.Token()))
	}
	for k, v := range d.Header() {
		req.Header.Set(k, v)
	}

	res, err := d.Do(req)
	if err != nil {
		return err
	}
	d.Logger().Debugf("%s: response status: %s", pkgName, res.Status)

	if err := d.UnmarshalBody(res, response); err != nil {
		return err
	}

	return nil
}
