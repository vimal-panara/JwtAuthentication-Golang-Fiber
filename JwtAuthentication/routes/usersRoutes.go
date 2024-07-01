package routes

import (
	"JwtAuthentication/controllers"
	"JwtAuthentication/middleware"
	"JwtAuthentication/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AddUserRoutes(app fiber.Router) {

	app.Post("/login", func(c *fiber.Ctx) error {
		fmt.Println("/login api")
		return controllers.Login(c)
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		fmt.Println("/singup api")
		return controllers.Signup(c)
	})

	app.Get("/get_all", func(c *fiber.Ctx) error {
		fmt.Println("/get_all api")
		if err := middleware.VerifyToken(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				StatusCode: fiber.StatusUnauthorized,
				Msg:        err.Error(),
				Data:       nil,
			})
		}
		return controllers.GetAllUsers(c)
	})

	app.Get("/get/:id", func(c *fiber.Ctx) error {
		fmt.Println("/get/id api")
		if err := middleware.VerifyToken(c); err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
				StatusCode: fiber.StatusUnauthorized,
				Msg:        err.Error(),
				Data:       nil,
			})
		}
		return controllers.GetUserById(c)
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		fmt.Println("/logout api")
		return controllers.Logout(c)
	})

}
