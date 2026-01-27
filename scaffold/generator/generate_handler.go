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

const HANDLER_DESTINATION_PATH = "core/delivery/http/handler"

func GenerateHandler(endpoints []models.Endpoint) (string, error) {
	tmpl, err := template.ParseFiles("scaffold/template/handler.tmpl")
	if err != nil {
		return "", err
	}
	if len(endpoints) == 0 {
		return "", errors.New("no endpoints to generate handler")
	}

	seen := make(map[string]bool)
	for _, ep := range endpoints {
		if !shouldGenerateHandler(ep, seen) {
			continue
		}
		if err := writeHandlerFile(tmpl, ep); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("generated %d handler(s)", len(seen)), nil
}

func shouldGenerateHandler(ep models.Endpoint, seen map[string]bool) bool {
	if ep.Handler == "" || ep.Usecase.Method == "" || ep.Usecase.Name == "" {
		return false
	}
	if seen[ep.Handler] {
		return false
	}
	seen[ep.Handler] = true
	return true
}

func writeHandlerFile(tmpl *template.Template, ep models.Endpoint) error {
	data := handlerTemplateData{
		PackageName:    "handler",
		HandlerName:    ep.Handler,
		StructName:     lowerFirst(strings.TrimSuffix(ep.Handler, "Handler")) + "Handler",
		UsecaseName:    ep.Usecase.Name,
		UsecaseField:   lowerFirst(ep.Usecase.Name),
		UsecaseGetter:  ep.Usecase.Name,
		UsecaseMethod:  ep.Usecase.Method,
		Method:         strings.ToUpper(ep.Method),
		RequestStruct:  ep.Request.Struct,
		ResponseStruct: ep.Response.Struct,
		ResponseField:  responseFieldName(ep.Response.Fields),
		RequestArgs:    requestArgNames(ep.Request.Fields),
	}

	fileName := utils.Snake(ep.Handler) + "_handler.go"
	dst := filepath.Join(HANDLER_DESTINATION_PATH, fileName)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format handler failed: %w", err)
	}
	return os.WriteFile(dst, formatted, 0o644)
}

type handlerTemplateData struct {
	PackageName    string
	HandlerName    string
	StructName     string
	UsecaseName    string
	UsecaseField   string
	UsecaseGetter  string
	UsecaseMethod  string
	Method         string
	RequestStruct  string
	ResponseStruct string
	ResponseField  string
	RequestArgs    []string
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func responseFieldName(fields []models.FieldSpec) string {
	if len(fields) == 1 {
		return utils.Pascal(fields[0].Name)
	}
	return ""
}

func requestArgNames(fields []models.FieldSpec) []string {
	args := make([]string, 0, len(fields))
	for _, f := range fields {
		args = append(args, utils.Pascal(f.Name))
	}
	return args
}
