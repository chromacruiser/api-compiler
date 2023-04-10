package main

import (
	"github.com/chromacruiser/api-compiler/internal/api"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api.RegisterHandlers(e, api.Handlers{})

	e.Logger.Fatal(e.Start(":8080"))
}
