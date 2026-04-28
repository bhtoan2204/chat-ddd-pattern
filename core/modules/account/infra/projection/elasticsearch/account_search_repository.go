package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	accountprojection "wechat-clone/core/modules/account/application/projection"
	"wechat-clone/core/modules/account/domain/entity"
	valueobject "wechat-clone/core/modules/account/domain/value_object"
	accounttypes "wechat-clone/core/modules/account/types"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"

	es8 "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

type accountSearchRepository struct {
	client *es8.Client
	index  string
}

type accountSearchHit struct {
	Source accountDocument `json:"_source"`
}

type accountSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []accountSearchHit `json:"hits"`
	} `json:"hits"`
}

func NewAccountSearchRepository(cfg config.ElasticsearchConfig, client *es8.Client) (accountprojection.SearchRepository, error) {
	if !cfg.Enabled || client == nil {
		return nil, nil
	}

	repo := &accountSearchRepository{
		client: client,
		index:  resolveAccountSearchIndex(cfg),
	}
	if err := repo.ensureIndex(context.Background()); err != nil {
		return nil, stackErr.Error(err)
	}
	return repo, nil
}

func (r *accountSearchRepository) SearchUsers(ctx context.Context, q string, limit, offset int) ([]*entity.Account, int64, error) {
	if r == nil || r.client == nil {
		return []*entity.Account{}, 0, nil
	}
	q = strings.TrimSpace(q)
	if q == "" {
		return []*entity.Account{}, 0, nil
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	body, err := json.Marshal(searchUsersQuery(q, limit, offset))
	if err != nil {
		return nil, 0, stackErr.Error(fmt.Errorf("marshal account search query failed: %w", err))
	}

	req := esapi.SearchRequest{
		Index: []string{r.index},
		Body:  bytes.NewReader(body),
	}
	res, err := req.Do(ctx, r.client)
	if err != nil {
		return nil, 0, stackErr.Error(fmt.Errorf("search account projection failed: %w", err))
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return []*entity.Account{}, 0, nil
	}
	if res.IsError() {
		return nil, 0, stackErr.Error(fmt.Errorf("search account projection returned status %s: %s", res.Status(), readBody(res.Body)))
	}

	var response accountSearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, 0, stackErr.Error(fmt.Errorf("decode account search response failed: %w", err))
	}

	accounts := make([]*entity.Account, 0, len(response.Hits.Hits))
	for _, hit := range response.Hits.Hits {
		account, err := documentToAccount(hit.Source)
		if err != nil {
			return nil, 0, stackErr.Error(err)
		}
		accounts = append(accounts, account)
	}
	return accounts, response.Hits.Total.Value, nil
}

func (r *accountSearchRepository) ensureIndex(ctx context.Context) error {
	projection := &accountSearchProjection{client: r.client, index: r.index}
	return stackErr.Error(projection.ensureIndex(ctx))
}

func searchUsersQuery(q string, limit, offset int) map[string]interface{} {
	return map[string]interface{}{
		"from": offset,
		"size": limit,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{"term": map[string]interface{}{"status": accounttypes.AccountStatusActive.String()}},
				},
				"should": []interface{}{
					map[string]interface{}{"term": map[string]interface{}{"username": map[string]interface{}{"value": q, "boost": 10}}},
					map[string]interface{}{"prefix": map[string]interface{}{"username": map[string]interface{}{"value": q, "boost": 7}}},
					map[string]interface{}{"term": map[string]interface{}{"email": map[string]interface{}{"value": strings.ToLower(q), "boost": 6}}},
					map[string]interface{}{"prefix": map[string]interface{}{"email": map[string]interface{}{"value": strings.ToLower(q), "boost": 5}}},
					map[string]interface{}{"match_phrase_prefix": map[string]interface{}{"display_name": map[string]interface{}{"query": q, "boost": 4}}},
					map[string]interface{}{"match": map[string]interface{}{"display_name": map[string]interface{}{"query": q, "boost": 2}}},
				},
				"minimum_should_match": 1,
			},
		},
		"sort": []interface{}{
			"_score",
			map[string]interface{}{"last_login_at": map[string]interface{}{"order": "desc", "missing": "_last"}},
			map[string]interface{}{"created_at": map[string]interface{}{"order": "desc"}},
		},
	}
}

func documentToAccount(document accountDocument) (*entity.Account, error) {
	email, err := valueobject.NewEmail(document.Email)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	passwordHash, err := valueobject.NewHashedPassword("projection-redacted-password-hash")
	if err != nil {
		return nil, stackErr.Error(err)
	}
	status, err := accounttypes.ParseAccountStatus(document.Status)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return &entity.Account{
		ID:                document.ID,
		Email:             email,
		PasswordHash:      passwordHash,
		DisplayName:       document.DisplayName,
		Username:          cloneStringPtr(document.Username),
		AvatarObjectKey:   cloneStringPtr(document.AvatarObjectKey),
		Status:            status,
		EmailVerifiedAt:   cloneTimePtr(document.EmailVerifiedAt),
		LastLoginAt:       cloneTimePtr(document.LastLoginAt),
		PasswordChangedAt: cloneTimePtr(document.PasswordChangedAt),
		CreatedAt:         document.CreatedAt.UTC(),
		UpdatedAt:         document.UpdatedAt.UTC(),
		BannedReason:      document.BannedReason,
		BannedUntil:       cloneTimePtr(document.BannedUntil),
	}, nil
}

var _ accountprojection.SearchRepository = (*accountSearchRepository)(nil)
