package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	LoggerLogCountView = &view.View{
		Name:        "go.logger/log",
		Description: "Count of logs by level",
		Measure:     stat.StatLoggerLogCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagLoggerLogLevel},
	}
)

func overrideLoggerView() {

}

func initLoggerView() []*view.View {
	return []*view.View{
		LoggerLogCountView,
	}
}
