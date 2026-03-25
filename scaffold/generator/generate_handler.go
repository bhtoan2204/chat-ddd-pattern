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

func GenerateHandler(endpoints []models.Endpoint) (string, error) {
	tmpl, err := template.ParseFiles("scaffold/template/handler.tmpl")
	if err != nil {
		return "", err
	}
	if len(endpoints) == 0 {
		return "", errors.New("no endpoints to generate handler")
	}

	seen := make(map[string]bool)
	created := 0
	skipped := 0
	for _, ep := range endpoints {
		if !shouldGenerateHandler(ep) {
			continue
		}
		module, err := moduleForUsecase(ep.Usecase.Name)
		if err != nil {
			return "", err
		}
		key := module.FsRoot + ":handler:" + ep.Handler
		if seen[key] {
			continue
		}
		seen[key] = true
		written, err := writeHandlerFile(tmpl, module, ep)
		if err != nil {
			return "", err
		}
		if written {
			created++
		} else {
			skipped++
		}
	}

	return fmt.Sprintf("generated %d handler(s), skipped %d existing file(s)", created, skipped), nil
}

func shouldGenerateHandler(ep models.Endpoint) bool {
	if ep.Handler == "" || ep.Usecase.Method == "" || ep.Usecase.Name == "" {
		return false
	}
	return true
}

func writeHandlerFile(tmpl *template.Template, module modulePaths, ep models.Endpoint) (bool, error) {
	busKind := busKindForEndpoint(ep)
	data := handlerTemplateData{
		PackageName:      "handler",
		HandlerName:      ep.Handler,
		StructName:       lowerFirst(strings.TrimSuffix(ep.Handler, "Handler")) + "Handler",
		BusField:         busKind + "Bus",
		BusPackage:       busKind,
		UsecaseMethod:    ep.Usecase.Method,
		Method:           strings.ToUpper(ep.Method),
		RequestStruct:    ep.Request.Struct,
		BusImport:        module.ImportRoot + "/application/" + busKind,
		RequestDtoImport: module.ImportRoot + "/application/dto/in",
	}

	fileName := utils.Snake(ep.Handler) + "_handler.go"
	dst := filepath.Join(module.FsRoot, "transport/http/handler", fileName)
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return false, err
	}
	if fileExists(dst) {
		return false, nil
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return false, err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return false, fmt.Errorf("format handler failed: %w", err)
	}
	if err := os.WriteFile(dst, formatted, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

type handlerTemplateData struct {
	PackageName      string
	HandlerName      string
	StructName       string
	BusField         string
	BusPackage       string
	UsecaseMethod    string
	Method           string
	RequestStruct    string
	BusImport        string
	RequestDtoImport string
}

func busKindForEndpoint(ep models.Endpoint) string {
	if strings.EqualFold(ep.Method, "GET") {
		return "query"
	}
	return "command"
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
