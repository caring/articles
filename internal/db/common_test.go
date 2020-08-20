package db

import (
	"github.com/DATA-DOG/go-sqlmock"
)

// NewTestDB creates a testable store instance with a mocked sql driver
// and provides a test utility for making assertions against and setting query
// response values.
//
// queries executed against the test DB do not interact in any way with a real DB
// see https://github.com/DATA-DOG/go-sqlmock
func NewTestDB(stmts map[string]string) (*Store, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	for _, s := range stmts {
		mock.ExpectPrepare(s)
	}

	prepared, err := prepareStmts(db, stmts)
	if err != nil {
		return nil, nil, err
	}

	s := Store{
		db:    db,
		stmts: prepared,
	}

	return &s, mock, nil
}
