package sql

import (
	"context"
	"database/sql"

	tags "go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"

	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"github.com/jmoiron/sqlx"
)

type Command interface {
	// getStats
	GetStats() sql.DBStats
	// getTagMutator
	GetTagMutator() []tags.Mutator
	// Close all connections
	Close() error
	// Rebind rebinds query to db
	Rebind(query string) string
	// Ping ping to db
	Ping(ctx context.Context) error
	// Select selects from db where the result is mapped to dest
	Select(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error
	// Prepare prepares statement to db
	Prepare(ctx context.Context, name string, query string) (CommandStmt, error)
	// Prepare preprares named statement to db
	PrepareNamed(ctx context.Context, name string, query string) (CommandNamedStmt, error)
	// QueryRow returns single row from db
	QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error)
	// Query returns multiple rows from db
	Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error)
	// NamedQuery take named query as arg and returns multiple rows from db
	NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error)
	// Get returns query result from db and map them to dest
	Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error
	// Exec query against db
	Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error)
	// NamedExec exec named query against db
	NamedExec(ctx context.Context, name string, query string, arg interface{}) (sql.Result, error)
	// BeginTx begin transaction to db
	BeginTx(ctx context.Context, name string, opts *sql.TxOptions) (CommandTx, error)
	// Unsafe
	Unsafe() Command
}

type command struct {
	db         *sqlx.DB
	tagMutator []tags.Mutator
}

func initCommand(db *sqlx.DB, mutator []tags.Mutator) Command {
	return &command{
		db:         db,
		tagMutator: mutator,
	}
}

func (x *command) GetStats() sql.DBStats {
	return x.db.Stats()
}

func (x *command) GetTagMutator() []tags.Mutator {
	return x.tagMutator
}

func (x *command) Close() error {
	return x.db.Close()
}

// getDBWithMutatedContext returns mutated context based on node selection
// it returns db object that will be responsible in query execution
func (x *command) getDBWithMutatedContext(ctx context.Context, queryName string, isTx bool) (context.Context, error) {
	var (
		err         error
		tagMutators []tags.Mutator
	)

	if isTx {
		return ctx, nil
	}

	span := octrace.FromContext(ctx)
	span.AddAttributes(octrace.StringAttribute(tag.TagSQLQuery.Name(), queryName))
	tagMutators = append(tagMutators, x.tagMutator...)
	tagMutators = append(tagMutators, tags.Upsert(tag.TagSQLQuery, queryName))
	ctx, err = tags.New(ctx, tagMutators...)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

// Rebind rebinds query to db
func (x *command) Rebind(query string) string {
	return x.db.Rebind(query)
}

// Unsafe
func (x *command) Unsafe() Command {
	x.db = x.db.Unsafe()
	return x
}

// Ping ping to db
func (x *command) Ping(ctx context.Context) error {
	ctx, err := x.getDBWithMutatedContext(ctx, `ping`, false)
	if err != nil {
		return err
	}
	return x.db.PingContext(ctx)
}

// QueryRow returns single row from db
func (x *command) QueryRow(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Row, error) {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return nil, err
	}
	return x.db.QueryRowxContext(ctx, query, args...), nil
}

// Query returns multiple rows from db
func (x *command) Query(ctx context.Context, name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return nil, err
	}
	return x.db.QueryxContext(ctx, query, args...)
}

// Select selects from db where the result is mapped to dest
func (x *command) Select(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return err
	}
	return x.db.SelectContext(ctx, dest, query, args...)
}

// Prepare prepares statement to db
func (x *command) Prepare(ctx context.Context, name string, query string) (CommandStmt, error) {
	stmt, err := x.db.PreparexContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(ctx, name, x.tagMutator, stmt), nil
}

// Prepare preprares named statement to db
func (x *command) PrepareNamed(ctx context.Context, name string, query string) (CommandNamedStmt, error) {
	nstmt, err := x.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	return initNamedStmt(ctx, name, x.tagMutator, nstmt), nil
}

// Get returns query result from db and map them to dest
func (x *command) Get(ctx context.Context, name string, query string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return err
	}
	return x.db.GetContext(ctx, dest, query, args...)
}

// NamedQuery takes named query as arg and returns multiple rows from db
func (x *command) NamedQuery(ctx context.Context, name string, query string, arg interface{}) (*sqlx.Rows, error) {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return nil, err
	}
	return x.db.NamedQueryContext(ctx, query, arg)
}

// Exec query against db
func (x *command) Exec(ctx context.Context, name string, query string, args ...interface{}) (sql.Result, error) {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return nil, err
	}
	return x.db.ExecContext(ctx, query, args...)
}

// NamedExec exec named query against db
func (x *command) NamedExec(ctx context.Context, name string, query string, arg interface{}) (sql.Result, error) {
	ctx, err := x.getDBWithMutatedContext(ctx, name, false)
	if err != nil {
		return nil, err
	}
	return x.db.NamedExecContext(ctx, query, arg)
}

// BeginTx begin transaction to db
func (x *command) BeginTx(ctx context.Context, name string, opts *sql.TxOptions) (CommandTx, error) {
	tx, err := x.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return initTx(ctx, name, x.tagMutator, tx, opts), nil
}
