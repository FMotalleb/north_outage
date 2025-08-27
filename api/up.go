package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func init() {
	api.GET("/up", up)
}

func up(c echo.Context) error {
	return c.String(http.StatusOK, "1")
}
