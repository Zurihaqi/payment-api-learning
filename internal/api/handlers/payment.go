package handlers

import (
	"fmt"
	"net/http"

	"payment-api-learning/internal/models"
	"payment-api-learning/internal/storage"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	storage storage.Storage
}

func NewPaymentHandler(storage storage.Storage) *PaymentHandler {
	return &PaymentHandler{storage: storage}
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var paymentData models.Payment
	if err := c.ShouldBindJSON(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	paymentData.From = userID.(string)

	if err := h.validatePayment(&paymentData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.processPaymentTransaction(&paymentData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	h.storage.LogActivity(paymentData.From, "payment", fmt.Sprintf("Payment of %.2f to %s", paymentData.Amount, paymentData.To))
	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}

func (h *PaymentHandler) validatePayment(payment *models.Payment) error {
	_, err := h.storage.GetCustomerByID(payment.From)
	if err != nil {
		return fmt.Errorf("invalid sender")
	}

	_, err = h.storage.GetCustomerByID(payment.To)
	if err != nil {
		return fmt.Errorf("invalid recipient")
	}

	if payment.Amount <= 0 {
		return fmt.Errorf("invalid amount")
	}

	return nil
}

func (h *PaymentHandler) processPaymentTransaction(_ *models.Payment) error {
	//Simulasi payment berhasil
	return nil
}
