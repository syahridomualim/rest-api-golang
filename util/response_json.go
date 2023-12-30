package util

import (
	"github.com/labstack/echo/v4"
	"os"
)

func RespJSONData(c echo.Context, sc int, i interface{}, err error) error {
	c.Set("completionStatus", sc)
	c.Set("instanceId", os.Hostname)
	if err != nil {
		c.Set("error", err.Error())
	}
	return c.JSON(sc, i)
}
