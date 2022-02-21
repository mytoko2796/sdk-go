package tag

import tags "go.opencensus.io/tag"

var (
	TagElasticHost, _       = tags.NewKey(`go.elastic.host`)
	TagElasticMethod, _     = tags.NewKey(`go.elastic.method`)
	TagElasticURL, _        = tags.NewKey(`go.elastic.url`)
	TagElasticPort, _       = tags.NewKey(`go.elastic.port`)
	TagElasticStatusCode, _ = tags.NewKey(`go.elastic.status`)
)
