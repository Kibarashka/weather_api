package handler

import (
	"errors"
	"log"
	"net/http"
	"weather/project/domain"
	"weather/project/service"

	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(ss service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{subscriptionService: ss}
}

func (h *SubscriptionHandler) Subscribe(c *gin.Context) {
	var input domain.SubscriptionInput

	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Subscribe handler: failed to bind input: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	_, err := h.subscriptionService.Subscribe(input)
	if err != nil {
		log.Printf("Subscribe handler: error from subscriptionService for email %s: %v", input.Email, err)
		if errors.Is(err, domain.ErrEmailAlreadySubscribed) {
			c.JSON(http.StatusConflict, gin.H{"error": domain.ErrEmailAlreadySubscribed.Error()})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process subscription request"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription request successful. A confirmation email has been sent (simulated)."})
}

func (h *SubscriptionHandler) ConfirmSubscription(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		log.Println("ConfirmSubscription handler: token parameter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Confirmation token is required"})
		return
	}

	err := h.subscriptionService.ConfirmSubscription(token)
	if err != nil {
		log.Printf("ConfirmSubscription handler: error from subscriptionService for token %s: %v", token, err)
		if errors.Is(err, domain.ErrTokenInvalidOrExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTokenInvalidOrExpired.Error()}) // 404 as per Swagger for not found
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm subscription"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription confirmed successfully"})
}

func (h *SubscriptionHandler) Unsubscribe(c *gin.Context) {
	token := c.Param("token")
	if token == "" {
		log.Println("Unsubscribe handler: token parameter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsubscribe token is required"})
		return
	}

	err := h.subscriptionService.UnsubscribeByToken(token)
	if err != nil {
		log.Printf("Unsubscribe handler: error from subscriptionService for token %s: %v", token, err)
		if errors.Is(err, domain.ErrTokenInvalidOrExpired) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrTokenInvalidOrExpired.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unsubscribe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Unsubscribed successfully"})
}
