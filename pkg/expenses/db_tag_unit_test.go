//go:build unit

package expenses

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	t.Run("should insert success when no error", func(t *testing.T) {
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, req.Title, req.Amount, req.Note, pq.Array(req.Tags))
		get := mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id, title, amount, note, tags"))
		get.ExpectQuery().WithArgs(req.Title, req.Amount, req.Note, pq.Array(req.Tags)).WillReturnRows(row)

		dataMgmt := New(db)
		result, err := dataMgmt.Insert(req)

		assert.Nil(t, err)
		assert.Equal(t, req.Title, result.Title)
		assert.Equal(t, req.Amount, result.Amount)
		assert.Equal(t, req.Note, result.Note)
		assert.Equal(t, req.Tags, result.Tags)
		assert.Equal(t, int64(1), result.Id)
	})

	t.Run("should return error when error", func(t *testing.T) {
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		get := mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id, title, amount, note, tags"))
		get.ExpectQuery().WithArgs(req.Title, req.Amount, req.Note, pq.Array(req.Tags)).WillReturnError(&pq.Error{Message: "error connection db"})

		dataMgmt := New(db)
		result, err := dataMgmt.Insert(req)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestSearchById(t *testing.T) {
	t.Run("should search success when no error", func(t *testing.T) {
		id := int64(1)
		mockData := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(id, mockData.Title, mockData.Amount, mockData.Note, pq.Array(mockData.Tags))
		get := mock.ExpectPrepare(regexp.QuoteMeta("select id, title, amount, note, tags from expenses where id = $1"))
		get.ExpectQuery().WithArgs(id).WillReturnRows(row)

		dataMgmt := New(db)
		result, err := dataMgmt.SearchById(id)

		assert.Nil(t, err)
		assert.Equal(t, mockData.Title, result.Title)
		assert.Equal(t, mockData.Amount, result.Amount)
		assert.Equal(t, mockData.Note, result.Note)
		assert.Equal(t, mockData.Tags, result.Tags)
		assert.Equal(t, id, result.Id)
	})

	t.Run("should return error when error", func(t *testing.T) {
		id := int64(1)
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		get := mock.ExpectPrepare(regexp.QuoteMeta("select id, title, amount, note, tags from expenses where id = $1"))
		get.ExpectQuery().WithArgs(id).WillReturnError(&pq.Error{Message: "error connection db"})

		dataMgmt := New(db)
		result, err := dataMgmt.SearchById(id)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("should update success when no error", func(t *testing.T) {
		id := int64(2)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(id, req.Title, req.Amount, req.Note, pq.Array(req.Tags))
		get := mock.ExpectPrepare(regexp.QuoteMeta("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id, title, amount, note, tags"))
		get.ExpectQuery().WithArgs(req.Title, req.Amount, req.Note, pq.Array(req.Tags), id).WillReturnRows(row)

		dataMgmt := New(db)
		result, err := dataMgmt.Update(id, req)

		assert.Nil(t, err)
		assert.Equal(t, req.Title, result.Title)
		assert.Equal(t, req.Amount, result.Amount)
		assert.Equal(t, req.Note, result.Note)
		assert.Equal(t, req.Tags, result.Tags)
		assert.Equal(t, id, result.Id)
	})

	t.Run("should return error when error", func(t *testing.T) {
		id := int64(2)
		req := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		get := mock.ExpectPrepare(regexp.QuoteMeta("UPDATE expenses SET title = $1, amount = $2, note = $3, tags = $4 WHERE id = $5 RETURNING id, title, amount, note, tags"))
		get.ExpectQuery().WithArgs(req.Title, req.Amount, req.Note, pq.Array(req.Tags), id).WillReturnError(&pq.Error{Message: "error connection db"})

		dataMgmt := New(db)
		result, err := dataMgmt.Update(id, req)

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}

func TestSearchAll(t *testing.T) {
	t.Run("should search success when no error", func(t *testing.T) {
		mockData := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"})
		row.AddRow(1, mockData.Title, mockData.Amount, mockData.Note, pq.Array(mockData.Tags))
		row.AddRow(2, mockData.Title, mockData.Amount, mockData.Note, pq.Array(mockData.Tags))
		get := mock.ExpectPrepare(regexp.QuoteMeta("select id, title, amount, note, tags from expenses"))
		get.ExpectQuery().WithArgs().WillReturnRows(row)

		dataMgmt := New(db)
		result, err := dataMgmt.SearchAll()

		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)
	})

	t.Run("should return error when error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		defer db.Close()
		assert.NoError(t, err)
		get := mock.ExpectPrepare(regexp.QuoteMeta("select id, title, amount, note, tags from expenses"))
		get.ExpectQuery().WithArgs().WillReturnError(&pq.Error{Message: "error connection db"})

		dataMgmt := New(db)
		result, err := dataMgmt.SearchAll()

		assert.NotNil(t, err)
		assert.Nil(t, result)
	})
}
