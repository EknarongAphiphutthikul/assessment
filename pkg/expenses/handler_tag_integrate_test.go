//go:build integration

package expenses

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func setup(t *testing.T) (*sql.DB, func()) {
	// Setup server
	eh := echo.New()
	db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
	if err != nil {
		assert.Error(t, err)
	}
	go func(e *echo.Echo, db *sql.DB) {
		logRus := logrus.New()
		storage := New(db)
		service := NewService(storage, logRus)
		handler := NewHandler(service, logRus)

		e.POST("/expenses", handler.AddExpenses)
		e.GET("/expenses/:id", handler.SearchExpensesById)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh, db)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}
	return db, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = eh.Shutdown(ctx)
		assert.NoError(t, err)
		err = db.Close()
		assert.NoError(t, err)
	}
}

func TestAddExpensesHandlerIntegratetion(t *testing.T) {
	_, teardown := setup(t)
	defer teardown()
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(string(body)))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	respBody := &ExpensesResponse{}
	err = json.Unmarshal(byteBody, &respBody)
	assert.NoError(t, err)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, int64(0), respBody.Id)
		assert.Equal(t, reqBody.Title, respBody.Title)
		assert.Equal(t, reqBody.Amount, respBody.Amount)
		assert.Equal(t, reqBody.Note, respBody.Note)
		assert.Equal(t, reqBody.Tags, respBody.Tags)
	}
}

func TestSearchExpensesByIdIntegratetion(t *testing.T) {
	db, teardown := setup(t)
	defer teardown()
	// Arrange
	stmt, err := db.Prepare("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id")
	assert.NoError(t, err)
	defer stmt.Close()

	mockData := ExpensesRequest{
		Title:  "mockTitle",
		Amount: 10,
		Note:   "mockNote",
		Tags:   []string{"mockTags"},
	}
	row := stmt.QueryRow(mockData.Title, mockData.Amount, mockData.Note, pq.Array(mockData.Tags))

	var id int64
	err = row.Scan(&id)
	assert.NoError(t, err)

	targetUrl, err := url.Parse(fmt.Sprintf("http://localhost:%d/expenses", serverPort))
	assert.NoError(t, err)
	targetUrl = targetUrl.JoinPath(strconv.FormatInt(id, 10))

	req, err := http.NewRequest(http.MethodGet, targetUrl.String(), nil)
	fmt.Println(req.URL.String())
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	respBody := &ExpensesResponse{}
	err = json.Unmarshal(byteBody, &respBody)
	assert.NoError(t, err)

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, id, respBody.Id)
		assert.Equal(t, mockData.Title, respBody.Title)
		assert.Equal(t, mockData.Amount, respBody.Amount)
		assert.Equal(t, mockData.Note, respBody.Note)
		assert.Equal(t, mockData.Tags, respBody.Tags)
	}
}
