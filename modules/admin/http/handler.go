package http

import (
	"go_whatsapp/helper"
	"go_whatsapp/models"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{
		db: db,
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}

// HandleLogin handles admin login
func (h *AdminHandler) HandleLogin(c *fiber.Ctx) error {
	// Parse request body
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(helper.APIResponse("Invalid request data", http.StatusBadRequest, "ERROR", nil))
	}

	// Find admin user by username
	var admin models.AdminUser
	if err := h.db.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		return c.Status(http.StatusUnauthorized).JSON(helper.APIResponse("Invalid credentials", http.StatusUnauthorized, "ERROR", nil))
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(helper.APIResponse("Invalid credentials", http.StatusUnauthorized, "ERROR", nil))
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": admin.ID,
		"username": admin.Username,
		"name": admin.Name,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to generate token", http.StatusInternalServerError, "ERROR", nil))
	}

	// Return the token
	return c.Status(http.StatusOK).JSON(helper.APIResponse("Login successful", http.StatusOK, "SUCCESS", fiber.Map{
		"token": tokenString,
		"user": fiber.Map{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"email":    admin.Email,
		},
	}))
}

// HandleRegister handles admin registration (this would typically be restricted)
func (h *AdminHandler) HandleRegister(c *fiber.Ctx) error {
	// Parse request body
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(helper.APIResponse("Invalid request data", http.StatusBadRequest, "ERROR", nil))
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to hash password", http.StatusInternalServerError, "ERROR", nil))
	}

	// Create new admin user
	admin := models.AdminUser{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to database
	if err := h.db.Create(&admin).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(helper.APIResponse("Failed to create admin user", http.StatusInternalServerError, "ERROR", nil))
	}

	return c.Status(http.StatusCreated).JSON(helper.APIResponse("Admin user created successfully", http.StatusCreated, "SUCCESS", fiber.Map{
		"id":       admin.ID,
		"username": admin.Username,
		"name":     admin.Name,
		"email":    admin.Email,
	}))
}

// HandleGetProfile gets the admin profile
func (h *AdminHandler) HandleGetProfile(c *fiber.Ctx) error {
	// Get user ID from JWT token
	user := c.Locals("user").(jwt.MapClaims)
	userID := uint(user["sub"].(float64))

	// Find admin user by ID
	var admin models.AdminUser
	if err := h.db.First(&admin, userID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(helper.APIResponse("Admin user not found", http.StatusNotFound, "ERROR", nil))
	}

	return c.Status(http.StatusOK).JSON(helper.APIResponse("Admin profile retrieved successfully", http.StatusOK, "SUCCESS", fiber.Map{
		"id":       admin.ID,
		"username": admin.Username,
		"name":     admin.Name,
		"email":    admin.Email,
	}))
}
