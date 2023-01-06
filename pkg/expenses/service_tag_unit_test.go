//go:build unit

package expenses

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
)

type DBCaseSuccess struct {
	insertWasCalled     bool
	searchByIdWasCalled bool
	updateWasCalled     bool
	searchAllWasCalled  bool
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

func (db *DBCaseSuccess) SearchById(id int64) (*ExpensesResponse, error) {
	db.searchByIdWasCalled = true
	resp := &ExpensesResponse{
		Id:     id,
		Title:  "mockTitle",
		Amount: 10,
		Note:   "mockNote",
		Tags:   []string{"mockTags"},
	}
	return resp, nil
}

func (db *DBCaseSuccess) Update(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	db.updateWasCalled = true
	resp := &ExpensesResponse{
		Id:     id,
		Title:  req.Title,
		Amount: req.Amount,
		Note:   req.Note,
		Tags:   req.Tags,
	}
	return resp, nil
}

func (db *DBCaseSuccess) SearchAll() ([]ExpensesResponse, error) {
	db.searchAllWasCalled = true
	resp := []ExpensesResponse{
		{
			Id:     1,
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		},
		{
			Id:     2,
			Title:  "mockTitle2",
			Amount: 9000,
			Note:   "mockNote2",
			Tags:   []string{"mockTags2"},
		},
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
	insertWasCalled     bool
	searchByIdWasCalled bool
	updateWasCalled     bool
	searchAllWasCalled  bool
}

func (db *DBCaseError) Insert(req ExpensesRequest) (*ExpensesResponse, error) {
	db.insertWasCalled = true
	return nil, &Err{}
}

func (db *DBCaseError) SearchById(id int64) (*ExpensesResponse, error) {
	db.searchByIdWasCalled = true
	return nil, &Err{}
}

func (db *DBCaseError) Update(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	db.updateWasCalled = true
	return nil, &Err{}
}

func (db *DBCaseError) SearchAll() ([]ExpensesResponse, error) {
	db.searchAllWasCalled = true
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

func TestSearchExpensesById(t *testing.T) {
	t.Run("should return ExpensesResponse when no error that storage.SearchById()", func(t *testing.T) {
		storage := &DBCaseSuccess{}
		log := logrus.New()
		service := NewService(storage, log)
		id := int64(43)

		resp, err := service.SearchExpensesById(id)

		assert.Equal(t, true, storage.searchByIdWasCalled)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.Equal(t, id, resp.Id)
		assert.NotEmpty(t, resp.Title)
		assert.NotEmpty(t, resp.Amount)
		assert.NotEmpty(t, resp.Note)
		assert.NotEmpty(t, resp.Tags)
	})

	t.Run("should return error when  error that storage.SearchById()", func(t *testing.T) {
		storage := &DBCaseError{}
		log := logrus.New()
		service := NewService(storage, log)
		id := int64(43)

		resp, err := service.SearchExpensesById(id)

		assert.Equal(t, true, storage.searchByIdWasCalled)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})
}

func TestUpdateExpenses(t *testing.T) {
	t.Run("should return ExpensesResponse when no error that storage.Update()", func(t *testing.T) {
		storage := &DBCaseSuccess{}
		log := logrus.New()
		service := NewService(storage, log)
		id := int64(43)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}

		resp, err := service.UpdateExpenses(id, req)

		assert.Equal(t, true, storage.updateWasCalled)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.Equal(t, id, resp.Id)
		assert.Equal(t, req.Title, resp.Title)
		assert.Equal(t, req.Amount, resp.Amount)
		assert.Equal(t, req.Note, resp.Note)
		assert.Equal(t, req.Tags, resp.Tags)
	})

	t.Run("should return error when  error that storage.Update()", func(t *testing.T) {
		storage := &DBCaseError{}
		log := logrus.New()
		service := NewService(storage, log)
		id := int64(43)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}

		resp, err := service.UpdateExpenses(id, req)

		assert.Equal(t, true, storage.updateWasCalled)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})
}

func TestSearchExpensesAll(t *testing.T) {
	t.Run("should return ExpensesResponse when no error that storage.SearchAll()", func(t *testing.T) {
		storage := &DBCaseSuccess{}
		log := logrus.New()
		service := NewService(storage, log)

		resp, err := service.SearchExpensesAll()

		assert.Equal(t, true, storage.searchAllWasCalled)
		assert.NotNil(t, resp)
		assert.Nil(t, err)
		assert.Len(t, resp, 2)
	})

	t.Run("should return error when  error that storage.SearchAll()", func(t *testing.T) {
		storage := &DBCaseError{}
		log := logrus.New()
		service := NewService(storage, log)

		resp, err := service.SearchExpensesAll()

		assert.Equal(t, true, storage.searchAllWasCalled)
		assert.NotNil(t, err)
		assert.Nil(t, resp)
	})
}
