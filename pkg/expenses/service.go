package expenses

import (
	"net/http"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
)

type Storage interface {
	Insert(req ExpensesRequest) (*ExpensesResponse, error)
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
