package swagger

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"go-socket/scaffold/models"
)

func TestBuildDocument(t *testing.T) {
	doc, err := BuildDocument(&models.APISpec{
		Version:  1,
		BasePath: "/api/v1",
		Endpoints: []models.Endpoint{
			{
				Name:    "PaymentProcessWebhook",
				Method:  "POST",
				Path:    "/payment/webhooks/:provider",
				Handler: "ProcessWebhookHandler",
				Usecase: models.Usecase{Name: "PaymentUsecase", Method: "ProcessWebhook"},
				Request: models.Payload{
					Struct: "ProcessWebhookRequest",
					Fields: []models.FieldSpec{
						{Name: "provider", Type: "string", Source: "path", Required: true},
						{Name: "signature", Type: "string", Source: "header", Header: "X-Signature"},
						{Name: "payload", Type: "string", Source: "raw_body", Required: true},
					},
				},
				Response: models.Payload{
					Struct: "ProcessWebhookResponse",
					Fields: []models.FieldSpec{
						{Name: "status", Type: "string"},
					},
				},
			},
			{
				Name:          "PaymentCreatePayment",
				Method:        "POST",
				Path:          "/payment/intents",
				Auth:          true,
				SuccessStatus: 201,
				Handler:       "CreatePaymentHandler",
				Usecase:       models.Usecase{Name: "PaymentUsecase", Method: "CreatePayment"},
				Request: models.Payload{
					Struct: "CreatePaymentRequest",
					Fields: []models.FieldSpec{
						{Name: "provider", Type: "string", Required: true},
						{Name: "amount", Type: "int64", Required: true},
						{Name: "metadata", Type: "object"},
					},
				},
				Response: models.Payload{
					Struct: "CreatePaymentResponse",
					Fields: []models.FieldSpec{
						{Name: "transaction_id", Type: "string"},
					},
				},
			},
			{
				Name:    "RoomList",
				Method:  "GET",
				Path:    "/room/list",
				Handler: "ListRoomsHandler",
				Usecase: models.Usecase{Name: "RoomUsecase", Method: "ListRooms"},
				Request: models.Payload{
					Struct: "ListRoomsRequest",
					Fields: []models.FieldSpec{
						{Name: "page", Type: "int"},
						{Name: "limit", Type: "int"},
					},
				},
				Response: models.Payload{
					Struct: "ListRoomsResponse",
					Fields: []models.FieldSpec{
						{Name: "page", Type: "int"},
					},
				},
			},
		},
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got := doc.Servers[0].URL; got != "/api/v1" {
		t.Fatalf("unexpected server url: %s", got)
	}

	webhookPath := doc.Paths["/payment/webhooks/{provider}"]
	if webhookPath == nil || webhookPath.Post == nil {
		t.Fatalf("expected normalized webhook path and POST operation")
	}
	if webhookPath.Post.RequestBody == nil {
		t.Fatalf("expected raw body request body")
	}
	if webhookPath.Post.Parameters[0].In != "header" || webhookPath.Post.Parameters[1].In != "path" {
		t.Fatalf("expected header and path parameters, got %+v", webhookPath.Post.Parameters)
	}

	createPath := doc.Paths["/payment/intents"]
	if createPath == nil || createPath.Post == nil {
		t.Fatalf("expected create payment operation")
	}
	if len(createPath.Post.Security) != 1 {
		t.Fatalf("expected bearer auth security, got %+v", createPath.Post.Security)
	}
	if _, exists := doc.Components.SecuritySchemes["BearerAuth"]; !exists {
		t.Fatalf("expected bearer auth security scheme")
	}
	if got := createPath.Post.Responses["201"]; got == nil {
		t.Fatalf("expected 201 response")
	}

	roomListPath := doc.Paths["/room/list"]
	if roomListPath == nil || roomListPath.Get == nil {
		t.Fatalf("expected room list GET operation")
	}
	if len(roomListPath.Get.Parameters) != 2 {
		t.Fatalf("expected query parameters for GET, got %+v", roomListPath.Get.Parameters)
	}
	if roomListPath.Get.RequestBody != nil {
		t.Fatalf("did not expect request body for GET")
	}

	t.Run("request body schema excludes header parameters", func(t *testing.T) {
		loginDoc, err := BuildDocument(&models.APISpec{
			Version:  1,
			BasePath: "/api/v1",
			Endpoints: []models.Endpoint{
				{
					Name:    "AuthLogin",
					Method:  "POST",
					Path:    "/auth/login",
					Handler: "LoginHandler",
					Usecase: models.Usecase{Name: "AuthUsecase", Method: "Login"},
					Request: models.Payload{
						Struct: "LoginRequest",
						Fields: []models.FieldSpec{
							{Name: "email", Type: "string", Required: true},
							{Name: "password", Type: "string", Required: true},
							{Name: "device_uid", Type: "string", Source: "header", Header: "X-Device-UID", Required: true},
							{Name: "user_agent", Type: "string", Source: "header", Header: "User-Agent"},
						},
					},
					Response: models.Payload{
						Struct: "LoginResponse",
						Fields: []models.FieldSpec{
							{Name: "access_token", Type: "string"},
						},
					},
				},
			},
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		loginPath := loginDoc.Paths["/auth/login"]
		if loginPath == nil || loginPath.Post == nil {
			t.Fatalf("expected login POST operation")
		}
		if len(loginPath.Post.Parameters) != 2 {
			t.Fatalf("expected 2 header parameters, got %+v", loginPath.Post.Parameters)
		}

		loginSchema := loginDoc.Components.Schemas["Auth_LoginRequest"]
		if loginSchema == nil {
			t.Fatalf("expected auth.LoginRequest schema")
		}
		if _, exists := loginSchema.Properties["device_uid"]; exists {
			t.Fatalf("did not expect header field device_uid in request body schema")
		}
		if _, exists := loginSchema.Properties["user_agent"]; exists {
			t.Fatalf("did not expect header field user_agent in request body schema")
		}
		if _, exists := loginSchema.Properties["email"]; !exists {
			t.Fatalf("expected email in request body schema")
		}
		if _, exists := loginSchema.Properties["password"]; !exists {
			t.Fatalf("expected password in request body schema")
		}
	})
}

func TestGenerateWritesOpenAPIJSON(t *testing.T) {
	tmpDir := t.TempDir()
	specDir := filepath.Join(tmpDir, "api")
	outDir := filepath.Join(tmpDir, "swagger")

	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatalf("create spec dir failed: %v", err)
	}

	spec := `version: 1
basePath: /api/v1
endpoints:
  - name: HealthPing
    method: GET
    path: /health/ping
    handler: HealthHandler
    usecase:
      name: HealthUsecase
      method: Ping
    request:
      struct: PingRequest
      fields:
        - name: verbose
          type: bool
    response:
      struct: PingResponse
      fields:
        - name: ok
          type: bool
`
	if err := os.WriteFile(filepath.Join(specDir, "health.yaml"), []byte(spec), 0o644); err != nil {
		t.Fatalf("write spec file failed: %v", err)
	}

	result, err := Generate(specDir, outDir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.OutputPath == "" {
		t.Fatalf("expected output path")
	}
	if _, err := os.Stat(result.OutputPath); err != nil {
		t.Fatalf("expected generated file to exist, got %v", err)
	}

	var doc Document
	if err := json.Unmarshal(result.JSON, &doc); err != nil {
		t.Fatalf("expected valid json, got %v", err)
	}
	if doc.Paths["/health/ping"] == nil {
		t.Fatalf("expected /health/ping path in generated document")
	}
}
