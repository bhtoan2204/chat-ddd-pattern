package projection

import (
	"context"

	"wechat-clone/core/modules/room/infra/projection/cassandra/views"
	sharedcassandra "wechat-clone/core/shared/infra/cassandra"
	"wechat-clone/core/shared/pkg/stackErr"

	"github.com/gocql/gocql"
)

const projectionMigrationSource = "file://migration/cassandra/room_projection"

func runProjectionMigrations(ctx context.Context, session *gocql.Session, tables views.ProjectionTableNames) error {
	tool := sharedcassandra.NewMigrateTool()
	if err := tool.MigrateFromSource(ctx, session, tables.SchemaMigrations, projectionMigrationSource); err != nil {
		return stackErr.Error(err)
	}
	return nil
}
