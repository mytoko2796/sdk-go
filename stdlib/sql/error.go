package sql

import errors "github.com/mytoko2796/sdk-go/stdlib/error"

// Ecode defines package internal error code
const (
	// SQL Error Codes
	EcodeBadSQLDriver = errors.Code(iota)
	EcodeInvalidQueryRegistration
	EcodeQueryNotFound
	EcodeBadSQLURI
	EcodeBadSQLOpen
	EcodeBadSQLConnection
)

const (
	errSQL                  string = `%sSQL Error`
	errSQLQueryRegistration string = `Query should not be nil`
	errSQLQueryExist        string = `Cannot register query with the same name %s`
	errInitSQLDBLeader      string = `Cannot init SQL database leader`
	errInitSQLDBFollower    string = `Cannot init SQL database follower`
	errTelemetrySetContext  string = `Cannot mutate Tag Value`
)
