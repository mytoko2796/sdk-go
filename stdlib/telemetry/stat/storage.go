package stat

import "go.opencensus.io/stats"

var (
	StatStorageLatency = stats.Float64(`go.storage/latency`, `latency of storage request`, stats.UnitMilliseconds)
)