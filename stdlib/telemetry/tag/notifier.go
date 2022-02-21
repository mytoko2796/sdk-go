package tag

import (
tags "go.opencensus.io/tag"
)

var (
	TagMailerDriver, _ = tags.NewKey(`go.mailer.driver`)
	TagMailerMethod, _  = tags.NewKey(`go.mailer.method`)
	TagMailerTemplate, _  = tags.NewKey(`go.mailer.template`)
	TagMailerStatusCode, _ = tags.NewKey(`go.mailer.status`)

	TagSMSDriver,_ = tags.NewKey(`go.sms.driver`)
	TagSMSStatusCode, _ = tags.NewKey(`go.sms.status`)
)

