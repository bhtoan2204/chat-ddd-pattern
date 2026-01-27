package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"go-socket/scaffold/models"
	"go-socket/scaffold/utils"
)

const REQUEST_DESTINATION_PATH = "core/delivery/http/data/in"

func GenerateRequest(endpoints []models.Endpoint) (string, error) {
	tmpl, err := template.ParseFiles("scaffold/template/request.tmpl")
	if err != nil {
		return "", err
	}
	if len(endpoints) == 0 {
		return "", errors.New("no endpoints to generate request")
	}

	seen := make(map[string]bool)
	for _, ep := range endpoints {
		if ep.Request.Struct == "" {
			continue
		}
		if seen[ep.Request.Struct] {
			continue
		}
		seen[ep.Request.Struct] = true

		data := requestTemplateData{
			PackageName: "in",
			StructName:  ep.Request.Struct,
			Fields:      mapRequestFields(ep.Request.Fields),
		}

		fileName := utils.Snake(ep.Request.Struct) + "_request.go"
		dst := filepath.Join(REQUEST_DESTINATION_PATH, fileName)
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return "", err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return "", err
		}
		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return "", fmt.Errorf("format request DTO failed: %w", err)
		}
		if err := os.WriteFile(dst, formatted, 0o644); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("generated %d request DTO(s)", len(seen)), nil
}

type requestTemplateData struct {
	PackageName string
	StructName  string
	Fields      []requestField
}

type requestField struct {
	GoName     string
	Type       string
	JSONName   string
	BindingTag string
	Required   bool
	ZeroCheck  string
}

func mapRequestFields(fields []models.FieldSpec) []requestField {
	result := make([]requestField, 0, len(fields))
	for _, f := range fields {
		goName := utils.Pascal(f.Name)
		binding := ""
		if f.Required {
			if strings.Contains(strings.ToLower(f.Name), "email") {
				binding = "required,email"
			} else {
				binding = "required"
			}
		}
		result = append(result, requestField{
			GoName:     goName,
			Type:       utils.GoType(f.Type),
			JSONName:   f.Name,
			BindingTag: binding,
			Required:   f.Required,
			ZeroCheck:  utils.ZeroCheck(utils.GoType(f.Type), goName),
		})
	}
	return result
}
