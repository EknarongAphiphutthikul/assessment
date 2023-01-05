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
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func TestAddExpensesHandler(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		db, err := sql.Open("postgres", "postgresql://root:root@db/go-example-db?sslmode=disable")
		if err != nil {
			assert.Error(t, err)
		}
		logRus := logrus.New()
		storage := New(db)
		service := NewService(storage, logRus)
		handler := NewHandler(service, logRus)

		e.POST("/expenses", handler.AddExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
	}(eh)
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
