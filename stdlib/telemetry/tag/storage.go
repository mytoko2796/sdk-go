package tag

import (
tags "go.opencensus.io/tag"
)

var (
	TagStorageDriver, _ = tags.NewKey(`go.storage.driver`)
	TagStorageMethod, _  = tags.NewKey(`go.storage.method`)
	TagStorageStatusCode, _ = tags.NewKey(`go.storage.status`)
)

