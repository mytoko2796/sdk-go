package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	StorageLatencyView = &view.View{
		Name:        "go.storage/latency",
		Description: "Storage Request Latency",
		Measure:     stat.StatStorageLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagStorageDriver, tag.TagStorageMethod, tag.TagStorageStatusCode},
	}
)

func overrideStorageView() {
}

func initStorageView() []*view.View {
	return []*view.View{
		StorageLatencyView,
	}
}
