package stat

import "go.opencensus.io/stats"

var (
	StatLoggerLogCount = stats.Int64(`go.logger/count`, `Current number of log messages by level`, stats.UnitDimensionless)
)
