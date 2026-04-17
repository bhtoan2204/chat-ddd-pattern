package assembly

import (
	appCtx "wechat-clone/core/context"
	"wechat-clone/core/modules/room/application/projection/processor"
	roomCassandra "wechat-clone/core/modules/room/infra/projection/cassandra"
	roomElasticsearch "wechat-clone/core/modules/room/infra/projection/elasticsearch"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"
	modruntime "wechat-clone/core/shared/runtime"
)

func buildServingProjectionProcessor(cfg *config.Config, appCtx *appCtx.AppContext) (modruntime.Module, error) {
	servingProjector, err := roomCassandra.NewCassandraTimelineProjector(cfg.CassandraConfig, appCtx.GetCassandraSession())
	if err != nil {
		return nil, stackErr.Error(err)
	}

	searchIndexer, err := roomElasticsearch.NewElasticsearchMessageIndexer(cfg.ElasticsearchConfig, appCtx.GetElasticsearchClient())
	if err != nil {
		return nil, stackErr.Error(err)
	}

	return processor.NewProcessor(cfg, servingProjector, searchIndexer)
}
