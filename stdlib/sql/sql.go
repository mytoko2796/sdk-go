// sql package implements jmoiron/sqlx as an interface to database. This package currently supports
// postgres and mysql only. All methods require context to be passed to help terminating request whenever
// request time exceeds request context deadline. All queries should be indexed/ named as this packages
// monitor the sql activities. e.g. query latencies, connections, etc.
// See Details :
//	https://www.github.com/jmoiron/sqlx
//
package sql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"contrib.go.opencensus.io/integrations/ocsql"
	"github.com/cenkalti/backoff"
	mysql "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	errors "github.com/mytoko2796/sdk-go/stdlib/error"
	log "github.com/mytoko2796/sdk-go/stdlib/logger"
	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	tags "go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"
)

const (
	PGSQL string = `postgres`
	MYSQL string = `mysql`
)

const (
	infoSQL string = `SQL:`
	_OK     string = "[OK]"
	_FAILED string = "[FAILED]"

	defaultMaxConnectTimeout = 15 * time.Second
)

// Default errors from sql packages
var (
	ErrNoRows   = sql.ErrNoRows
	ErrTxDone   = sql.ErrTxDone
	ErrConnDone = sql.ErrConnDone
)

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 sql.NullInt64

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 sql.NullFloat64

// NullBool is an alias for sql.NullBool data type
type NullBool sql.NullBool

// NullTime is an alias for mysql.NullTime data type
type NullTime mysql.NullTime

// NullString is an alias for sql.NullString data type
type NullString sql.NullString

// SQL
type SQL interface {
	Leader() Command
	Follower() Command
	Driver() string
	// Stop stopping sql recorder and close all db connections
	Stop()
}

// sqlxImpl
type sqlxImpl struct {
	endOnce  *sync.Once
	logger   log.Logger
	leader   Command
	follower Command
	recorder []*recorder
	opt      Options
}

// Options
type Options struct {
	Enabled  bool
	Driver   string
	Leader   Config
	Follower Config
}

// Config
type Config struct {
	// Host URI
	Host string
	// Port database port
	Port int
	// DB database name
	DB string
	// User
	User string
	// Password
	Password string
	// SSL
	SSL bool
	// ConnOptions
	ConnOptions ConnOptions
	// TraceOptions
	TraceOptions TraceOptions
	// mockdb
	MockDB *sql.DB
}

var In = sqlx.In

// ConnOptions
type ConnOptions struct {
	MaxLifeTime time.Duration
	MaxIdle     int
	MaxOpen     int
}

// TraceOptions
type TraceOptions struct {
	// Enabled
	Enabled bool
	// RecordPeriod for recorders to poll db stats
	RecordPeriod time.Duration
	// Available Tracing Configurations
	AllowRoot    bool
	Ping         bool
	RowsNext     bool
	RowsClose    bool
	RowsAffected bool
	LastInsertID bool
	Query        bool
	QueryParams  bool
}

// Init initialize SQL Object and starts sql activity recorder. All established connections and
// must be stopped when application terminates to prevent any dangling connections.
func Init(logger log.Logger, opt Options) SQL {
	if !opt.Enabled {
		return nil
	}

	sql := &sqlxImpl{
		endOnce:  &sync.Once{},
		logger:   logger,
		opt:      opt,
		recorder: nil,
	}

	sql.initDB()
	sql.StartRecorder()
	return sql
}

func (x *sqlxImpl) Leader() Command {
	return x.leader
}

func (x *sqlxImpl) Follower() Command {
	return x.follower
}

func (x *sqlxImpl) Driver() string {
	return x.opt.Driver
}

// initDB initialize db if the follower db exists then x.follower will be set as new object that represents
// follower db connection pools. Otherwise, x.follower will be equal to x.leader
func (x *sqlxImpl) initDB() {
	db, tagMutators, err := x.connect(true)
	if err != nil {
		err = errors.Wrap(err, errInitSQLDBLeader)
		x.logger.Fatal(err)
	}
	x.leader = initCommand(db, tagMutators)

	x.logger.Info(_OK, infoSQL, fmt.Sprintf("[LEADER] driver=%s db=%s @%s:%v ssl=%v", x.opt.Driver, x.opt.Leader.DB, x.opt.Leader.Host, x.opt.Leader.Port, x.opt.Leader.SSL))

	if x.isFollowerEnabled() {
		db, tagMutators, err = x.connect(false)
		if err != nil {
			err = errors.Wrap(err, errInitSQLDBFollower)
			x.logger.Fatal(err)
		}
		x.logger.Info(_OK, infoSQL, fmt.Sprintf("[FOLLOWER] driver=%s db=%s @%s:%v ssl=%v", x.opt.Driver, x.opt.Follower.DB, x.opt.Follower.Host, x.opt.Follower.Port, x.opt.Leader.SSL))
		x.follower = initCommand(db, tagMutators)
	}
	x.follower = x.leader
}

