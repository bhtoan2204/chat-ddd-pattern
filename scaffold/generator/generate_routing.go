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
)

const ROUTING_DESTINATION_PATH = "core/delivery/http/routing.go"

func GenerateRouting(spec *models.APISpec) (string, error) {
	if spec == nil {
		return "", errors.New("api spec is nil")
	}
	if len(spec.Endpoints) == 0 {
		return "", errors.New("no endpoints to generate routing")
	}

	tmpl, err := template.ParseFiles("scaffold/template/routing.tmpl")
	if err != nil {
		return "", err
	}
	data := routingTemplateData{
		PackageName:   "http",
		PublicRoutes:  buildRoutingLines(spec.Endpoints, false),
		PrivateRoutes: buildRoutingLines(spec.Endpoints, true),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format routing failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(ROUTING_DESTINATION_PATH), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(ROUTING_DESTINATION_PATH, formatted, 0o644); err != nil {
		return "", err
	}
	return "generated routing.go", nil
}

type routingTemplateData struct {
	PackageName   string
	PublicRoutes  []string
	PrivateRoutes []string
}

func buildRoutingLines(endpoints []models.Endpoint, auth bool) []string {
	lines := make([]string, 0, len(endpoints))
	seen := make(map[string]bool)
	for _, ep := range endpoints {
		if ep.Auth != auth {
			continue
		}
		method := strings.ToUpper(ep.Method)
		path := ep.Path
		if path == "" {
			continue
		}
		key := method + ":" + path
		if seen[key] {
			continue
		}
		seen[key] = true
		lines = append(lines, fmt.Sprintf("\troutes.%s(%q, h.Handle())", method, path))
	}
	return lines
}
