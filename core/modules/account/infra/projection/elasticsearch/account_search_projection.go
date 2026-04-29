package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	accountprojection "wechat-clone/core/modules/account/application/projection"
	"wechat-clone/core/modules/account/domain/entity"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

const defaultAccountSearchIndex = "accounts_v1"

type accountSearchProjection struct {
	client *es8.Client
	index  string
}

func NewAccountSearchProjection(cfg config.ElasticsearchConfig, client *es8.Client) (accountprojection.SearchProjection, error) {
	if !cfg.Enabled || client == nil {
		return nil, nil
	}

	projection := &accountSearchProjection{
		client: client,
		index:  resolveAccountSearchIndex(cfg),
	}
	if err := projection.ensureIndex(context.Background()); err != nil {
		return nil, stackErr.Error(err)
	}
	return projection, nil
}

func (p *accountSearchProjection) SyncAccount(ctx context.Context, account *entity.Account) error {
	if p == nil || p.client == nil || account == nil {
		return nil
	}

	document := accountDocument{
		ID:                strings.TrimSpace(account.ID),
		Email:             account.Email.Value(),
		DisplayName:       account.DisplayName,
		Username:          cloneStringPtr(account.Username),
		AvatarObjectKey:   cloneStringPtr(account.AvatarObjectKey),
		Status:            account.Status.String(),
		EmailVerifiedAt:   cloneTimePtr(account.EmailVerifiedAt),
		LastLoginAt:       cloneTimePtr(account.LastLoginAt),
		PasswordChangedAt: cloneTimePtr(account.PasswordChangedAt),
		CreatedAt:         account.CreatedAt.UTC(),
		UpdatedAt:         account.UpdatedAt.UTC(),
		BannedReason:      account.BannedReason,
		BannedUntil:       cloneTimePtr(account.BannedUntil),
	}

	body, err := json.Marshal(document)
	if err != nil {
		return stackErr.Error(fmt.Errorf("marshal account search projection failed: %w", err))
	}

	req := esapi.IndexRequest{
		Index:      p.index,
		DocumentID: document.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "false",
	}
	res, err := req.Do(ctx, p.client)
	if err != nil {
		return stackErr.Error(fmt.Errorf("index account search projection failed: %w", err))
	}
	defer res.Body.Close()

	if res.IsError() {
		return stackErr.Error(fmt.Errorf("index account search projection returned status %s: %s", res.Status(), readBody(res.Body)))
	}
	return nil
}

func (p *accountSearchProjection) DeleteAccount(ctx context.Context, accountID string) error {
	if p == nil || p.client == nil || strings.TrimSpace(accountID) == "" {
		return nil
	}

	req := esapi.DeleteRequest{
		Index:      p.index,
		DocumentID: strings.TrimSpace(accountID),
		Refresh:    "false",
	}
	res, err := req.Do(ctx, p.client)
	if err != nil {
		return stackErr.Error(fmt.Errorf("delete account search projection failed: %w", err))
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return nil
	}
	if res.IsError() {
		return stackErr.Error(fmt.Errorf("delete account search projection returned status %s: %s", res.Status(), readBody(res.Body)))
	}
	return nil
}

func (p *accountSearchProjection) ensureIndex(ctx context.Context) error {
	if p == nil || p.client == nil {
		return nil
	}

	existsReq := esapi.IndicesExistsRequest{Index: []string{p.index}}
	existsRes, err := existsReq.Do(ctx, p.client)
	if err != nil {
		return stackErr.Error(fmt.Errorf("check account search index failed: %w", err))
	}
	defer existsRes.Body.Close()

	switch existsRes.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusNotFound:
	default:
		return stackErr.Error(fmt.Errorf("check account search index returned status %s: %s", existsRes.Status(), readBody(existsRes.Body)))
	}

	body, err := json.Marshal(accountSearchIndexDefinition())
	if err != nil {
		return stackErr.Error(fmt.Errorf("marshal account search index definition failed: %w", err))
	}

	createReq := esapi.IndicesCreateRequest{
		Index: p.index,
		Body:  bytes.NewReader(body),
	}
	createRes, err := createReq.Do(ctx, p.client)
	if err != nil {
		return stackErr.Error(fmt.Errorf("create account search index failed: %w", err))
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		payload := readBody(createRes.Body)
		if createRes.StatusCode == http.StatusBadRequest && strings.Contains(payload, "resource_already_exists_exception") {
			return nil
		}
		return stackErr.Error(fmt.Errorf("create account search index returned status %s: %s", createRes.Status(), payload))
	}
	return nil
}

func accountSearchIndexDefinition() map[string]interface{} {
	return map[string]interface{}{
		"settings": map[string]interface{}{
			"analysis": map[string]interface{}{
				"analyzer": map[string]interface{}{
					"account_text": map[string]interface{}{
						"tokenizer": "standard",
						"filter":    []string{"lowercase", "asciifolding"},
					},
				},
			},
		},
		"mappings": map[string]interface{}{
			"dynamic": "strict",
			"properties": map[string]interface{}{
				"id":                map[string]interface{}{"type": "keyword"},
				"email":             map[string]interface{}{"type": "keyword"},
				"display_name":      map[string]interface{}{"type": "text", "analyzer": "account_text", "fields": map[string]interface{}{"keyword": map[string]interface{}{"type": "keyword", "ignore_above": 256}}},
				"username":          map[string]interface{}{"type": "keyword"},
				"avatar_object_key": map[string]interface{}{"type": "keyword", "ignore_above": 1024},
				"status":            map[string]interface{}{"type": "keyword"},
				"email_verified_at": map[string]interface{}{"type": "date"},
				"last_login_at":     map[string]interface{}{"type": "date"},
				"password_changed_at": map[string]interface{}{
					"type": "date",
				},
				"created_at":    map[string]interface{}{"type": "date"},
				"updated_at":    map[string]interface{}{"type": "date"},
				"banned_reason": map[string]interface{}{"type": "keyword", "ignore_above": 512},
				"banned_until":  map[string]interface{}{"type": "date"},
			},
		},
	}
}

func resolveAccountSearchIndex(cfg config.ElasticsearchConfig) string {
	index := strings.TrimSpace(cfg.AccountIndex)
	if index == "" {
		return defaultAccountSearchIndex
	}
	return index
}

func cloneStringPtr(value *string) *string {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func cloneTimePtr(value *time.Time) *time.Time {
	if value == nil {
		return nil
	}
	cloned := value.UTC()
	return &cloned
}
