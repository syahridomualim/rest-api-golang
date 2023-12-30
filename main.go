package main

import (
	"employee-golang/config"
	"employee-golang/controller"
	"github.com/labstack/echo/v4"
)

func init() {
	config.InitConfig(true)
}

func main() {
	e := *echo.New()
	config.InitSwagger(&e)
	controller.EmployeeController(&e)
}
