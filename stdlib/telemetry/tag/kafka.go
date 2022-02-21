package tag

import tags "go.opencensus.io/tag"

var (
	TagKafkaSubTopic, _   = tags.NewKey(`go.kafka.subscribed.topic`)
	TagKafkaPubTopic, _   = tags.NewKey(`go.kafka.published.topic`)
	TagKafkaPartition, _  = tags.NewKey(`go.kafka.partition`)
	TagKafkaStatusCode, _ = tags.NewKey(`go.kafka.status.code`)
)
