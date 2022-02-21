package view

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	ViewSQLLatency              = ocsql.SQLClientLatencyView
	ViewSQLClientCalls          = ocsql.SQLClientCallsView
	ViewSQLClientOpenConns      = ocsql.SQLClientOpenConnectionsView
	ViewSQLClientIdleConns      = ocsql.SQLClientIdleConnectionsView
	ViewSQLClientActiveConns    = ocsql.SQLClientActiveConnectionsView
	ViewSQLClientWaitCount      = ocsql.SQLClientWaitCountView
	ViewSQLClientWaitDuration   = ocsql.SQLClientWaitDurationView
	ViewSQLClientIdleClosed     = ocsql.SQLClientIdleClosedView
	ViewSQLClientLifetimeClosed = ocsql.SQLClientLifetimeClosedView
)

func overrideSQLView() {
	ViewSQLLatency.Aggregation = DefaultDBMsDistribution
	ViewSQLLatency.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
		tag.TagSQLQuery,
		tag.TagSQLGoSQLMethod,
		tag.TagSQLGoSQLStatus,
	}

	ViewSQLClientCalls.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
		tag.TagSQLQuery,
		tag.TagSQLGoSQLMethod,
		tag.TagSQLGoSQLStatus,
	}
	ViewSQLClientOpenConns.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientIdleConns.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientActiveConns.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientWaitCount.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientWaitDuration.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientIdleClosed.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
	ViewSQLClientLifetimeClosed.TagKeys = []tags.Key{
		tag.TagSQLHost,
		tag.TagSQLDriver,
		tag.TagSQLDB,
	}
}

func initSQLView() []*view.View {
	return []*view.View{
		ViewSQLLatency,
		ViewSQLClientCalls,
		ViewSQLClientOpenConns,
		ViewSQLClientIdleConns,
		ViewSQLClientActiveConns,
		ViewSQLClientWaitCount,
		ViewSQLClientWaitDuration,
		ViewSQLClientIdleClosed,
		ViewSQLClientLifetimeClosed,
	}
}
