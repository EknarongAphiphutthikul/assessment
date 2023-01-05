package expenses

import (
	"database/sql"

	"github.com/lib/pq"
)

type IDataMgmt interface {
	QueryRow(query string, args ...any) *sql.Row
}

type DataMgmt struct {
	dataMgmt IDataMgmt
}

func New(d IDataMgmt) *DataMgmt {
	return &DataMgmt{d}
}

func (mgmt DataMgmt) Insert(req ExpensesRequest) (*ExpensesResponse, error) {
	row := mgmt.dataMgmt.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id, title, amount, note, tags", req.Title, req.Amount, req.Note, pq.Array(req.Tags))

	result := &ExpensesResponse{}
	err := row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return nil, err
	}
	return result, nil
}
