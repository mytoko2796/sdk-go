package tag

import tags "go.opencensus.io/tag"

var (
	TagRedisDriver, _ = tags.NewKey(`go.redis.driver`)
	TagRedisMethod, _ = tags.NewKey(`go.redis.method`)
	TagRedisHost, _   = tags.NewKey(`go.redis.host`)
)
