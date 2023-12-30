package config

import (
	_ "employee-golang/docs"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	swagger "github.com/swaggo/echo-swagger"
	"net/http"
)

func InitSwagger(e *echo.Echo) {

	e.GET("/swagger/*", swagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
}
