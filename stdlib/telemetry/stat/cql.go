package stat

import (
	"go.opencensus.io/stats"
)

var (
	StatCQLLatency        = stats.Float64(`go.cql/query/latency`, `The latency of calls in milliseconds`, stats.UnitMilliseconds)
	StatCQLAttempt        = stats.Int64(`go.cql/query/attempt`, `Number of attempt for each calls`, stats.UnitDimensionless)
	StatCQLConnectLatency = stats.Float64(`go.cql/connect/latency`, `The latency of connects in milliseconds`, stats.UnitMilliseconds)
)
