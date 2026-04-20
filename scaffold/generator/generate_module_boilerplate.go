package generator

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"text/template"

	"wechat-clone/scaffold/models"
)

func GenerateModuleBoilerplate(spec *models.APISpec) (string, error) {
	if spec == nil {
		return "", errors.New("api spec is nil")
	}
	if len(spec.Endpoints) == 0 {
		return "", errors.New("no endpoints to generate module boilerplate")
	}

	groups, err := groupEndpointsByModule(spec.Endpoints)
	if err != nil {
		return "", err
	}

	reposTemplate, err := template.ParseFiles("scaffold/template/module_repos.tmpl")
	if err != nil {
		return "", err
	}
	repositoryTemplate, err := template.ParseFiles("scaffold/template/module_repository.tmpl")
	if err != nil {
		return "", err
	}
	serverBuilderTemplate, err := template.ParseFiles("scaffold/template/module_http_server_builder.tmpl")
	if err != nil {
		return "", err
	}

	created := 0
	skipped := 0
	for _, group := range groups {
		moduleAlreadyExists := fileExists(group.Module.FsRoot)
		files := []struct {
			path string
			kind string
			body []byte
		}{}

		reposBody, err := renderTemplate(reposTemplate, moduleReposTemplateData{})
		if err != nil {
			return "", err
		}
		files = append(files, struct {
			path string
			kind string
			body []byte
		}{
			path: filepath.Join(group.Module.FsRoot, "domain", "repos", "repos.go"),
			kind: "module-repos",
			body: reposBody,
		})

		repositoryBody, err := renderTemplate(repositoryTemplate, moduleRepositoryTemplateData{
			ModuleImportRoot: group.Module.ImportRoot,
		})
		if err != nil {
			return "", err
		}
		files = append(files, struct {
			path string
			kind string
			body []byte
		}{
			path: filepath.Join(group.Module.FsRoot, "infra", "persistent", "repository", "repos_impl.go"),
			kind: "module-repository",
			body: repositoryBody,
		})

		serverBuilderBody, err := renderTemplate(serverBuilderTemplate, buildModuleHTTPServerBuilderData(group))
		if err != nil {
			return "", err
		}
		files = append(files, struct {
			path string
			kind string
			body []byte
		}{
			path: filepath.Join(group.Module.FsRoot, "assembly", "server_builder.go"),
			kind: "module-http-server-builder",
			body: serverBuilderBody,
		})

		for _, file := range files {
			written, err := writeGeneratedModuleFile(file.path, file.kind, file.body, !moduleAlreadyExists)
			if err != nil {
				return "", err
			}
			if written {
				created++
			} else {
				skipped++
			}
		}
	}

	return fmt.Sprintf("generated %d module boilerplate file(s), skipped %d existing hand-written file(s)", created, skipped), nil
}

func renderTemplate(tmpl *template.Template, data any) ([]byte, error) {
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	return formatted, nil
}

func writeGeneratedModuleFile(path, kind string, body []byte, allowCreate bool) (bool, error) {
	if !fileExists(path) && !allowCreate {
		return false, nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	if fileExists(path) && !isGeneratedFile(path, kind) {
		return false, nil
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return false, err
	}
	return true, nil
}

type moduleReposTemplateData struct{}

type moduleRepositoryTemplateData struct {
	ModuleImportRoot string
}

type moduleHTTPServerBuilderTemplateData struct {
	CommandImport string
	CommandAlias  string
	QueryImport   string
	QueryAlias    string
	RepoImport    string
	RepoAlias     string
	ServerImport  string
	ServerAlias   string
	HasCommand    bool
	HasQuery      bool
	ReposVar      string
	Dispatchers   []moduleHTTPServerDispatcher
	ServerArgs    []string
}

type moduleHTTPServerDispatcher struct {
	VarName string
	Ctor    string
}

func buildModuleHTTPServerBuilderData(group moduleEndpoints) moduleHTTPServerBuilderTemplateData {
	data := moduleHTTPServerBuilderTemplateData{
		CommandImport: group.Module.ImportRoot + "/application/command",
		CommandAlias:  modulePackageName(group.Module.ImportRoot) + "command",
		QueryImport:   group.Module.ImportRoot + "/application/query",
		QueryAlias:    modulePackageName(group.Module.ImportRoot) + "query",
		RepoImport:    group.Module.ImportRoot + "/infra/persistent/repository",
		RepoAlias:     modulePackageName(group.Module.ImportRoot) + "repo",
		ServerImport:  group.Module.ImportRoot + "/transport/server",
		ServerAlias:   modulePackageName(group.Module.ImportRoot) + "server",
		ReposVar:      modulePackageName(group.Module.ImportRoot) + "Repos",
	}

	seen := make(map[string]bool)
	for _, ep := range group.Endpoints {
		dispatcherName := dispatcherParamName(ep)
		if seen[dispatcherName] {
			continue
		}
		seen[dispatcherName] = true

		ctorPackage := data.CommandAlias
		if applicationHandlerKind(ep) == "query" {
			ctorPackage = data.QueryAlias
			data.HasQuery = true
		} else {
			data.HasCommand = true
		}

		data.Dispatchers = append(data.Dispatchers, moduleHTTPServerDispatcher{
			VarName: dispatcherName,
			Ctor:    fmt.Sprintf("%s.New%s(appContext, %s)", ctorPackage, ep.Usecase.Method, data.ReposVar),
		})
		data.ServerArgs = append(data.ServerArgs, dispatcherName)
	}

	return data
}
