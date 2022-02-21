package stat

import (
	"go.opencensus.io/stats"
)

var (
	// Kafka Producer Stat
	StatOrchProducerLatency = stats.Float64(`go.bpm/producer/latency`, `Latency of writing requests to bpm brokers in milliseconds`, stats.UnitMilliseconds)
	StatOrchProducerCount   = stats.Int64(`go.bpm/producer/count`, `Total number of writing requests to bpm brokers`, stats.UnitDimensionless)
	StatOrchProducerBytes   = stats.Int64(`go.bpm/producer/size`, `Total number of bytes transmitted to bpm brokers`, stats.UnitBytes)

	// Kafka Consumer Stat
	StatOrchConsumerLatency = stats.Float64(`go.bpm/consumer/latency`, `Latency of request per job type (handler) in milliseconds`, stats.UnitMilliseconds)
	StatOrchConsumerCount   = stats.Int64(`go.bpm/consumer/count`, `Total number of requests received from bpm brokers`, stats.UnitDimensionless)
	StatOrchConsumerBytes   = stats.Int64(`go.bpm/consumer/size`, `Total number of request bytes received from bpm brokers`, stats.UnitBytes)
)
