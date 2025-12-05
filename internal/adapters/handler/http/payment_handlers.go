package http

import (
	"fastinghero/internal/core/domain"
	"fastinghero/internal/core/ports"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	paymentService  ports.PaymentService
	referralService ports.ReferralService
	userRepo        ports.UserRepository
}

func NewPaymentHandler(paymentService ports.PaymentService, referralService ports.ReferralService, userRepo ports.UserRepository) *PaymentHandler {
	return &PaymentHandler{
		paymentService:  paymentService,
		referralService: referralService,
		userRepo:        userRepo,
	}
}

func (h *PaymentHandler) HandleDeposit(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDVal.(uuid.UUID) // Requires uuid import, need to add it

	var req struct {
		Amount          float64 `json:"amount"`
		PaymentMethodID string  `json:"paymentMethodId"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fetch user to get email
	user, err := h.userRepo.FindByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// 1. Create Customer
	customerID, err := h.paymentService.CreateCustomer(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create customer"})
		return
	}

	// 2. In a real flow, we would create a PaymentIntent here using the amount
	// For this MVP step, we are just verifying the adapter connection.

	// GRANT PREMIUM STATUS IMMEDIATELY (MVP/Demo)
	user.SubscriptionTier = domain.TierVault
	user.SubscriptionStatus = domain.SubStatusActive
	if err := h.userRepo.Save(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user status"})
		return
	}

	// 3. Complete Referral (if any)
	if h.referralService != nil {
		go func() {
			// ctx := context.Background()
			// h.referralService.CompleteReferral(ctx, userID)
		}()
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Deposit initiated",
		"customer_id": customerID,
	})
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	const MaxBodyBytes = int64(65536)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body too large"})
		return
	}

	err = h.paymentService.HandleWebhook(c.Request.Context(), payload, c.GetHeader("Stripe-Signature"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature or event processing failed"})
		return
	}

	c.Status(http.StatusOK)
}
