package expenses

import (
	"database/sql"

	"github.com/lib/pq"
)

type DataMgmt struct {
	sql *sql.DB
}

func New(s *sql.DB) *DataMgmt {
	return &DataMgmt{s}
}

func (mgmt DataMgmt) Insert(req ExpensesRequest) (int64, error) {
	row := mgmt.sql.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", req.Title, req.Amount, req.Note, pq.Array(req.Tags))

	var id int64
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
