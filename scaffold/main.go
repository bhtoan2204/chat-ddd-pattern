package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"wechat-clone/scaffold/generator"
	"wechat-clone/scaffold/models"
	scaffoldswagger "wechat-clone/scaffold/swagger"
	"wechat-clone/scaffold/utils"
)

const API_SPEC_DIR = "scaffold/api"
const ASSEMBLY_SPEC_PATH = "scaffold/assembly/modules.yaml"
const PROTO_SPEC_DIR = "scaffold/proto"

const (
	commandAll      = "all"
	commandAPI      = "api"
	commandModule   = "module"
	commandAssembly = "assembly"
	commandProto    = "proto"
	commandSwagger  = "swagger"
)

type cliConfig struct {
	command string
	module  string
	help    bool
}

func main() {
	cfg, err := parseCLI(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n\n", err)
		printHelp(os.Stderr)
		os.Exit(2)
	}
	if cfg.help {
		printHelp(os.Stdout)
		return
	}

	if err := run(cfg); err != nil {
		log.Fatal(err)
	}
}

func parseCLI(args []string) (cliConfig, error) {
	cfg := cliConfig{command: commandAll}
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		cfg.command = strings.ToLower(strings.TrimSpace(args[0]))
		args = args[1:]
	}

	flags := flag.NewFlagSet("scaffold", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)
	flags.BoolVar(&cfg.help, "help", false, "show scaffold generator help")
	flags.BoolVar(&cfg.help, "h", false, "show scaffold generator help")
	flags.StringVar(&cfg.module, "module", "", "module name for module-scoped generation")
	if err := flags.Parse(args); err != nil {
		return cliConfig{}, err
	}
	if flags.NArg() > 0 {
		return cliConfig{}, fmt.Errorf("unexpected argument %q", flags.Arg(0))
	}

	switch cfg.command {
	case commandAll, commandAPI, commandModule, commandAssembly, commandProto, commandSwagger, "help":
		if cfg.command == "help" {
			cfg.help = true
		}
	default:
		return cliConfig{}, fmt.Errorf("unsupported generate command %q", cfg.command)
	}

	if cfg.command == commandModule && strings.TrimSpace(cfg.module) == "" {
		return cliConfig{}, errors.New("module generation requires --module <name>")
	}
	if cfg.command != commandModule && cfg.module != "" {
		return cliConfig{}, fmt.Errorf("--module is only supported with %q", commandModule)
	}

	cfg.module = normalizeModuleName(cfg.module)
	return cfg, nil
}

func printHelp(output *os.File) {
	fmt.Fprint(output, `Scaffold generator

Usage:
  go run scaffold/main.go [command] [flags]

Commands:
  all                    Generate API scaffold, assembly builders, proto outputs, and swagger JSON (default)
  api                    Generate API scaffold outputs from all scaffold/api specs
  module --module <name> Generate API scaffold and assembly builders for one module
  assembly               Generate assembly builders from scaffold/assembly/modules.yaml
  proto                  Generate protobuf Go files from scaffold/proto
  swagger                Generate OpenAPI JSON from scaffold/api specs
  help                   Show this help

Examples:
  go run scaffold/main.go all
  go run scaffold/main.go swagger
  go run scaffold/main.go proto
  go run scaffold/main.go module --module payment
  go run scaffold/main.go module --module room

Make examples:
  make generate
  make generate WHAT=swagger
  make generate WHAT=proto
  make generate WHAT=module MODULE=payment
`)
}

func run(cfg cliConfig) error {
	switch cfg.command {
	case commandAll:
		if err := runAPI(""); err != nil {
			return fmt.Errorf("generate api scaffold failed: %w", err)
		}
		if err := runAssembly(""); err != nil {
			return fmt.Errorf("generate assembly failed: %w", err)
		}
		if err := runProto(); err != nil {
			return fmt.Errorf("generate proto failed: %w", err)
		}
		if err := runSwagger(); err != nil {
			return fmt.Errorf("generate swagger failed: %w", err)
		}
	case commandAPI:
		if err := runAPI(""); err != nil {
			return fmt.Errorf("generate api scaffold failed: %w", err)
		}
	case commandModule:
		if err := runAPI(cfg.module); err != nil {
			return fmt.Errorf("generate module api scaffold failed: %w", err)
		}
		if err := runAssembly(cfg.module); err != nil {
			return fmt.Errorf("generate module assembly failed: %w", err)
		}
	case commandAssembly:
		if err := runAssembly(""); err != nil {
			return fmt.Errorf("generate assembly failed: %w", err)
		}
	case commandProto:
		if err := runProto(); err != nil {
			return fmt.Errorf("generate proto failed: %w", err)
		}
	case commandSwagger:
		if err := runSwagger(); err != nil {
			return fmt.Errorf("generate swagger failed: %w", err)
		}
	}

	return nil
}

