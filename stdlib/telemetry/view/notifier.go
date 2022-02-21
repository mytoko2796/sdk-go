package view

import (
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/stat"
	"github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"go.opencensus.io/stats/view"
	tags "go.opencensus.io/tag"
)

var (

	MailerCountView = &view.View{
		Name:        "go.mailer/count",
		Description: "Count of processed email",
		Measure:     stat.StatMailerCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagMailerDriver, tag.TagMailerStatusCode, tag.TagMailerMethod},
	}

	MailerLatencyView = &view.View{
		Name:        "go.mailer/latency",
		Description: "Send Email Latency",
		Measure:     stat.StatMailerLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagMailerDriver, tag.TagMailerMethod, tag.TagSMSStatusCode},
	}

	SMSCountView = &view.View{
		Name:        "go.sms/count",
		Description: "Count of processed SMS",
		Measure:     stat.StatSMSCount,
		Aggregation: view.Count(),
		TagKeys:     []tags.Key{tag.TagSMSDriver, tag.TagSMSStatusCode},
	}

	SMSLatencyView = &view.View{
		Name:        "go.sms/latency",
		Description: "Send SMS Latency",
		Measure:     stat.StatSMSLatency,
		Aggregation: DefaultHTTPMsDistribution,
		TagKeys:     []tags.Key{tag.TagSMSDriver, tag.TagSMSStatusCode},
	}
)

func overrideNotifierView() {
}

func initNotifierView() []*view.View {
	return []*view.View{
		MailerCountView,
		MailerLatencyView,
		SMSCountView,
		SMSLatencyView,
	}
}
