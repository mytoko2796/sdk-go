package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	KafkaProducerQueueCountView = &view.View{
		Name:        "go.kafka/producer/queue_count",
		Description: "Current number of messages in producer queues",
		Measure:     stat.StatKafkaProducerQueueCount,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{},
	}

	KafkaProducerQueueSizeView = &view.View{
		Name:        "go.kafka/producer/queue_size",
		Description: "go.kafka/producer/queue_size`, `Current total size of messages in producer queues",
		Measure:     stat.StatKafkaProducerQueueSize,
		Aggregation: view.LastValue(),
		TagKeys:     []tags.Key{},
	}

	KafkaProducerLatencyView = &view.View{
		Name:        "go.kafka/producer/latency",
		Description: "Latency of requests to kafka brokers in milliseconds",
		Measure:     stat.StatKafkaProducerLatency,
		Aggregation: DefaultMessagingMsDistribution,
		TagKeys:     []tags.Key{tag.TagKafkaPubTopic},
	}

	KafkaProducerCountView = &view.View{
		Name:        "go.kafka/producer/count",
		Description: "Total number of requests to Kafka brokers",
		Measure:     stat.StatKafkaProducerCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagKafkaPubTopic, tag.TagKafkaStatusCode},
	}

	KafkaProducerSizeView = &view.View{
		Name:        "go.kafka/producer/size",
		Description: "Total number of bytes transmitted to Kafka brokers",
		Measure:     stat.StatKafkaProducerBytes,
		Aggregation: DefaultMessagingSizeDistribution,
		TagKeys:     []tags.Key{tag.TagKafkaPubTopic},
	}

	KafkaConsumerLatencyView = &view.View{
		Name:        "go.kafka/consumer/latency",
		Description: "Latency of process topic handling in milliseconds",
		Measure:     stat.StatKafkaConsumerLatency,
		Aggregation: DefaultMessagingMsDistribution,
		TagKeys:     []tags.Key{tag.TagKafkaSubTopic, tag.TagKafkaStatusCode},
	}

	KafkaConsumerCountView = &view.View{
		Name:        "go.kafka/consumer/count",
		Description: "Total number of request received by consumers from Kafka brokers",
		Measure:     stat.StatKafkaConsumerLatency,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagKafkaSubTopic, tag.TagKafkaStatusCode},
	}

	KafkaConsumerSizeView = &view.View{
		Name:        "go.kafka/consumer/size",
		Description: "Total number of bytes received by consumers from Kafka brokers",
		Measure:     stat.StatKafkaConsumerLatency,
		Aggregation: DefaultMessagingSizeDistribution,
		TagKeys:     []tags.Key{tag.TagKafkaSubTopic, tag.TagKafkaStatusCode},
	}
)

func overrideKafkaView() {

}

func initKafkaView() []*view.View {
	return []*view.View{
		KafkaProducerQueueCountView,
		KafkaProducerQueueSizeView,
		KafkaProducerLatencyView,
		KafkaProducerCountView,
		KafkaProducerSizeView,
		KafkaConsumerLatencyView,
		KafkaConsumerCountView,
		KafkaConsumerSizeView,
	}
}
