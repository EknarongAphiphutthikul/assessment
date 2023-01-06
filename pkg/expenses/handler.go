package expenses

import (
	"net/http"
	"strconv"

	"github.com/EknarongAphiphutthikul/assessment/pkg/common"
	"github.com/labstack/echo/v4"
)

type Services interface {
	AddExpenses(req ExpensesRequest) (*ExpensesResponse, error)
	SearchExpensesById(id int64) (*ExpensesResponse, error)
	UpdateExpenses(id int64, req ExpensesRequest) (*ExpensesResponse, error)
}
type Handler struct {
	log     common.Log
	service Services
}

func NewHandler(s Services, l common.Log) *Handler {
	return &Handler{service: s, log: l}
}

func (h Handler) AddExpenses(c echo.Context) error {
	req := ExpensesRequest{}
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	resp, err := h.service.AddExpenses(req)
	if err != nil {
		if cmErr, ok := err.(*common.Error); ok {
			return c.NoContent(cmErr.Code)
		}
		h.log.Errorf("Handler AddExpenses Error : %s", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h Handler) SearchExpensesById(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	resp, err := h.service.SearchExpensesById(id)
	if err != nil {
		if cmErr, ok := err.(*common.Error); ok {
			return c.NoContent(cmErr.Code)
		}
		h.log.Errorf("Handler SearchExpensesById Error : %s", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, resp)
}

func (h Handler) UpdateExpenses(c echo.Context) error {
	paramId := c.Param("id")
	id, err := strconv.ParseInt(paramId, 10, 64)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	req := ExpensesRequest{}
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	resp, err := h.service.UpdateExpenses(id, req)
	if err != nil {
		if cmErr, ok := err.(*common.Error); ok {
			return c.NoContent(cmErr.Code)
		}
		h.log.Errorf("Handler UpdateExpenses Error : %s", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, resp)
}
