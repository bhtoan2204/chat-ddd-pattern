package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path"
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
		if err := writeHandlerFile(tmpl, module, ep); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("generated %d handler(s)", len(seen)), nil
}

func shouldGenerateHandler(ep models.Endpoint) bool {
	if ep.Handler == "" || ep.Usecase.Method == "" || ep.Usecase.Name == "" {
		return false
	}
	return true
}

func writeHandlerFile(tmpl *template.Template, module modulePaths, ep models.Endpoint) error {
	data := handlerTemplateData{
		PackageName:      "handler",
		HandlerName:      ep.Handler,
		StructName:       lowerFirst(strings.TrimSuffix(ep.Handler, "Handler")) + "Handler",
		UsecaseName:      ep.Usecase.Name,
		UsecaseField:     lowerFirst(ep.Usecase.Name),
		UsecaseMethod:    ep.Usecase.Method,
		Method:           strings.ToUpper(ep.Method),
		RequestStruct:    ep.Request.Struct,
		UsecaseImport:    path.Join(module.ImportRoot, "application/usecase"),
		RequestDtoImport: path.Join(module.ImportRoot, "application/dto/in"),
	}

	fileName := utils.Snake(ep.Handler) + "_handler.go"
	dst := filepath.Join(module.FsRoot, "transport/http/handler", fileName)
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
	PackageName      string
	HandlerName      string
	StructName       string
	UsecaseName      string
	UsecaseField     string
	UsecaseMethod    string
	Method           string
	RequestStruct    string
	UsecaseImport    string
	RequestDtoImport string
}

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}
