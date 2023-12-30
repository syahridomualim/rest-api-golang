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

	e.Logger.Fatal(e.Start(":8080"))
}

func (controller *Controller) InsertEmployee(c echo.Context) error {
	rq := new(model.Employee)
	response := &model.GenericResponse[any]{}

	c.Echo().Validator = &CustomValidator{
		Validator: validator.New(),
	}
	errBind := c.Bind(&rq)
	if errBind != nil {
		response.Status = "BAD_REQUEST"
		response.Code = 400
		response.Data = "Body required"
		logrus.Printf("Body request is required %v", errBind)
		return util.RespJSONData(c, response.Code, response, errBind)
	}
	errValidate := c.Validate(rq)
	if errValidate != nil {
		response.Status = "BAD_REQUEST"
		response.Code = 400
		response.Data = errValidate.Error()
		logrus.Printf("Error: %s", errValidate)
		return util.RespJSONData(c, response.Code, response, errValidate)
	}

	body, err := controller.Service.InsertEmployee(rq)
	switch {
	case err != nil && err.Error() == "employee already exists":
		response.Status = "CONFLICTED"
		response.Code = 409
		response.Data = err.Error()
		logrus.Printf("Error: %v", err)
		return util.RespJSONData(c, response.Code, response, errValidate)
	case err != nil:
		response.Status = "INTERNAL_ERROR"
		response.Code = 500
		response.Data = err.Error()
		logrus.Printf("Error getting employees %s", err)
		return util.RespJSONData(c, response.Code, response, err)
	default:
		logrus.Print("Success getting employees")
	}
	response.Status = "Success"
	response.Code = 200
	response.Data = body
	return util.RespJSONData(c, response.Code, response, err)
}

func (controller *Controller) UpdateEmployee(c echo.Context) error {
	rq := new(model.Employee)
	response := model.GenericResponse[string]{}
	c.Echo().Validator = &CustomValidator{
		validator.New(),
	}
	errBind := c.Bind(rq)
	if errBind != nil {
		response.Status = "BAD_REQUEST"
		response.Code = 400
		response.Data = "Body required"
		logrus.Printf("Body request is required %v", errBind)
		return util.RespJSONData(c, response.Code, response, errBind)
	}
	errValidate := c.Validate(rq)
	if errValidate != nil {
		response.Status = "BAD_REQUEST"
		response.Code = 400
		response.Data = errValidate.Error()
		logrus.Printf("Error: %s", errValidate)
		return util.RespJSONData(c, response.Code, response, errValidate)
	}

	body, err := controller.Service.UpdateEmployee(rq)
	if err != nil {
		response.Status = "INTERNAL_ERROR"
		response.Code = 500
		response.Data = err.Error()
		logrus.Printf("Error updating employees %s", err)
		return util.RespJSONData(c, response.Code, response, err)
	}
	response.Status = "Success"
	response.Code = 200
	response.Data = body
	return util.RespJSONData(c, response.Code, response, err)
}
func (controller *Controller) GetEmployee(c echo.Context) error {
	rs := new(model.GenericResponse[[]*model.Employee])
	response, err := controller.Service.GetEmployees()
	if err != nil {
		rs.Status = "INTERNAL_ERROR"
		rs.Code = 500
		rs.Data = nil
		logrus.Printf("Error getting employees %v", err)
		return util.RespJSONData(c, rs.Code, rs, err)
	}

	rs.Status = "Success"
	rs.Code = 200
	rs.Data = response
	return util.RespJSONData(c, rs.Code, rs, err)
}

func (controller *Controller) GetEmployeeById(c echo.Context) error {
	id := c.Param(`id`)
	rs := new(model.GenericResponse[any])
	response, err := controller.Service.GetEmployeeById(id)

	switch {
	case errors.Is(err, sql.ErrNoRows):
		rs.Status = "NOT_FOUND"
		rs.Code = 404
		rs.Data = "Data Not Found"
		logrus.Warningf("%v", err)
		return util.RespJSONData(c, rs.Code, rs, err)
	case false:
		rs.Status = "INTERNAL_ERROR"
		rs.Code = 500
		rs.Data = err
		logrus.Printf("Error getting employees %s", err)
		return util.RespJSONData(c, rs.Code, rs, err)
	}

	rs.Status = "Success"
	rs.Code = 200
	rs.Data = response
	return util.RespJSONData(c, rs.Code, rs, err)
}
