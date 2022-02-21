package tag

import (
	tags "go.opencensus.io/tag"
)

var (
	TagMongoError, _  = tags.NewKey(`go.mongo.error`)
	TagMongoMethod, _ = tags.NewKey(`go.mongo.method`)
	TagMongoStatus, _ = tags.NewKey(`go.mongo.status`)

	TagMongoDB, _    = tags.NewKey(`go.mongo.db`)
	TagMongoHost, _  = tags.NewKey(`go.mongo.host`)
	TagMongoQuery, _ = tags.NewKey(`go.mongo.query`)
)
