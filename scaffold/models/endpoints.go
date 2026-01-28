package models

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
	Name     string   `json:"name" yaml:"name"`
	Type     string   `json:"type" yaml:"type"`
	Items    *Payload `json:"items,omitempty" yaml:"items,omitempty"`
	Required bool     `json:"required,omitempty" yaml:"required,omitempty"`
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

func LoadAPISpecDir(dir string) (*APISpec, error) {
	pattern := filepath.Join(dir, "*.yaml")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no api spec files found in %s", dir)
	}
	sort.Strings(files)
	merged := &APISpec{}
	for _, file := range files {
		spec, err := LoadAPISpec(file)
		if err != nil {
			return nil, err
		}
		if merged.Version == 0 {
			merged.Version = spec.Version
		}
		if merged.BasePath == "" {
			merged.BasePath = spec.BasePath
		}
		if spec.BasePath != "" && merged.BasePath != spec.BasePath {
			return nil, fmt.Errorf("basePath mismatch in %s", file)
		}
		merged.Endpoints = append(merged.Endpoints, spec.Endpoints...)
	}
	return merged, nil
}
