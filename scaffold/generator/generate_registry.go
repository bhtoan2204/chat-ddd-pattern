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

const REGISTRY_DESTINATION_PATH = "core/delivery/http/registry.go"

func GenerateRegistry(spec *models.APISpec) (string, error) {
	if spec == nil {
		return "", errors.New("api spec is nil")
	}
	if len(spec.Endpoints) == 0 {
		return "", errors.New("no endpoints to generate registry")
	}

	if fileExists(REGISTRY_DESTINATION_PATH) {
		return updateRegistry(spec)
	}

	tmpl, err := template.ParseFiles("scaffold/template/registry.tmpl")
	if err != nil {
		return "", err
	}
	data := registryTemplateData{
		PackageName: "http",
		Entries:     buildRegistryEntries(spec),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format registry failed: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(REGISTRY_DESTINATION_PATH), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(REGISTRY_DESTINATION_PATH, formatted, 0o644); err != nil {
		return "", err
	}
	return "generated registry.go", nil
}

type registryTemplateData struct {
	PackageName string
	Entries     []string
}

func buildRegistryEntries(spec *models.APISpec) []string {
	base := strings.TrimRight(spec.BasePath, "/")
	entries := make([]string, 0, len(spec.Endpoints))
	seen := make(map[string]bool)
	for _, ep := range spec.Endpoints {
		if ep.Method == "" || ep.Path == "" || ep.Handler == "" {
			continue
		}
		key := strings.ToUpper(ep.Method) + ":" + base + ep.Path
		if seen[key] {
			continue
		}
		seen[key] = true
		entries = append(entries, fmt.Sprintf("\t\t%q: {\n\t\t\thandler: handler.New%s(usecase),\n\t\t},", key, ep.Handler))
	}
	return entries
}

func updateRegistry(spec *models.APISpec) (string, error) {
	data, err := os.ReadFile(REGISTRY_DESTINATION_PATH)
	if err != nil {
		return "", err
	}
	content := string(data)
	inserts := make([]string, 0)
	for _, entry := range buildRegistryEntries(spec) {
		routeKey := extractRouteKey(entry)
		if routeKey == "" {
			continue
		}
		if strings.Contains(content, fmt.Sprintf("%q:", routeKey)) {
			continue
		}
		inserts = append(inserts, entry)
	}
	updated, ok := insertBeforeMapClose(content, inserts)
	if !ok {
		return "", fmt.Errorf("failed to update registry.go")
	}
	formatted, err := format.Source([]byte(updated))
	if err != nil {
		return "", fmt.Errorf("format registry failed: %w", err)
	}
	if err := os.WriteFile(REGISTRY_DESTINATION_PATH, formatted, 0o644); err != nil {
		return "", err
	}
	return "updated registry.go", nil
}

func insertBeforeMapClose(content string, inserts []string) (string, bool) {
	start := strings.Index(content, "return map[string]routingConfig{")
	if start == -1 {
		return "", false
	}
	sub := content[start:]
	end := strings.LastIndex(sub, "\t}")
	if end == -1 {
		return "", false
	}
	insertAt := start + end
	block := strings.Join(inserts, "\n") + "\n"
	return content[:insertAt] + block + content[insertAt:], true
}

func extractRouteKey(entry string) string {
	first := strings.Index(entry, "\"")
	if first == -1 {
		return ""
	}
	rest := entry[first+1:]
	second := strings.Index(rest, "\"")
	if second == -1 {
		return ""
	}
	return rest[:second]
}
