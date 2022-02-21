package tag

import (
	"contrib.go.opencensus.io/integrations/ocsql"
	tags "go.opencensus.io/tag"
)

var (
	TagSQLGoSQLError  = ocsql.GoSQLError
	TagSQLGoSQLMethod = ocsql.GoSQLMethod
	TagSQLGoSQLStatus = ocsql.GoSQLStatus

	TagSQLDriver, _ = tags.NewKey(`go.sql.driver`)
	TagSQLDB, _     = tags.NewKey(`go.sql.db`)
	TagSQLHost, _   = tags.NewKey(`go.sql.host`)
	TagSQLQuery, _  = tags.NewKey(`go.sql.query`)
)
