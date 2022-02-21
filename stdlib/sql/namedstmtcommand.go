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

type commandnamedstmt struct {
	ctx        context.Context
	name       string
	nstmt      *sqlx.NamedStmt
	tagMutator []tags.Mutator
}

type CommandNamedStmt interface {
	Close() error
	// Select selects from db where the result is mapped to dest
	Select(name string, dest interface{}, arg interface{}) error
	// QueryRow returns single row from db
	QueryRow(name string, arg interface{}) (*sqlx.Row, error)
	// Query returns multiple rows from db
	Query(name string, arg interface{}) (*sqlx.Rows, error)
	// Get returns query result from db and map them to dest
	Get(name string, dest interface{}, arg interface{}) error
	// Exec query against db
	Exec(name string, arg interface{}) (sql.Result, error)
	// MustExec query against db
	MustExec(name string, arg interface{}) (sql.Result, error)
	// Unsafe
	Unsafe() CommandNamedStmt
}

func initNamedStmt(ctx context.Context, name string, mutator []tags.Mutator, nstmt *sqlx.NamedStmt) CommandNamedStmt {
	return &commandnamedstmt{
		ctx:        ctx,
		name:       name,
		nstmt:      nstmt,
		tagMutator: mutator,
	}
}

// getDBWithMutatedContext returns mutated context based on node selection
// it returns db object that will be responsible in query execution
func (x *commandnamedstmt) getNamedStmtWithMutatedContext(queryName string) (context.Context, error) {
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

func (x *commandnamedstmt) Close() error {
	return x.nstmt.Close()
}

func (x *commandnamedstmt) Unsafe() CommandNamedStmt {
	x.nstmt = x.nstmt.Unsafe()
	return x
}

func (x *commandnamedstmt) Select(name string, dest interface{}, arg interface{}) error {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.nstmt.SelectContext(ctx, dest, arg)
}

func (x *commandnamedstmt) QueryRow(name string, arg interface{}) (*sqlx.Row, error) {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.nstmt.QueryRowxContext(ctx, arg), nil
}

func (x *commandnamedstmt) Query(name string, arg interface{}) (*sqlx.Rows, error) {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.nstmt.QueryxContext(ctx, arg)
}

func (x *commandnamedstmt) Get(name string, dest interface{}, arg interface{}) error {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return err
	}
	return x.nstmt.GetContext(ctx, dest, arg)
}
func (x *commandnamedstmt) Exec(name string, arg interface{}) (sql.Result, error) {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.nstmt.ExecContext(ctx, arg)
}

func (x *commandnamedstmt) MustExec(name string, arg interface{}) (sql.Result, error) {
	ctx, err := x.getNamedStmtWithMutatedContext(name)
	if err != nil {
		return nil, err
	}
	return x.nstmt.MustExecContext(ctx, arg), nil
}