func runAPI(moduleName string) error {
	apiSpec, err := loadAPISpec(moduleName)
	if err != nil {
		return fmt.Errorf("load API spec failed: %w", err)
	}

	steps := []struct {
		name string
		run  func() (string, error)
	}{
		{name: "requests", run: func() (string, error) { return generator.GenerateRequest(apiSpec.Endpoints) }},
		{name: "responses", run: func() (string, error) { return generator.GenerateResponse(apiSpec.Endpoints) }},
		{name: "module boilerplate", run: func() (string, error) { return generator.GenerateModuleBoilerplate(apiSpec) }},
		{name: "application handlers", run: func() (string, error) { return generator.GenerateApplicationHandler(apiSpec.Endpoints) }},
		{name: "handlers", run: func() (string, error) { return generator.GenerateHandler(apiSpec.Endpoints) }},
		{name: "routing", run: func() (string, error) { return generator.GenerateRouting(apiSpec) }},
		{name: "registry", run: func() (string, error) { return generator.GenerateRegistry(apiSpec) }},
	}
	for _, step := range steps {
		msg, err := step.run()
		if err != nil {
			return fmt.Errorf("generate %s failed: %w", step.name, err)
		}
		fmt.Println(msg)
	}

	return nil
}

func runAssembly(moduleName string) error {
	assemblySpec, err := models.LoadAssemblySpec(ASSEMBLY_SPEC_PATH)
	if err != nil {
		return fmt.Errorf("load assembly spec failed: %w", err)
	}
	if moduleName != "" {
		assemblySpec = filterAssemblySpec(assemblySpec, moduleName)
		if len(assemblySpec.Modules) == 0 {
			return fmt.Errorf("module %q is not registered in %s", moduleName, ASSEMBLY_SPEC_PATH)
		}
	}

	msg, err := generator.GenerateAssembly(assemblySpec)
	if err != nil {
		return fmt.Errorf("generate assembly builders failed: %w", err)
	}
	fmt.Println(msg)
	return nil
}

func runProto() error {
	msg, err := generator.GenerateProto(PROTO_SPEC_DIR)
	if err != nil {
		return fmt.Errorf("generate proto files failed: %w", err)
	}
	fmt.Println(msg)
	return nil
}

func runSwagger() error {
	swaggerSpec, err := scaffoldswagger.GenerateDefault()
	if err != nil {
		return fmt.Errorf("generate swagger json failed: %w", err)
	}
	fmt.Printf("generated swagger json at %s\n", swaggerSpec.OutputPath)
	return nil
}

func loadAPISpec(moduleName string) (*models.APISpec, error) {
	apiSpec, err := models.LoadAPISpecDir(API_SPEC_DIR)
	if err != nil {
		return nil, err
	}
	if moduleName == "" {
		return apiSpec, nil
	}

	filtered := *apiSpec
	filtered.Endpoints = make([]models.Endpoint, 0)
	for _, endpoint := range apiSpec.Endpoints {
		if endpointModuleName(endpoint) == moduleName || endpointSpecFileName(endpoint) == moduleName {
			filtered.Endpoints = append(filtered.Endpoints, endpoint)
		}
	}
	if len(filtered.Endpoints) == 0 {
		return nil, fmt.Errorf("module %q has no endpoints in %s", moduleName, API_SPEC_DIR)
	}
	return &filtered, nil
}

func filterAssemblySpec(spec *models.AssemblySpec, moduleName string) *models.AssemblySpec {
	filtered := &models.AssemblySpec{Modules: make([]models.AssemblyModule, 0)}
	for _, module := range spec.Modules {
		if normalizeModuleName(module.Name) == moduleName {
			filtered.Modules = append(filtered.Modules, module)
		}
	}
	return filtered
}

func endpointModuleName(endpoint models.Endpoint) string {
	moduleName := strings.TrimSuffix(strings.TrimSpace(endpoint.Usecase.Name), "Usecase")
	switch moduleName {
	case "Auth":
		return "account"
	case "Message":
		return "room"
	default:
		return normalizeModuleName(utils.Snake(moduleName))
	}
}

func endpointSpecFileName(endpoint models.Endpoint) string {
	path := strings.Trim(strings.TrimSpace(endpoint.Path), "/")
	if path == "" {
		return ""
	}
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return ""
	}
	candidate := strings.Trim(parts[0], "{}:")
	return normalizeModuleName(filepath.Base(candidate))
}

func normalizeModuleName(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.ReplaceAll(value, "-", "_")
	return value
}
