package expenses

import (
	"net/http"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
)

type Storage interface {
	Insert(req ExpensesRequest) (*ExpensesResponse, error)
	SearchById(id int64) (*ExpensesResponse, error)
	Update(id int64, req ExpensesRequest) (*ExpensesResponse, error)
	SearchAll() ([]ExpensesResponse, error)
}

type Service struct {
	log     common.Log
	storage Storage
}

func NewService(s Storage, l common.Log) *Service {
	return &Service{storage: s, log: l}
}

func (s Service) AddExpenses(req ExpensesRequest) (*ExpensesResponse, error) {
	resp, err := s.storage.Insert(req)
	if err != nil {
		s.log.Errorf("Insert Expenses Error : %s", err)
		return nil, &common.Error{Code: http.StatusInternalServerError, Desc: "Insert Expenses Error", OriginalError: err}
	}
	return resp, nil
}

func (s Service) SearchExpensesById(id int64) (*ExpensesResponse, error) {
	resp, err := s.storage.SearchById(id)
	if err != nil {
		s.log.Errorf("Search Expenses By Id Error : %s", err)
		return nil, &common.Error{Code: http.StatusInternalServerError, Desc: "Search Expenses By Id Error", OriginalError: err}
	}
	return resp, nil
}

func (s Service) UpdateExpenses(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	resp, err := s.storage.Update(id, req)
	if err != nil {
		s.log.Errorf("Update Expenses Error : %s", err)
		return nil, &common.Error{Code: http.StatusInternalServerError, Desc: "Update Expenses Error", OriginalError: err}
	}
	return resp, nil
}

func (s Service) SearchExpensesAll() ([]ExpensesResponse, error) {
	resp, err := s.storage.SearchAll()
	if err != nil {
		s.log.Errorf("Search Expenses All Error: %s", err)
		return nil, &common.Error{Code: http.StatusInternalServerError, Desc: "Search Expenses All Error", OriginalError: err}
	}
	return resp, nil
}
