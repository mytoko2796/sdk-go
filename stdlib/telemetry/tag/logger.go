package tag

import (
	tags "go.opencensus.io/tag"
)

var (
	//client tag
	TagLoggerLogLevel, _ = tags.NewKey(`log.level`)
)
