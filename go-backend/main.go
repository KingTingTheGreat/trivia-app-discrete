package main

import (
	"go-backend/controllers"
	"go-backend/shared"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Password string

// takes a function that returns an error and wraps it in a function that returns a JSON response
func restWrap(f func(echo.Context) error) func(echo.Context) error {
	return func(c echo.Context) error {
		err := f(c)
		if err != nil {
			return c.JSON(500, "Server error")
		}
		return nil
	}
}

func main() {
	shared.LoadPassword()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*:3000"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	go controllers.ResetBuzzers()
	go controllers.BroadcastLeaderboard()
	go controllers.BroadcastBuzzedIn()
	go controllers.BroadcastPlayerList()

	// player onboarding
	e.POST("/token", restWrap(controllers.PostToken))
	e.POST("/verify", restWrap(controllers.PostVerify))

	// player actions
	e.GET("/buzz", controllers.BuzzWs)

	// host controls
	e.POST("/reset", restWrap(controllers.Reset))
	e.PUT("/player", restWrap(controllers.UpdatePlayer))
	e.DELETE("/player", restWrap(controllers.DeletePlayer))

	// data retrieval
	e.GET("/players", controllers.GetPlayers)
	e.GET("/leaderboard", controllers.Leaderboard)
	e.GET("/buzzed-in", controllers.BuzzedIn)

	e.Logger.Fatal(e.Start(":8080"))
}
