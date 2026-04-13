package generator

import (
	"bytes"
	"go/format"
	"path/filepath"
	"strings"
	"testing"
	"text/template"

	"go-socket/scaffold/models"
)

func TestRequestTemplateFormats(t *testing.T) {
	tmpl, err := template.ParseFiles(filepath.Join("..", "template", "request.tmpl"))
	if err != nil {
		t.Fatalf("parse template failed: %v", err)
	}

	searchFields := mapRequestFields([]models.FieldSpec{
		{Name: "room_id", Type: "string", Required: true},
		{Name: "q", Type: "string"},
		{Name: "limit", Type: "int"},
	})
	searchData := requestTemplateData{
		PackageName:       "in",
		StructName:        "SearchChatMentionsRequest",
		Fields:            searchFields,
		NeedsErrors:       requestNeedsErrors(searchFields),
		NeedsStrings:      requestNeedsStrings(searchFields, nil),
		HasNormalize:      requestHasNormalize(searchFields),
		AdditionalStructs: nil,
	}

	var searchBuf bytes.Buffer
	if err := tmpl.Execute(&searchBuf, searchData); err != nil {
		t.Fatalf("execute search template failed: %v", err)
	}

	searchFormatted, err := format.Source(searchBuf.Bytes())
	if err != nil {
		t.Fatalf("format search output failed: %v\n%s", err, searchBuf.String())
	}

	searchOutput := string(searchFormatted)
	if !strings.Contains(searchOutput, "Q      string `json:\"q\" form:\"q\"`") {
		t.Fatalf("expected short field name Q for q, got:\n%s", searchOutput)
	}
	if !strings.Contains(searchOutput, "r.Q = strings.TrimSpace(r.Q)") {
		t.Fatalf("expected Normalize to trim Q, got:\n%s", searchOutput)
	}

	messageFields := []models.FieldSpec{
		{Name: "room_id", Type: "string", Required: true},
		{
			Name: "mentions",
			Type: "array",
			Items: &models.Payload{
				Struct: "SendChatMessageMentionRequest",
				Fields: []models.FieldSpec{
					{Name: "account_id", Type: "string"},
				},
			},
		},
	}
	parentFields := mapRequestFields(messageFields)
	additional := mapRequestNestedStructs(t.TempDir(), filepath.Join(t.TempDir(), "send_chat_message_request.go"), messageFields)
	messageData := requestTemplateData{
		PackageName:       "in",
		StructName:        "SendChatMessageRequest",
		Fields:            parentFields,
		AdditionalStructs: additional,
		NeedsErrors:       requestNeedsErrors(parentFields),
		NeedsStrings:      requestNeedsStrings(parentFields, additional),
		HasNormalize:      requestHasNormalize(parentFields),
	}

	var messageBuf bytes.Buffer
	if err := tmpl.Execute(&messageBuf, messageData); err != nil {
		t.Fatalf("execute message template failed: %v", err)
	}

	messageFormatted, err := format.Source(messageBuf.Bytes())
	if err != nil {
		t.Fatalf("format message output failed: %v\n%s", err, messageBuf.String())
	}

	messageOutput := string(messageFormatted)
	if !strings.Contains(messageOutput, "func (r *SendChatMessageMentionRequest) Normalize()") {
		t.Fatalf("expected nested Normalize method, got:\n%s", messageOutput)
	}
	if !strings.Contains(messageOutput, "r.AccountID = strings.TrimSpace(r.AccountID)") {
		t.Fatalf("expected nested Normalize to trim account id, got:\n%s", messageOutput)
	}
	if !strings.Contains(messageOutput, "for idx := range r.Mentions {") || !strings.Contains(messageOutput, "r.Mentions[idx].Normalize()") {
		t.Fatalf("expected parent Normalize to cascade nested Normalize, got:\n%s", messageOutput)
	}
}
