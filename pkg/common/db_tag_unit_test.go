//go:build unit

package common

import (
	"database/sql"
	"testing"
)

type mockDB struct {
	query        string
	lastInsertID int64
	rowsAffected int64
}

func (m *mockDB) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m *mockDB) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func (m *mockDB) Exec(query string, args ...any) (sql.Result, error) {
	m.query = query
	return m, nil
}

func TestPrepareDB(t *testing.T) {
	mock := &mockDB{}
	sqlStm := `CREATE TABLE IF NOT EXISTS expenses (id SERIAL PRIMARY KEY, title TEXT,	amount FLOAT,	note TEXT,	tags TEXT[]	);`

	prepareDB(sqlStm, mock)

	if mock.query != sqlStm {
		t.Error("should have been call db.Exec with query but it not.")
	}
}
