//go:build unit

package expenses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

type ServiceSuccess struct {
	addExpensesWasCalled        bool
	searchExpensesByIdWasCalled bool
	updateExpensesWasCalled     bool
	searchExpensesAllWasCalled  bool
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

func (s *ServiceSuccess) SearchExpensesById(id int64) (*ExpensesResponse, error) {
	s.searchExpensesByIdWasCalled = true
	resp := &ExpensesResponse{
		Id:     id,
		Title:  "mockTitle",
		Amount: 10,
		Note:   "mockNote",
		Tags:   []string{"mockTags"},
	}
	return resp, nil
}

func (s *ServiceSuccess) UpdateExpenses(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	s.updateExpensesWasCalled = true
	resp := &ExpensesResponse{
		Id:     id,
		Title:  req.Title,
		Amount: req.Amount,
		Note:   req.Note,
		Tags:   req.Tags,
	}
	return resp, nil
}

func (s *ServiceSuccess) SearchExpensesAll() ([]ExpensesResponse, error) {
	s.searchExpensesAllWasCalled = true
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

type ServiceError struct {
	addExpensesWasCalled        bool
	searchExpensesByIdWasCalled bool
	updateExpensesWasCalled     bool
	searchExpensesAllWasCalled  bool
	statusCodeError             int
}

func (s *ServiceError) AddExpenses(req ExpensesRequest) (*ExpensesResponse, error) {
	s.addExpensesWasCalled = true
	return nil, &common.Error{Code: s.statusCodeError}
}

func (s *ServiceError) SearchExpensesById(id int64) (*ExpensesResponse, error) {
	s.searchExpensesByIdWasCalled = true
	return nil, &common.Error{Code: s.statusCodeError}
}

func (s *ServiceError) UpdateExpenses(id int64, req ExpensesRequest) (*ExpensesResponse, error) {
	s.updateExpensesWasCalled = true
	return nil, &common.Error{Code: s.statusCodeError}
}

func (s *ServiceError) SearchExpensesAll() ([]ExpensesResponse, error) {
	s.searchExpensesAllWasCalled = true
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

func TestSearchExpensesByIdHandler(t *testing.T) {
	t.Run("should return http status code = 200 and ExpensesResponse when no error that service.SearchExpensesById()", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		id := int64(23)
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		service := &ServiceSuccess{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.SearchExpensesById(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			resp := &ExpensesResponse{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, id, resp.Id)
			assert.NotEmpty(t, resp.Title)
			assert.NotEmpty(t, resp.Amount)
			assert.NotEmpty(t, resp.Note)
			assert.NotEmpty(t, resp.Tags)
		}
	})

	t.Run("should return http status code = 500 when error that service.SearchExpensesById()", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		id := int64(23)
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		service := &ServiceError{statusCodeError: http.StatusInternalServerError}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.SearchExpensesById(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, service.statusCodeError, rec.Code)
		}
	})

	t.Run("should return http status code = 400 when not send param :id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceError{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.SearchExpensesById(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestUpdateExpensesHandler(t *testing.T) {
	t.Run("should return http status code = 200 and ExpensesResponse when no error that service.UpdateExpenses()", func(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		id := int64(23)
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		service := &ServiceSuccess{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err = handler.UpdateExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			resp := &ExpensesResponse{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Equal(t, id, resp.Id)
			assert.Equal(t, reqBody.Title, resp.Title)
			assert.Equal(t, reqBody.Amount, resp.Amount)
			assert.Equal(t, reqBody.Note, resp.Note)
			assert.Equal(t, reqBody.Tags, resp.Tags)
		}
	})

	t.Run("should return http status code = 500 when error that service.UpdateExpenses()", func(t *testing.T) {
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
		req := httptest.NewRequest(http.MethodPut, "/expenses", strings.NewReader(string(body)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		id := int64(23)
		c := e.NewContext(req, rec)
		c.SetPath("/:id")
		c.SetParamNames("id")
		c.SetParamValues(strconv.FormatInt(id, 10))

		service := &ServiceError{statusCodeError: http.StatusBadGateway}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err = handler.UpdateExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, service.statusCodeError, rec.Code)
		}
	})

	t.Run("should return http status code = 400 when not send param :id", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/expenses/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceError{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.UpdateExpenses(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		}
	})
}

func TestSearchExpensesAllHandler(t *testing.T) {
	t.Run("should return http status code = 200 and ExpensesResponse when no error that service.SearchExpensesAll()", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceSuccess{}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.SearchExpensesAll(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, http.StatusOK, rec.Code)
			resp := []ExpensesResponse{}
			json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.Len(t, resp, 2)
		}
	})

	t.Run("should return http status code = 500 when error that service.SearchExpensesAll()", func(t *testing.T) {
		// Arrange
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)

		service := &ServiceError{statusCodeError: http.StatusTooManyRequests}
		log := logrus.New()
		handler := NewHandler(service, log)

		// Act
		err := handler.SearchExpensesAll(c)

		// Assertions
		if assert.NoError(t, err) {
			assert.Equal(t, service.statusCodeError, rec.Code)
		}
	})
}
