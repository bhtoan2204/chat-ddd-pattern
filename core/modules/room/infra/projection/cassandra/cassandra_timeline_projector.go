package projection

import (
	roomprojection "go-socket/core/modules/room/application/projection"
	"go-socket/core/shared/config"
	"go-socket/core/shared/pkg/stackErr"

	"github.com/gocql/gocql"
)

func NewCassandraTimelineProjector(cfg config.CassandraConfig, session *gocql.Session) (roomprojection.ServingProjector, error) {
	store, err := NewCassandraProjectionStore(cfg, session)
	if err != nil {
		return nil, stackErr.Error(err)
	}
	return store, nil
}
