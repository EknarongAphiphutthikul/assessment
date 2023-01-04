package expenses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	storage Storage
}

func NewHandler(s Storage) *Handler {
	return &Handler{storage: s}
}

func (h Handler) AddExpenses(c echo.Context) error {
	req := new(ExpensesRequest)
	if err := c.Bind(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	id, err := h.storage.Insert(*req)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := ExpensesResponse{
		Id:     id,
		Title:  req.Title,
		Amount: req.Amount,
		Note:   req.Note,
		Tags:   req.Tags,
	}

	return c.JSON(http.StatusCreated, resp)
}
