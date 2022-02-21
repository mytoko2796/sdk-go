package stat

import (
	"go.opencensus.io/stats"
)

var (
	// see internal metrics exported by librdkafka
	// https://github.com/edenhill/librdkafka/blob/master/STATISTICS.md

	// Librdkafka Kafka Producer Internal Queue Stat
	StatKafkaProducerQueueCount = stats.Int64(`go.kafka/producer/qcount`, `Current number of messages in producer queues`, stats.UnitDimensionless)
	StatKafkaProducerQueueSize  = stats.Int64(`go.kafka/producer/qsize`, `Current total size of messages in producer queues`, stats.UnitBytes)

	// Kafka Producer Stat
	StatKafkaProducerLatency = stats.Float64(`go.kafka/producer/latency`, `Latency of writing requests to kafka brokers in milliseconds`, stats.UnitMilliseconds)
	StatKafkaProducerCount   = stats.Int64(`go.kafka/producer/count`, `Total number of writing requests to Kafka brokers`, stats.UnitDimensionless)
	StatKafkaProducerBytes   = stats.Int64(`go.kafka/producer/size`, `Total number of bytes transmitted to Kafka brokers`, stats.UnitBytes)

	// Kafka Consumer Stat
	StatKafkaConsumerLatency = stats.Float64(`go.kafka/consumer/latency`, `Latency of request per topic (handler) in milliseconds`, stats.UnitMilliseconds)
	StatKafkaConsumerCount   = stats.Int64(`go.kafka/consumer/count`, `Total number of requests received from Kafka brokers`, stats.UnitDimensionless)
	StatKafkaConsumerBytes   = stats.Int64(`go.kafka/consumer/size`, `Total number of request bytes received from Kafka brokers`, stats.UnitBytes)
)
