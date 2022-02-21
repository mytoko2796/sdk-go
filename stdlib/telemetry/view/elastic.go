package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	ElasticRoundtripLatencyView = &view.View{
		Name:        "go.elastic/roundtrip/latency",
		Description: "HTTP Elastic Client Request Latency",
		Measure:     stat.StatElasticRoundtripLatency,
		Aggregation: DefaultDBMsDistribution,
		TagKeys:     []tags.Key{tag.TagElasticHost, tag.TagElasticMethod, tag.TagElasticURL},
	}

	ElasticReceivedBytesView = &view.View{
		Name:        "go.elastic/response/size",
		Description: "HTTP Elastic Client Response Received Size",
		Measure:     stat.StatElasticReceivedBytes,
		Aggregation: DefaultHTTPSizeDistribution,
		TagKeys:     []tags.Key{tag.TagElasticHost, tag.TagElasticMethod, tag.TagElasticURL},
	}

	ElasticSentBytesView = &view.View{
		Name:        "go.elastic/request/size",
		Description: "HTTP Elastic Client Request Sent Size",
		Measure:     stat.StatElasticSentBytes,
		Aggregation: DefaultHTTPSizeDistribution,
		TagKeys:     []tags.Key{tag.TagElasticHost, tag.TagElasticMethod, tag.TagElasticURL},
	}

	ElasticResponseCountByStatusCodeView = &view.View{
		Name:        "go.elastic/response/status",
		Description: "HTTP Elastic Client Response Count By Status Code",
		Measure:     stat.StatElasticRoundtripLatency,
		Aggregation: ochttp.ServerResponseCountByStatusCode.Aggregation,
		TagKeys:     []tags.Key{tag.TagElasticHost, tag.TagElasticMethod, tag.TagElasticURL, tag.TagElasticStatusCode},
	}
)

func overrideElasticView() {

}

func initElasticView() []*view.View {
	return []*view.View{
		ElasticRoundtripLatencyView,
		ElasticReceivedBytesView,
		ElasticSentBytesView,
		ElasticResponseCountByStatusCodeView,
	}
}
