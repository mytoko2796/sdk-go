package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	tag "github.com/mytoko2796/sdk-go/stdlib/telemetry/tag"
	tags "go.opencensus.io/tag"
	octrace "go.opencensus.io/trace"
)

type commandstmt struct {
	ctx        context.Context
	name       string
	stmt       *sqlx.Stmt
	tagMutator []tags.Mutator
}

type CommandStmt interface {
	Close() error
	// Select selects from db where the result is mapped to dest
	Select(name string, dest interface{}, args ...interface{}) error
	// QueryRow returns single row from db
	QueryRow(name string, args ...interface{}) (*sqlx.Row, error)
	// Query returns multiple rows from db
	Query(name string, args ...interface{}) (*sqlx.Rows, error)
	// Get returns query result from db and map them to dest
	Get(name string, dest interface{}, args ...interface{}) error
	// Exec query against db
	Exec(name string, args ...interface{}) (sql.Result, error)
	// MustExec query against db
	MustExec(name string, args ...interface{}) (sql.Result, error)
	// Unsafe
	Unsafe() CommandStmt
}

func initStmt(ctx context.Context, name string, mutator []tags.Mutator, stmt *sqlx.Stmt) CommandStmt {
	return &commandstmt{
		ctx:        ctx,
		name:       name,
		stmt:       stmt,
		tagMutator: mutator,
	}
}

// getDBWithMutatedContext returns mutated context based on node selection
// it returns db object that will be responsible in query execution
func (x *commandstmt) getStmtWithMutatedContext(queryName string) (context.Context, error) {
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

func (x *commandstmt) Close() error {
	return x.stmt.Close()
}

func (x *commandstmt) Unsafe() CommandStmt {
	x.stmt = x.stmt.Unsafe()
	return x
}

func (x *commandstmt) Select(name string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.stmt.SelectContext(ctx, dest, args...)
}

func (x *commandstmt) QueryRow(name string, args ...interface{}) (*sqlx.Row, error) {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.stmt.QueryRowxContext(ctx, args...), nil
}

func (x *commandstmt) Query(name string, args ...interface{}) (*sqlx.Rows, error) {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.stmt.QueryxContext(ctx, args...)
}

func (x *commandstmt) Get(name string, dest interface{}, args ...interface{}) error {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.stmt.GetContext(ctx, dest, args...)
}
func (x *commandstmt) Exec(name string, args ...interface{}) (sql.Result, error) {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.stmt.ExecContext(ctx, args...)
}

func (x *commandstmt) MustExec(name string, args ...interface{}) (sql.Result, error) {
	ctx, err := x.getStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.stmt.MustExecContext(ctx, args...), nil
}
