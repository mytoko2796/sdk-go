package stat

import (
	"go.opencensus.io/stats"
)

var (
	//server stats
	StatFFlagCall    = stats.Int64(`go.fflag/calls`, `number of calls`, stats.UnitDimensionless)
	StatFFlagLatency = stats.Float64(`go.fflag/latency`, `latency of http request`, stats.UnitMilliseconds)
)
