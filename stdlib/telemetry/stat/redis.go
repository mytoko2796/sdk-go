package stat

import (
	"go.opencensus.io/stats"
)

var (
	StatRedisLatency     = stats.Float64(`go.redis/latency`, `The latency of calls in milliseconds`, stats.UnitMilliseconds)
	StatRedisIdleConn    = stats.Int64(`go.redis/connections/idle`, `Count of idle connections in the pool`, stats.UnitDimensionless)
	StatRedisStaleConn   = stats.Int64(`go.redis/connections/stale`, `Count of stale connections removed from the pool`, stats.UnitDimensionless)
	StatRedisTotalConn   = stats.Int64(`go.redis/connections/total`, `Count of total connections in the pool`, stats.UnitDimensionless)
	StatRedisTimeoutConn = stats.Int64(`go.redis/connections/timeout`, `Number of times a wait timeout occurred`, stats.UnitDimensionless)
	StatRedisConnMisses  = stats.Int64(`go.redis/connections/misses`, `Number of times free connection was NOT found in the pool`, stats.UnitDimensionless)
	StatRedisConnHits    = stats.Int64(`go.redis/connections/hits`, `Number of times free connection was found in the pool`, stats.UnitDimensionless)
)
