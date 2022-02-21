package sql

import (
	"context"
	"database/sql"
	"fmt"

	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	"github.com/jmoiron/sqlx"
	tags "go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"
)

type commandtx struct {
	ctx        context.Context
	name       string
	tx         *sqlx.Tx
	tagMutator []tags.Mutator
}

type CommandTx interface {
	// Commit tx
	Commit() error
	// Rollback tx
	Rollback() error
	// Rebind rebinds query to db
	Rebind(query string) string
	// BindNamed
	BindNamed(query string, arg interface{}) (string, []interface{}, error)
	// Select selects from db where the result is mapped to dest
	Select(name string, query string, dest interface{}, args ...interface{}) error
	// Prepare prepares statement to db
	Prepare(name string, query string) (CommandStmt, error)
	// Prepare prepares named statement to db
	PrepareNamed(name string, query string) (CommandNamedStmt, error)
	// QueryRow returns single row from db
	QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error)
	// Query returns multiple rows from db
	Query(name string, query string, args ...interface{}) (*sqlx.Rows, error)
	// NamedQuery take named query as arg and returns multiple rows from db
	NamedQuery(name string, query string, arg interface{}) (*sqlx.Rows, error)
	// Get returns query result from db and map them to dest
	Get(name string, query string, dest interface{}, args ...interface{}) error
	// Exec query against db
	Exec(name string, query string, args ...interface{}) (sql.Result, error)
	// NamedExec exec named query against db
	NamedExec(name string, query string, arg interface{}) (sql.Result, error)
	// Stmt
	Stmt(name string, stmt *sqlx.Stmt) CommandStmt
	// Unsafe
	Unsafe() CommandTx
}

var _ CommandTx = (*commandtx)(nil)

func initTx(ctx context.Context, name string, mutator []tags.Mutator, tx *sqlx.Tx, opts *sql.TxOptions) CommandTx {
	return &commandtx{
		ctx:        ctx,
		name:       name,
		tx:         tx,
		tagMutator: mutator,
	}
}

// getTxWithMutatedContext returns mutated context based on node selection
// it returns db object that will be responsible in query execution
func (x *commandtx) getTxWithMutatedContext(queryName string) (context.Context, error) {
	var (
		err         error
		tagMutators []tags.Mutator
	)
	span := octrace.FromContext(x.ctx)
	span.AddAttributes(octrace.StringAttribute(fmt.Sprintf("%s:%s", tag.TagSQLQuery.Name(), x.name), queryName))
	tagMutators = append(tagMutators, x.tagMutator...)
	tagMutators = append(tagMutators, tags.Upsert(tag.TagSQLQuery, fmt.Sprintf("%s:%s:%s", tag.TagSQLQuery.Name(), x.name, queryName)))
	ctx, err := tags.New(x.ctx, tagMutators...)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func (x *commandtx) Commit() error {
	return x.tx.Commit()
}

func (x *commandtx) Unsafe() CommandTx {
	x.tx = x.tx.Unsafe()
	return x
}

func (x *commandtx) Rollback() error {
	return x.tx.Rollback()
}

// Rebind rebinds query to db
func (x *commandtx) Rebind(query string) string {
	return x.tx.Rebind(query)
}

func (x *commandtx) BindNamed(query string, arg interface{}) (string, []interface{}, error) {
	return x.tx.BindNamed(query, arg)
}

// QueryRow returns single row from db
func (x *commandtx) QueryRow(name string, query string, args ...interface{}) (*sqlx.Row, error) {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.tx.QueryRowxContext(ctx, query, args...), nil
}

// Query returns multiple rows from db
func (x *commandtx) Query(name string, query string, args ...interface{}) (*sqlx.Rows, error) {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.tx.QueryxContext(ctx, query, args...)
}

// Select selects from db where the result is mapped to dest
func (x *commandtx) Select(name string, query string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.tx.SelectContext(ctx, dest, query, args...)
}

// Prepare prepares statement to db
func (x *commandtx) Prepare(name string, query string) (CommandStmt, error) {
	stmt, err := x.tx.PreparexContext(x.ctx, query)
	if err != nil {
		return nil, err
	}
	return initStmt(x.ctx, name, x.tagMutator, stmt), nil
}

// Prepare preprares named statement to db
func (x *commandtx) PrepareNamed(name string, query string) (CommandNamedStmt, error) {
	nstmt, err := x.tx.PrepareNamedContext(x.ctx, query)
	if err != nil {
		return nil, err
	}
	return initNamedStmt(x.ctx, name, x.tagMutator, nstmt), nil
}

// Get returns query result from db and map them to dest
func (x *commandtx) Get(name string, query string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.tx.GetContext(ctx, dest, query, args...)
}

// NamedQuery takes named query as arg and returns multiple rows from db
func (x *commandtx) NamedQuery(name string, query string, arg interface{}) (*sqlx.Rows, error) {
	_, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.tx.NamedQuery(query, arg)
}

// Exec query against db
func (x *commandtx) Exec(name string, query string, args ...interface{}) (sql.Result, error) {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return nil, err
	}

	return x.tx.ExecContext(ctx, query, args...)
}

// NamedExec exec named query against db
func (x *commandtx) NamedExec(name string, query string, arg interface{}) (sql.Result, error) {
	ctx, err := x.getTxWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.tx.NamedExecContext(ctx, query, arg)
}

// NamedStmt name stmt query against db
func (x *commandtx) Stmt(name string, stmt *sqlx.Stmt) CommandStmt {
	return initStmt(x.ctx, name, x.tagMutator, stmt)
}
