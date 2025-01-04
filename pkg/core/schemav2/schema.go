package schemav2

import (
	"os"

	"gopkg.in/yaml.v3"
)

func New(path string) (*RootSchema, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := &RootSchema{}
	if err := yaml.Unmarshal(bytes, &s); err != nil {
		return s, err
	}

	return s, nil
}

type RequestBody struct {
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
}
type Responses map[string]Response

type Response struct {
	Description string                 `json:"description"`
	Content     map[string]interface{} `json:"content"`
}

type Content struct {
	Schema RootSchema `json:"schema"`
}

type Parameter struct {
	Name        string     `json:"name"`
	In          string     `json:"in"`
	Description string     `json:"description"`
	Required    bool       `json:"required"`
	Schema      RootSchema `json:"schema"`
}

type Operation struct {
	Summary     string      `json:"summary"`
	Description string      `json:"description"`
	OperationId string      `json:"operationId"`
	Parameters  []Parameter `json:"parameters"`
	RequestBody RequestBody `json:"requestBody"`
	Responses   Responses   `json:"responses"`
}

type PathItem struct {
	Get    Operation `json:"get"`
	Post   Operation `json:"post"`
	Put    Operation `json:"put"`
	Delete Operation `json:"delete"`
}

type Info struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Server struct {
	Url         string `json:"url"`
	Description string `json:"description"`
}

type Components struct {
	Schemas map[string]Schema `json:"schemas"`
}

type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Ref        string              `json:"$ref"`
}

type Property struct {
	Type   string `json:"type"`
	Format string `json:"format"`
	Sample string `json:"sample"`
}

type RootSchema struct {
	OpenAPI    string              `json:"openapi"`
	Info       Info                `json:"info"`
	Paths      map[string]PathItem `json:"paths"`
	Servers    []Server            `json:"servers"`
	Components Components          `json:"components"`
}