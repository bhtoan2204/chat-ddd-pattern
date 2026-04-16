package generator

import (
	"bytes"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"go-socket/scaffold/models"
)

func TestRegistryTemplateIncludesSocketHookWhenWebsocketTransportExists(t *testing.T) {
	moduleRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(moduleRoot, "transport", "websocket"), 0o755); err != nil {
		t.Fatalf("mkdir websocket transport failed: %v", err)
	}

	group := moduleEndpoints{
		Module: modulePaths{
			FsRoot:     moduleRoot,
			ImportRoot: "go-socket/core/modules/room",
		},
		Endpoints: []models.Endpoint{
			{
				Auth: true,
				Usecase: models.Usecase{
					Name:   "RoomUsecase",
					Method: "CreateRoom",
				},
				Request: models.Payload{Struct: "CreateRoomRequest"},
				Response: models.Payload{
					Struct: "CreateRoomResponse",
				},
			},
		},
	}

	tmpl, err := template.ParseFiles(filepath.Join("..", "template", "registry.tmpl"))
	if err != nil {
		t.Fatalf("parse template failed: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, buildRegistryTemplateData(group)); err != nil {
		t.Fatalf("execute template failed: %v", err)
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		t.Fatalf("format output failed: %v\n%s", err, buf.String())
	}

	output := string(formatted)
	if !strings.Contains(output, `"go-socket/core/modules/room/transport/websocket"`) {
		t.Fatalf("expected websocket transport import, got:\n%s", output)
	}
	if !strings.Contains(output, "func (s *roomHTTPServer) RegisterSocketRoutes(") {
		t.Fatalf("expected RegisterSocketRoutes method, got:\n%s", output)
	}
	if !strings.Contains(output, "socketHandler gin.HandlerFunc") {
		t.Fatalf("expected socket handler dependency, got:\n%s", output)
	}
}
