package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	//server views
	ViewRequestCount              = ochttp.ServerRequestCountView
	ViewRequestCountByMethod      = ochttp.ServerRequestCountByMethod
	ViewServerRequestBytes        = ochttp.ServerRequestBytesView
	ViewResponseBytes             = ochttp.ServerResponseBytesView
	ViewResponseCountByStatusCode = ochttp.ServerResponseCountByStatusCode
	ViewLatency                   = ochttp.ServerLatencyView
)

func overrideServerView() {
	ViewRequestCount.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute, tag.TagHTTPStatus}
	ViewRequestCountByMethod.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute}
	ViewServerRequestBytes.Aggregation = DefaultHTTPSizeDistribution
	ViewServerRequestBytes.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute}
	ViewResponseBytes.Aggregation = DefaultHTTPSizeDistribution
	ViewResponseBytes.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute, tag.TagHTTPStatus}
	ViewResponseCountByStatusCode.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute, tag.TagHTTPStatus}
	ViewLatency.Aggregation = DefaultHTTPMsDistribution
	ViewLatency.TagKeys = []tags.Key{tag.TagHTTPMethod, tag.TagKeyServerRoute, tag.TagHTTPStatus}
}

func initServerView() []*view.View {
	return []*view.View{
		ViewRequestCount,
		ViewRequestCountByMethod,
		ViewServerRequestBytes,
		ViewResponseBytes,
		ViewResponseCountByStatusCode,
		ViewLatency,
	}
}
