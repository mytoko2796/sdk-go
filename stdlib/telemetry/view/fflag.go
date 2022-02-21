package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	FFLagFlagLatencyView = &view.View{
		Name:        "go.fflag/roundtrip/latency",
		Description: "HTTP Client Request Evaluation Latency by Flag",
		Measure:     stat.StatFFlagLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagFFlagFlagID, tag.TagFFlagFeature, tag.TagHTTPClientStatus, tag.TagFFlagCacheMiss},
	}

	FFlagFlagCompletedCount = &view.View{
		Name:        "go.fflag/flag/completed_count",
		Description: "HTTP Client Request Evaluation Latency by Flag, HTTP Status",
		Measure:     stat.StatFFlagLatency,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagFFlagFlagID, tag.TagFFlagFeature, tag.TagHTTPClientStatus, tag.TagFFlagCacheMiss},
	}

	FFlagSegmentView = &view.View{
		Name:        "go.fflag/segment",
		Description: "Total count of segment allocation by flag",
		Measure:     stat.StatFFlagCall,
		Aggregation: view.Sum(),
		TagKeys:     []tags.Key{tag.TagFFlagFlagID, tag.TagFFlagFeature, tag.TagFFlagSegmentID},
	}

	FFlagVariantView = &view.View{
		Name:        "go.fflag/variant",
		Description: "Total count of given variants by flag",
		Measure:     stat.StatFFlagCall,
		Aggregation: view.Sum(),
		TagKeys:     []tags.Key{tag.TagFFlagFlagID, tag.TagFFlagFeature, tag.TagFFlagVariantID},
	}
)

func overrideFFlagView() {

}

func initFFlagView() []*view.View {
	return []*view.View{
		FFLagFlagLatencyView,
		FFlagFlagCompletedCount,
		FFlagSegmentView,
		FFlagVariantView,
	}
}
