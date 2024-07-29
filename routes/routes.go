package route

import (
	"github.com/b0nbon1/VidFlux/pkg/videos"
	"github.com/gofiber/fiber/v2"
)

//	@BasePath	/api/v1
 func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to VideoFlux 👋!")
	})
 	v1 := app.Group("/api/v1")
	videos.VideosRoutes(v1)
 }
