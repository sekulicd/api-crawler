package rest

import (
	"github.com/labstack/echo"
	"net/http"
	"runtime"
	"strconv"
)

func (s *service) hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func (s *service) getAllSchools(c echo.Context) error {
	schools, err := s.collegeService.GetAllSchools()
	if err != nil {
		c.Error(err)
		return err
	}
	c.JSON(200, schools)
	return nil
}

func (s *service) num(c echo.Context) error {
	return c.String(http.StatusOK, strconv.Itoa(runtime.NumGoroutine()))
}
