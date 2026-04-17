package projection

import (
	roomprojection "wechat-clone/core/modules/room/application/projection"
	"wechat-clone/core/shared/config"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gocql/gocql"
)

func NewCassandraTimelineProjector(cfg config.CassandraConfig, session *gocql.Session) (roomprojection.ServingProjector, error) {
	store, err := NewCassandraProjectionStore(cfg, session)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return store, nil
}
