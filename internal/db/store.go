package db

import (
	"context"
	"database/sql"

	"github.com/caring/go-packages/pkg/errors"
	_ "github.com/caring/go-packages/pkg/uuid"
	// anonymous import so package exports are not exposed
	_ "github.com/go-sql-driver/mysql"
)

type ctxKey struct{}

var txCtxKey = ctxKey{}

// Store represents a connection and a collection
// of statements that we will use to interface with
// a backing store
type Store struct {
	db    *sql.DB
	stmts map[string]*sql.Stmt
}

// NewStore will give a pointer to a MySQL instance ready to run queries against.
func NewStore(dataSourceName string) (*Store, error) {
	unprepared := statements

	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	stmts, err := prepareStmts(db, unprepared)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	err = db.Ping()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	s := Store{
		db:    db,
		stmts: stmts,
	}

	return &s, nil
}

// prepareStmts will attempt to prepare each unprepared
// query on the database. If one fails, the function returns
// with an error.
func prepareStmts(db *sql.DB, unprepared map[string]string) (map[string]*sql.Stmt, error) {
	prepared := map[string]*sql.Stmt{}
	for k, v := range unprepared {
		stmt, err := db.Prepare(v)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		prepared[k] = stmt
	}
	return prepared, nil
}

// Close will close the connection to the underlying database
func (s *Store) Close() error {
	err := s.db.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Ping will check the connection to the underlying database
func (s *Store) Ping(ctx context.Context) error {
	if err = s.db.PingContext(ctx); err != nil {
    return err
  }
  return nil
}

// GetTx initializes a db transaction
func (s *Store) GetTx() (*sql.Tx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return tx, nil
}

// ToCtx stores a sql.Tx within a context
func ToCtx(ctx context.Context, tx *sql.Tx) context.Context {
	return context.WithValue(ctx, txCtxKey, tx)
}

// FromCtx extracts a sql.Tx from a context which has been stored by this package,
// returns a error if no tx is present
func FromCtx(ctx context.Context) (*sql.Tx, error) {
	val := ctx.Value(txCtxKey)
	if val != nil {
		return val.(*sql.Tx), nil
	}
	return nil, errors.New("No *sql.Tx present in context")
}
