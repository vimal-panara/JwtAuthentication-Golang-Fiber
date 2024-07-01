package controllers

import (
	"JwtAuthentication/handlers"
	"JwtAuthentication/helpers"
	"JwtAuthentication/models"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *fiber.Ctx) error {

	var userLoginData models.LoginRequest
	if err := c.BodyParser(&userLoginData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "invalid body recieved",
			Data:       nil,
		})
	}

	if userLoginData.UserName == "" || userLoginData.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "Either username or password is not recieveds",
			Data:       nil,
		})
	}

	user, err := handlers.FindOneUser(&models.User{Email: userLoginData.UserName, Mobile: userLoginData.UserName})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "Username is wrong",
			Data:       nil,
		})
	}

	encPass, err := helpers.GetEncryptedPassword(userLoginData.Password)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "password is wrong",
			Data:       nil,
		})
	} else {
		// fmt.Println(encPass)
		if encPass != user.Password {
			return c.Status(fiber.StatusBadRequest).JSON(models.Response{
				StatusCode: fiber.StatusBadRequest,
				Msg:        "password is wrong",
				Data:       nil,
			})
		}
	}

	//Generating new token for the user coming for login
	token, refToken, err := helpers.GenerateJwtToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        "unable to generate token",
			Data:       nil,
		})
	}

	user.Token = token
	user.RefreshToken = refToken
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := handlers.UpdateUserTokens(user, user.Id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        err.Error(),
			Data:       nil,
		})
	}

	c.Cookie(&fiber.Cookie{Name: "token", Value: token})
	c.Cookie(&fiber.Cookie{Name: "refreshToken", Value: refToken})
	return c.Status(fiber.StatusOK).JSON(models.Response{
		StatusCode: fiber.StatusOK,
		Msg:        "login successfull",
		Data:       nil,
	})
}

func Signup(c *fiber.Ctx) error {

	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "can not parse the request body",
			Data:       nil,
		})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        err.Error(),
			Data:       nil,
		})
	}

	token, refToken, err := helpers.GenerateJwtToken(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        err.Error(),
			Data:       nil,
		})
	}

	user.Token = token
	user.RefreshToken = refToken
	user.IsActive = true
	user.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	encPass, err := helpers.GetEncryptedPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        err.Error(),
			Data:       nil,
		})
	}

	user.Password = encPass
	result := handlers.AddOneUser(user)
	return c.Status(fiber.StatusOK).JSON(result)
}

func GetAllUsers(c *fiber.Ctx) error {
	users, err := handlers.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        err.Error(),
			Data:       nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(models.Response{
		StatusCode: fiber.StatusOK,
		Msg:        "users data found",
		Data:       users,
	})
}

func GetUserById(c *fiber.Ctx) error {

	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "Invalid id recieved",
			Data:       nil,
		})
	}
	// fmt.Println("Id: ", id)

	user, err := handlers.FindUserById(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
			StatusCode: fiber.StatusBadRequest,
			Msg:        "No user found for the given id",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.Response{
		StatusCode: fiber.StatusOK,
		Msg:        "user found",
		Data:       user,
	})
}

func Logout(c *fiber.Ctx) error {
	token := c.Cookies("token")
	refToken := c.Cookies("refreshToken")

	if token == "" || refToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
			StatusCode: fiber.StatusUnauthorized,
			Msg:        "token not found",
			Data:       nil,
		})
	}
	// fmt.Println(token, refToken)

	user, err := helpers.GetEmailMobileFromToken(token)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        "invalid token",
			Data:       nil,
		})
	}

	if err := handlers.UpdateUserTokens(user, primitive.NewObjectID()); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			StatusCode: fiber.StatusInternalServerError,
			Msg:        err.Error(),
			Data:       nil,
		})
	}

	c.Cookie(&fiber.Cookie{Name: "token", Value: ""})
	c.Cookie(&fiber.Cookie{Name: "refreshToken", Value: ""})

	return c.Status(fiber.StatusOK).JSON(models.Response{
		StatusCode: fiber.StatusOK,
		Msg:        "Logged out successfully",
		Data:       nil,
	})
}
