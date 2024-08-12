package app

import (
	"github.com/b0nbon1/VidFlux/database"
	_ "github.com/b0nbon1/VidFlux/docs"
	router "github.com/b0nbon1/VidFlux/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
)

//	@title			Swerv API
//	@version		1.0
//	@description	This is a sample server for Swerv API.
//	@termsOfService	http://swagger.io/terms/
//	@host			localhost:4500
//	@BasePath		/
func SetupAndRunApp() error {
	err := initdb.InitDb(GetEnvs("DB_URL"))
	if err != nil {
		return err
	}

	app := fiber.New(fiber.Config{
			StreamRequestBody: true,
		})

	app.Use(cors.New(cors.Config{
			AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
			AllowOrigins:     "*",
			AllowCredentials: false,
			AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		}))

	app.Static("/static", "./static")

	app.Get("/swagger/*", swagger.HandlerDefault)


	app.Use(recover.New())

	router.SetupRoutes(app)


	app.Listen(":" + GetEnvs("PORT"))

	return nil
}
