package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	MongoClientLatencyView = &view.View{
		Name:        "mongo/client/latency",
		Description: "The latency of the various calls",
		Measure:     stat.StatMongoLatency,
		Aggregation: DefaultDBMsDistribution,
		TagKeys:     []tags.Key{tag.TagMongoHost, tag.TagMongoDB, tag.TagMongoMethod, tag.TagMongoStatus, tag.TagMongoQuery},
	}

	MongoClientCallsView = &view.View{
		Name: "mongo/client/calls", Description: "The various calls",
		Measure:     stat.StatMongoLatency,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagMongoHost, tag.TagMongoDB, tag.TagMongoMethod, tag.TagMongoStatus, tag.TagMongoQuery},
	}
)

func overrideMongoView() {

}

func initMongoView() []*view.View {
	return []*view.View{
		MongoClientLatencyView,
		MongoClientCallsView,
	}
}
