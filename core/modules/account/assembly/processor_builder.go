package assembly

import (
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/account/application/projection/processor"
	accountrepo "wechat-clone/core/modules/account/infra/persistent/repository"
	accountes "wechat-clone/core/modules/account/infra/projection/elasticsearch"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
)

func buildSearchProjectionProcessor(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	searchProjection, err := accountes.NewAccountSearchProjection(cfg.ElasticsearchConfig, appCtx.GetElasticsearchClient())
	if err != nil {
		return nil, stackErr.Error(err)
	}

	searchRepository, err := accountes.NewAccountSearchRepository(cfg.ElasticsearchConfig, appCtx.GetElasticsearchClient())
	if err != nil {
		return nil, stackErr.Error(err)
	}
	accountReadRepo := accountrepo.NewAccountRepoImpl(appCtx.GetDB(), appCtx.GetCache(), true, nil, searchRepository)

	return processor.NewProcessor(cfg, accountReadRepo, searchProjection)
}

func buildProjectionRuntime(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	searchProjection, err := buildSearchProjectionProcessor(cfg, appCtx)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return modruntime.NewComposite(searchProjection), nil
}
