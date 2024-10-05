package handlers

import (
	"net/http"
	"strings"
	"time"

	"payment-api-learning/internal/models"
	"payment-api-learning/internal/storage"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type AuthHandler struct {
	storage        storage.Storage
	tokenBlacklist map[string]bool
}

func NewAuthHandler(storage storage.Storage) *AuthHandler {
	return &AuthHandler{
		storage:        storage,
		tokenBlacklist: make(map[string]bool),
	}
}

func (h *AuthHandler) IsTokenBlacklisted(token string) bool {
	return h.tokenBlacklist[token]
}

func (h *AuthHandler) Logout(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")

	userID := c.GetString("user_id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	h.tokenBlacklist[tokenString] = true

	h.storage.LogActivity(userID, "logout", "User logged out")
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var loginCustomer models.Customer
	if err := c.ShouldBindJSON(&loginCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	customer, err := h.storage.GetCustomerByUsername(loginCustomer.Username)
	if err != nil || customer.Password != loginCustomer.Password {
		h.storage.LogActivity("", "login", "Failed login attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": customer.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("zul-fahri-baihaqi-batch17enigma-jakarta"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	h.storage.LogActivity(customer.ID, "login", "Successful login")
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *AuthHandler) TokenBlacklist() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization token"})
			c.Abort()
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		if h.tokenBlacklist[tokenString] {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			c.Abort()
			return
		}

		c.Next()
	}
}
