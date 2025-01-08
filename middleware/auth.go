// middleware/auth_middleware.go

package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() fiber.Handler {
	jwtSecret := os.Getenv("JWT_SECRET")

	return func(c *fiber.Ctx) error {
		// Get the Authorization header value
		authHeader := c.Get("Authorization")

		// Check if the Authorization header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return fiber.ErrUnauthorized
		}

		// Extract the token from the Authorization header
		tokenString := authHeader[7:] // Remove "Bearer " prefix

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Replace "your-secret-key" with your actual secret key used for signing the token
			return []byte(jwtSecret), nil
		})
		if err != nil {
			return fiber.ErrUnauthorized
		}

		// Check if the token is valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// Set the user information in the context for use in subsequent handlers
			c.Locals("user", claims)
			return c.Next()
		}

		return fiber.ErrUnauthorized
	}
}

func HasRole(claims jwt.MapClaims) int {
	roleID, exists := claims["rid"].(int)
	if !exists {
		return 0
	}

	return roleID
}
