package main

import (
	"JwtAuthentication/helpers"
	initializer "JwtAuthentication/initializers"
	"JwtAuthentication/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func init() {
	initializer.LoadEnvFile()
	initializer.ConnectToDatabase()
	helpers.LoadEncryptionKeys()
}

func main() {
	// fmt.Println("Starting Authentication service...")
	app := fiber.New()
	app.Use(recover.New())

	// cfg := swagger.Config{
	// 	BasePath: "/",
	// 	FilePath: "./docs/swagger.json",
	// 	Path:     "swagger",
	// 	Title:    "Swagger API Docs",
	// }

	// app.Use(swagger.New(cfg))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Fiber version 2 is working fine")
	})

	userApp := app.Group("/user")
	routes.AddUserRoutes(userApp)

	app.Listen(":" + port)
}
