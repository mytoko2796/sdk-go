package stat

import (
	"go.opencensus.io/stats"
)

var (
	//client stats
	StatElasticSentBytes        = stats.Int64(`go.elastic/sent/size`, `Total number of bytes transmitted to Elastic Search`, stats.UnitBytes)
	StatElasticReceivedBytes    = stats.Int64(`go.elastic/received/size`, `Total number of bytes received by Elastic Search`, stats.UnitBytes)
	StatElasticRoundtripLatency = stats.Float64(`go.elastic/query/latency`, `The latency of calls in milliseconds`, stats.UnitMilliseconds)
)
