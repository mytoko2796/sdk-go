package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (
	OrchProducerLatencyView = &view.View{
		Name:        "go.bpm/producer/latency",
		Description: "Latency of orchestrator producer in milliseconds",
		Measure:     stat.StatOrchProducerLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchBPMMethod},
	}

	OrchProducerCountView = &view.View{
		Name:        "go.bpm/producer/count",
		Description: "Total number of requests to Orchestrator brokers",
		Measure:     stat.StatOrchProducerCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchBPMMethod, tag.TagOrchStatusCode},
	}

	OrchProducerSizeView = &view.View{
		Name:        "go.bpm/producer/size",
		Description: "Total number of bytes transmitted to Orchestrator brokers",
		Measure:     stat.StatOrchProducerBytes,
		Aggregation: DefaultHTTPSizeDistribution,
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchBPMMethod},
	}

	OrchConsumerLatencyView = &view.View{
		Name:        "go.bpm/consumer/latency",
		Description: "Latency of process job handling in milliseconds",
		Measure:     stat.StatOrchConsumerLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchJobType, tag.TagOrchBPMNProcessID, tag.TagOrchWorkflowVersion},
	}

	OrchConsumerCountView = &view.View{
		Name:        "go.bpm/consumer/count",
		Description: "Total number of request received by consumers from Orchestrator brokers",
		Measure:     stat.StatOrchConsumerCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchJobType, tag.TagOrchBPMNProcessID, tag.TagOrchWorkflowVersion, tag.TagOrchStatusCode},
	}

	OrchConsumerSizeView = &view.View{
		Name:        "go.bpm/consumer/size",
		Description: "Total number of bytes received by consumers from Orchestrator brokers",
		Measure:     stat.StatOrchConsumerLatency,
		Aggregation: DefaultHTTPSizeDistribution,
		TagKeys:     []tags.Key{tag.TagOrchEngine, tag.TagOrchJobType, tag.TagOrchBPMNProcessID, tag.TagOrchWorkflowVersion, tag.TagOrchStatusCode},
	}
)

func overrideOrchView() {

}

func initOrchView() []*view.View {
	return []*view.View{
		OrchProducerLatencyView,
		OrchProducerCountView,
		OrchProducerSizeView,
		OrchConsumerLatencyView,
		OrchConsumerCountView,
		OrchConsumerSizeView,
	}
}
