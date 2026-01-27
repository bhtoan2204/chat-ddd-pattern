package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go-socket/scaffold/models"
	"go-socket/scaffold/utils"
	"go/format"
	"os"
	"path/filepath"
	"text/template"
)

const RESPONSE_DESTINATION_PATH = "core/delivery/http/data/out"

func GenerateResponse(endpoints []models.Endpoint) (string, error) {
	tmpl, err := template.ParseFiles("scaffold/template/response.tmpl")
	if err != nil {
		return "", err
	}
	if len(endpoints) == 0 {
		return "", errors.New("no endpoints to generate response")
	}

	seen := make(map[string]bool)
	for _, ep := range endpoints {
		if ep.Response.Struct == "" {
			continue
		}
		if seen[ep.Response.Struct] {
			continue
		}
		seen[ep.Response.Struct] = true

		data := responseTemplateData{
			PackageName: "out",
			StructName:  ep.Response.Struct,
			Fields:      mapResponseFields(ep.Response.Fields),
		}

		fileName := utils.Snake(ep.Response.Struct) + "_response.go"
		dst := filepath.Join(RESPONSE_DESTINATION_PATH, fileName)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return "", fmt.Errorf("format response DTO failed: %w", err)
		}
		if err := os.WriteFile(dst, formatted, 0o644); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("generated %d response DTO(s)", len(seen)), nil
}

type responseTemplateData struct {
	PackageName string
	StructName  string
	Fields      []responseField
}

type responseField struct {
	GoName   string
	Type     string
	JSONName string
}

func mapResponseFields(fields []models.FieldSpec) []responseField {
	result := make([]responseField, 0, len(fields))
	for _, f := range fields {
		result = append(result, responseField{
			GoName:   utils.Pascal(f.Name),
			Type:     utils.GoType(f.Type),
			JSONName: f.Name,
		})
	}
	return result
}
