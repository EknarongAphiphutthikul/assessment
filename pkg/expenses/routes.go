package expenses

import (
	"github.com/EknarongAphiphutthikul/assessment/pkg/config"
	"github.com/labstack/echo/v4"
)

func Routes(echo *echo.Echo, ins *config.Instance) {
	expenDb := New(ins.DB)
	expenService := NewService(expenDb, ins.Log)
	expenHandler := NewHandler(expenService, ins.Log)

	echo.POST("/expenses", expenHandler.AddExpenses)
}
