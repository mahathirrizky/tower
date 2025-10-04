package controllers

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/user/tower-tracker-bima/backend/database"
	"github.com/user/tower-tracker-bima/backend/helper"
	"github.com/user/tower-tracker-bima/backend/models"
)

type RegisterInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := helper.HashPassword(input.Password)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := models.User{Name: input.Name, Email: input.Email, Password: hashedPassword}

	if err := database.DB.Create(&user).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to create user: "+err.Error())
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Registration successful", nil)
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	// Log raw request body for debugging
	body, _ := c.GetRawData()
	log.Printf("Login Request Raw Body: %s", string(body))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body)) // Restore body for subsequent reads

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Login Request Binding Error: %v", err)
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !helper.CheckPasswordHash(input.Password, user.Password) {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := helper.GenerateJWT(user.ID)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Login successful", gin.H{"token": token})
}

type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
}

func ChangePassword(c *gin.Context) {
	var input ChangePasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helper.SendErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Custom validation for new password
	if len(input.NewPassword) < 6 {
		helper.SendErrorResponse(c, http.StatusBadRequest, "New password must be at least 6 characters.")
		return
	}
	if ok, _ := regexp.MatchString(`[a-z]`, input.NewPassword); !ok {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Must have a lowercase letter.")
		return
	}
	if ok, _ := regexp.MatchString(`[A-Z]`, input.NewPassword); !ok {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Must have an uppercase letter.")
		return
	}
	if ok, _ := regexp.MatchString(`\d`, input.NewPassword); !ok {
		helper.SendErrorResponse(c, http.StatusBadRequest, "Must have a number.")
		return
	}

	userID := c.MustGet("user_id").(uint) // Get user_id from JWT claims
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusNotFound, "User not found")
		return
	}

	// Verify current password
	if !helper.CheckPasswordHash(input.CurrentPassword, user.Password) {
		helper.SendErrorResponse(c, http.StatusUnauthorized, "Invalid current password")
		return
	}

	// Hash new password
	hashedPassword, err := helper.HashPassword(input.NewPassword)
	if err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to hash new password")
		return
	}

	// Update password
	user.Password = hashedPassword
	if err := database.DB.Save(&user).Error; err != nil {
		helper.SendErrorResponse(c, http.StatusInternalServerError, "Failed to update password")
		return
	}

	helper.SendSuccessResponse(c, http.StatusOK, "Password updated successfully", nil)
}
