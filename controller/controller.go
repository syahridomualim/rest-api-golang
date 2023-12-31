package controller

import (
	"database/sql"
	model "employee-golang/model"
	"employee-golang/service"
	"employee-golang/util"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	Service service.IEmployeeService
}

func EmployeeController(e *echo.Echo) {
	handler := &Controller{
		Service: service.NewEmployeeService(),
	}

	apis := e.Group("/api/v1/employees")
	apis.Add("GET", "", handler.GetEmployee)
	apis.Add("GET", "/:id", handler.GetEmployeeById)
	apis.Add("POST", "", handler.InsertEmployee)
	apis.Add("PUT", "", handler.UpdateEmployee)
	apis.Add("DELETE", "/:id", handler.DeleteEmployee)

	e.Logger.Fatal(e.Start(":8080"))
}

func (controller *Controller) InsertEmployee(c echo.Context) error {
	rq := new(model.Employee)

	c.Echo().Validator = &CustomValidator{
		Validator: validator.New(),
	}
	errBind := c.Bind(&rq)
	if errBind != nil {
		logrus.Printf("Body request is required %v", errBind)
		return createErrorResponse(c, 400, "BAD_REQUEST", "Body required", "Body request is required", errBind)
	}
	errValidate := c.Validate(rq)
	if errValidate != nil {
		logrus.Printf("Error: %s", errValidate)
		return createErrorResponse(c, 400, "BAD_REQUEST", errValidate.Error(), "Error: "+errValidate.Error(), errValidate)
	}

	body, err := controller.Service.InsertEmployee(rq)
	switch {
	case err != nil && err.Error() == "employee already exists":
		logrus.Printf("Error: %v", err)
		return createErrorResponse(c, 409, "CONFLICTED", err.Error(), "Error: "+err.Error(), err)
	case err != nil:
		logrus.Printf("Error getting employees %s", err)
		return createErrorResponse(c, 500, "INTERNAL_ERROR", err.Error(), "Error getting employees", err)
	default:
		logrus.Print("Success getting employees")
	}

	return createSuccessResponse(c, 200, body)
}

func (controller *Controller) UpdateEmployee(c echo.Context) error {
	rq := new(model.Employee)

	c.Echo().Validator = &CustomValidator{
		validator.New(),
	}
	errBind := c.Bind(rq)
	if errBind != nil {
		logrus.Printf("Body request is required %v", errBind)
		return createErrorResponse(c, 400, "BAD_REQUEST", "Body required", "Body request is required", errBind)
	}
	errValidate := c.Validate(rq)
	if errValidate != nil {
		logrus.Printf("Error: %s", errValidate)
		return createErrorResponse(c, 400, "BAD_REQUEST", errValidate.Error(), "Error: "+errValidate.Error(), errValidate)
	}

	body, err := controller.Service.UpdateEmployee(rq)
	if err != nil {
		logrus.Printf("Error updating employees %s", err)
		return createErrorResponse(c, 500, "INTERNAL_ERROR", err.Error(), "Error getting employees", err)
	}
	return createSuccessResponse(c, 200, body)
}
func (controller *Controller) GetEmployee(c echo.Context) error {
	response, err := controller.Service.GetEmployees()
	if err != nil {
		logrus.Printf("Error getting employees %v", err)
		return createErrorResponse(c, 500, "INTERNAL_ERROR", err.Error(), "Error getting employees", err)
	}

	return createSuccessResponse(c, 200, response)
}

func (controller *Controller) GetEmployeeById(c echo.Context) error {
	id := c.Param(`id`)
	response, err := controller.Service.GetEmployeeById(id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return createErrorResponse(c, 404, "NOT_FOUND", "Data Not Found", "Data not found", err)
	case false:
		return createErrorResponse(c, 500, "INTERNAL_ERROR", err.Error(), "Error getting employees", err)
	}

	return createSuccessResponse(c, 200, response)
}

func (controller *Controller) DeleteEmployee(c echo.Context) error {
	id := c.Param(`id`)

	response, err := controller.Service.DeleteEmployee(id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		return createErrorResponse(c, 404, "NOT_FOUND", "Data Not Found", "Data not found", err)
	case !errors.Is(err, sql.ErrNoRows):
		return createErrorResponse(c, 500, "INTERNAL_ERROR", err.Error(), "Error getting employees", err)
	}

	return createSuccessResponse(c, 200, response)
}

func createErrorResponse(c echo.Context, code int, status, message string, logMessage string, err error) error {
	response := model.GenericResponse[any]{
		Code:   code,
		Status: status,
		Data:   message,
	}
	logrus.Printf("%s %v", logMessage, err)
	return util.RespJSONData(c, code, response, err)
}

func createSuccessResponse(c echo.Context, code int, data interface{}) error {
	response := model.GenericResponse[any]{
		Code:   code,
		Status: "Success",
		Data:   data,
	}
	return util.RespJSONData(c, code, response, nil)
}
