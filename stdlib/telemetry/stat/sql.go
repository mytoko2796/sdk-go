package stat

import "contrib.go.opencensus.io/integrations/ocsql"

var (
	StatSQLLatency             = ocsql.MeasureLatencyMs
	StatSQLMeasureWaitDuration = ocsql.MeasureWaitDuration

	StatSQLMeasureOpenConnection   = ocsql.MeasureOpenConnections
	StatSQLMeasureIdleConnection   = ocsql.MeasureIdleConnections
	StatSQLMeasureActiveConnection = ocsql.MeasureActiveConnections
	StatSQLMeasureWaitCount        = ocsql.MeasureWaitCount
	StatSQLMeasureIdleClosed       = ocsql.MeasureIdleClosed
	StatSQLMeasureLifetimeClosed   = ocsql.MeasureLifetimeClosed
)
