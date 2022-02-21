package tag

import tags "go.opencensus.io/tag"

var (
	TagCQLHost, _     = tags.NewKey(`go.cql.host`)
	TagCQLKeyspace, _ = tags.NewKey(`go.cql.keyspace`)
	TagCQLQuery, _    = tags.NewKey(`go.cql.query`)
	TagCQLMethod, _   = tags.NewKey(`go.cql.method`)
	TagCQLStatus, _   = tags.NewKey(`go.cql.status`)
)
