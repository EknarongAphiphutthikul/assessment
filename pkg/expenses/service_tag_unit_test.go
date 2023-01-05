//go:build unit

package expenses

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
)

type DBCaseSuccess struct {
	insertWasCalled bool
}

func (db *DBCaseSuccess) Insert(req ExpensesRequest) (*ExpensesResponse, error) {
	db.insertWasCalled = true
	resp := &ExpensesResponse{
		Id:     1,
		Title:  req.Title,
		Amount: req.Amount,
		Note:   req.Note,
		Tags:   req.Tags,
	}
	return resp, nil
}

type Err struct {
	msg string
}

func (e Err) Error() string {
	return "mock error" + e.msg
}

type DBCaseError struct {
	insertWasCalled bool
}

func (db *DBCaseError) Insert(req ExpensesRequest) (*ExpensesResponse, error) {
	db.insertWasCalled = true
	return nil, &Err{}
}

func TestAddExpenses(t *testing.T) {
	t.Run("should return ExpensesResponse when no error that storage.Insert()", func(t *testing.T) {
		storage := &DBCaseSuccess{}
		log := logrus.New()
		service := NewService(storage, log)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}

		resp, err := service.AddExpenses(req)

		assert.Equal(t, true, storage.insertWasCalled)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.NotEqual(t, int64(0), resp.Id)
		assert.Equal(t, req.Title, resp.Title)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Note, resp.Note)
		assert.Equal(t, req.Tags, resp.Tags)
	})

	t.Run("should return error when  error that storage.Insert()", func(t *testing.T) {
		storage := &DBCaseError{}
		log := logrus.New()
		service := NewService(storage, log)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}

		resp, err := service.AddExpenses(req)

		assert.Equal(t, true, storage.insertWasCalled)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})
}