// Connect return db connection pool with tags mutator which will be used in stats and tracings.
func (x *sqlxImpl) connect(toLeader bool) (*sqlx.DB, []tags.Mutator, error) {
	conf := x.opt.Leader
	if !toLeader {
		conf = x.opt.Follower
	}

	// see if mock db object is passed
	if toLeader {
		if x.opt.Leader.MockDB != nil {
			return sqlx.NewDb(x.opt.Leader.MockDB, x.opt.Driver), nil, nil
		}
	} else {
		if x.opt.Follower.MockDB != nil {
			return sqlx.NewDb(x.opt.Follower.MockDB, x.opt.Driver), nil, nil
		}
	}

	dbHost := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	tagMutators := []tags.Mutator{
		tags.Upsert(tag.TagSQLHost, dbHost),
		tags.Upsert(tag.TagSQLDriver, x.opt.Driver),
		tags.Upsert(tag.TagSQLDB, conf.DB),
	}

	var err error
	trace := conf.TraceOptions
	driverName, err := ocsql.Register(
		x.opt.Driver,
		ocsql.WithPing(trace.Ping),
		ocsql.WithAllowRoot(trace.AllowRoot),
		ocsql.WithRowsNext(trace.RowsNext),
		ocsql.WithRowsClose(trace.RowsClose),
		ocsql.WithRowsAffected(trace.RowsAffected),
		ocsql.WithLastInsertID(trace.LastInsertID),
		ocsql.WithQuery(trace.Query),
		ocsql.WithQueryParams(trace.QueryParams),
		ocsql.WithDisableErrSkip(true),
		ocsql.WithDefaultAttributes(
			octrace.StringAttribute(tag.TagSQLHost.Name(), dbHost),
			octrace.StringAttribute(tag.TagSQLDriver.Name(), x.opt.Driver),
			octrace.StringAttribute(tag.TagSQLDB.Name(), conf.DB),
		))
	if err != nil {
		return nil, tagMutators, errors.WrapWithCode(err, EcodeBadSQLDriver, errSQL, _FAILED)
	}

	uri, err := x.getURI(conf)
	if err != nil {
		return nil, tagMutators, errors.WrapWithCode(err, EcodeBadSQLURI, errSQL, _FAILED)
	}

	db, err := sql.Open(driverName, uri)
	if err != nil {
		return nil, tagMutators, errors.WrapWithCode(err, EcodeBadSQLOpen, errSQL, _FAILED)
	}

	sqlOpen := func() error {
		return db.PingContext(context.Background())
	}

	bo := backoff.NewExponentialBackOff()
	bo.MaxElapsedTime = defaultMaxConnectTimeout
	bo.MaxInterval = defaultMaxConnectTimeout
	bo.Multiplier = 1.5
	bo.RandomizationFactor = 0.5
	err = backoff.RetryNotify(
		sqlOpen,
		bo,
		backoff.Notify(func(err error, duration time.Duration) {
			if err != nil {
				err = errors.WrapWithCode(err, EcodeBadSQLConnection, errSQL, `[RETRY]`)
				x.logger.Error(err, duration)
			}
		}))
	if err != nil {
		return nil, tagMutators, errors.WrapWithCode(err, EcodeBadSQLConnection, errSQL, ``)
	}

	sqlxDB := sqlx.NewDb(db, x.opt.Driver)
	sqlxDB.SetMaxOpenConns(conf.ConnOptions.MaxOpen)
	sqlxDB.SetMaxIdleConns(conf.ConnOptions.MaxIdle)
	sqlxDB.SetConnMaxLifetime(conf.ConnOptions.MaxLifeTime)

	return sqlxDB, tagMutators, nil
}

// getURI returns formatted uri for particular database implementation
// currently only supports postgres and mysql
func (x *sqlxImpl) getURI(conf Config) (string, error) {
	switch x.opt.Driver {
	case PGSQL:
		ssl := `disable`
		if conf.SSL {
			ssl = `require`
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", conf.Host, conf.Port, conf.User, conf.Password, conf.DB, ssl), nil

	case MYSQL:
		ssl := `false`
		if conf.SSL {
			ssl = `true`
		}
		return fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?tls=%s", conf.User, conf.Password, conf.Host, conf.Port, conf.DB, ssl), nil

	default:
		return "", errors.New(`DB Driver is not supported [%s]`, x.opt.Driver)
	}
}
