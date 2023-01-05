//go:build unit

package expenses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ServiceSuccess struct {
	addExpensesWasCalled bool
}

func (s *ServiceSuccess) AddExpenses(req ExpensesRequest) (*ExpensesResponse, error) {
	s.addExpensesWasCalled = true
	resp := &ExpensesResponse{
		Id:     1,
		Title:  req.Title,
		Amount: req.Amount,
		Note:   req.Note,
		Tags:   req.Tags,
	}
	return resp, nil
}

type ServiceError struct {
	addExpensesWasCalled bool
	statusCodeError      int
}

func (s *ServiceError) AddExpenses(req ExpensesRequest) (*ExpensesResponse, error) {
	s.addExpensesWasCalled = true
	return nil, &common.Error{Code: s.statusCodeError}
}

func TestAddExpensesHandler(t *testing.T) {
	t.Run("should return http status code = 201 and ExpensesResponse when no error that service.AddExpenses()", func(t *testing.T) {
		// Arrange
		reqBody := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceSuccess{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err = handler.AddExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			resp := &ExpensesResponse{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NotEqual(t, int64(0), resp.Id)
			assert.Equal(t, reqBody.Title, resp.Title)
			assert.Equal(t, reqBody.Amount, resp.Amount)
			assert.Equal(t, reqBody.Note, resp.Note)
			assert.Equal(t, reqBody.Tags, resp.Tags)
		}
	})

	t.Run("should return http status code = 500 when error that service.AddExpenses()", func(t *testing.T) {
		// Arrange
		reqBody := ExpensesRequest{
			Title:  "mockTitle",
			Amount: 10,
			Note:   "mockNote",
			Tags:   []string{"mockTags"},
		}
		body, err := json.Marshal(reqBody)
		if err != nil {
			assert.Fail(t, "json marshal error")
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err = handler.AddExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, service.statusCodeError, rec.Code)
		}
	})
}
