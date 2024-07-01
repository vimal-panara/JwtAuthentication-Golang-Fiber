package middleware

import (
	"JwtAuthentication/handlers"
	"JwtAuthentication/helpers"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func VerifyToken(c *fiber.Ctx) error {
	token := c.Cookies("token")
	if token == "" {
		return errors.New("invalid token")
	}

	if err := helpers.ValidateJwtToken(token); err != nil {
		refToken := c.Cookies("refreshToken")
		if refTokErr := helpers.ValidateJwtToken(refToken); refTokErr != nil {
			return errors.New("invalid refresh token")
		} else {
			if tokenDbErr := handlers.CheckTokenInDb("", refToken); tokenDbErr != nil {
				return tokenDbErr
			}
			helpers.UpdateUserTokens(refToken)
		}
	} else {
		if tokenDbErr := handlers.CheckTokenInDb(token, ""); tokenDbErr != nil {
			return tokenDbErr
		}
		// helpers.UpdateUserTokens(token)
	}

	return nil
}
