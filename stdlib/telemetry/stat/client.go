package stat

import (
	"go.opencensus.io/plugin/ochttp"
)

var (
	//client stats
	StatClientSentBytes        = ochttp.ClientSentBytes
	StatClientReceivedBytes    = ochttp.ClientReceivedBytes
	StatClientRoundtripLatency = ochttp.ClientRoundtripLatency
)
