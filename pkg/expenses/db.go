package expenses

import (
	"database/sql"

	"github.com/lib/pq"
)

type DataMgmt struct {
	dataMgmt *sql.DB
}

func New(d *sql.DB) *DataMgmt {
	return &DataMgmt{d}
}

func (mgmt DataMgmt) Insert(req ExpensesRequest) (*ExpensesResponse, error) {
	stmt, err := mgmt.dataMgmt.Prepare("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id, title, amount, note, tags")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Title, req.Amount, req.Note, pq.Array(req.Tags))

	result := &ExpensesResponse{}
	err = row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (mgmt DataMgmt) SearchById(id int64) (*ExpensesResponse, error) {
	stmt, err := mgmt.dataMgmt.Prepare("select id, title, amount, note, tags from expenses where id = $1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows := stmt.QueryRow(id)

	result := &ExpensesResponse{}
	err = rows.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (mgmt DataMgmt) Update(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	stmt, err := mgmt.dataMgmt.Prepare("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id, title, amount, note, tags")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(req.Title, req.Amount, req.Note, pq.Array(req.Tags), id)

	result := &ExpensesResponse{}
	err = row.Scan(&result.Id, &result.Title, &result.Amount, &result.Note, pq.Array(&result.Tags))
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (mgmt DataMgmt) SearchAll() ([]ExpensesResponse, error) {
	stmt, err := mgmt.dataMgmt.Prepare("select id, title, amount, note, tags from expenses")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []ExpensesResponse{}
	for rows.Next() {
		exp := &ExpensesResponse{}
		err = rows.Scan(&exp.Id, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		if err != nil {
			return nil, err
		}
		result = append(result, *exp)
	}

	return result, nil
}
