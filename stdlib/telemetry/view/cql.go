package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	CQLClientLatencyView = &view.View{
		Name:        "go.cql/client/latency",
		Description: "The distribution of latencies of various calls in milliseconds",
		Measure:     stat.StatCQLLatency,
		Aggregation: DefaultDBMsDistribution,
		TagKeys:     []tags.Key{tag.TagCQLHost, tag.TagCQLKeyspace, tag.TagCQLQuery, tag.TagCQLMethod, tag.TagCQLStatus},
	}

	CQLClientAttemptView = &view.View{
		Name:        "go.cql/client/attempt",
		Description: "sum of attempt for various calls",
		Measure:     stat.StatCQLAttempt,
		Aggregation: view.Sum(),
		TagKeys:     []tags.Key{tag.TagCQLHost, tag.TagCQLKeyspace, tag.TagCQLQuery, tag.TagCQLMethod, tag.TagCQLStatus},
	}

	CQLClientConnectLatencyView = &view.View{
		Name:        "go.cql/connect/latency",
		Description: "The distribution of connect latencies of various calls in milliseconds",
		Measure:     stat.StatCQLConnectLatency,
		Aggregation: DefaultConnectMsDistribution,
		TagKeys:     []tags.Key{tag.TagCQLHost},
	}
)

func overrideCQLView() {

}

func initCQLView() []*view.View {
	return []*view.View{
		CQLClientLatencyView,
		CQLClientAttemptView,
		CQLClientConnectLatencyView,
	}
}
