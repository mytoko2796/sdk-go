package view

import (
	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"

	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	//client views
	ViewClientSentBytes        = ochttp.ClientSentBytesDistribution
	ViewClientReceivedBytes    = ochttp.ClientReceivedBytesDistribution
	ViewClientRoundtripLatency = ochttp.ClientRoundtripLatencyDistribution
	ViewClientCompletedCount   = ochttp.ClientCompletedCount
)

func overrideClientView() {
	ViewClientSentBytes.Aggregation = DefaultHTTPSizeDistribution
	ViewClientSentBytes.TagKeys = []tags.Key{
		tag.TagHTTPClientMethod,
		tag.TagHTTPClientStatus,
		tag.TagKeyClientRoute}

	ViewClientReceivedBytes.Aggregation = DefaultHTTPSizeDistribution
	ViewClientReceivedBytes.TagKeys = []tags.Key{
		tag.TagHTTPClientMethod,
		tag.TagHTTPClientStatus,
		tag.TagKeyClientRoute}

	ViewClientRoundtripLatency.Aggregation = DefaultHTTPMsDistribution
	ViewClientRoundtripLatency.TagKeys = []tags.Key{
		tag.TagHTTPClientMethod,
		tag.TagHTTPClientStatus,
		tag.TagKeyClientRoute}

	ViewClientCompletedCount.TagKeys = []tags.Key{
		tag.TagHTTPClientMethod,
		tag.TagHTTPClientStatus,
		tag.TagKeyClientRoute}
}

func initClientView() []*view.View {
	return []*view.View{
		ViewClientSentBytes,
		ViewClientReceivedBytes,
		ViewClientRoundtripLatency,
		ViewClientCompletedCount,
	}
}
