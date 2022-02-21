package stat

import "go.opencensus.io/stats"

var (
	StatMongoLatency = stats.Float64("latency", "The latency in milliseconds", "ms")
)
