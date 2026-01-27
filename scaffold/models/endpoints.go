package models

import (
	"os"

	"github.com/goccy/go-yaml"
)

type APISpec struct {
	Version   int        `json:"version" yaml:"version"`
	BasePath  string     `json:"basePath" yaml:"basePath"`
	Endpoints []Endpoint `json:"endpoints" yaml:"endpoints"`
}

type Endpoint struct {
	Name     string  `json:"name" yaml:"name"`
	Method   string  `json:"method" yaml:"method"`
	Path     string  `json:"path" yaml:"path"`
	Handler  string  `json:"handler" yaml:"handler"`
	Auth     bool    `json:"auth,omitempty" yaml:"auth,omitempty"`
	Usecase  Usecase `json:"usecase" yaml:"usecase"`
	Request  Payload `json:"request" yaml:"request"`
	Response Payload `json:"response" yaml:"response"`
}

type Usecase struct {
	Name   string `json:"name" yaml:"name"`
	Method string `json:"method" yaml:"method"`
}

type Payload struct {
	Struct string      `json:"struct" yaml:"struct"`
	Fields []FieldSpec `json:"fields" yaml:"fields"`
}

type FieldSpec struct {
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Required bool   `json:"required,omitempty" yaml:"required,omitempty"`
}

func LoadAPISpec(path string) (*APISpec, error) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var apiSpec APISpec
	err = yaml.Unmarshal(yamlFile, &apiSpec)
	if err != nil {
		return nil, err
	}
	return &apiSpec, nil
}
