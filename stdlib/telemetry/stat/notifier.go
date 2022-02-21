package stat

import "go.opencensus.io/stats"

var (
	StatMailerCount = stats.Int64(`go.mailer/count`, `Current number of email sent`, stats.UnitDimensionless)
	StatMailerLatency = stats.Float64(`go.mailer/latency`, `latency of send email request`, stats.UnitMilliseconds)

	StatSMSCount = stats.Int64(`go.sms/count`, `Current number of sms sent`, stats.UnitDimensionless)
	StatSMSLatency = stats.Float64(`go.sms/latency`, `latency of send sms request`, stats.UnitMilliseconds)
)