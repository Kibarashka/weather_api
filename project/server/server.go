package server

import (
	"log"
	"net/http"
	"weather/project/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	weatherHandler *handler.WeatherHandler,
	subscriptionHandler *handler.SubscriptionHandler,
) *gin.Engine {

	router := gin.Default()

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "message": "Weather API is healthy"})
	})

	router.StaticFile("/swagger.yaml", "./swagger.yaml")

	apiGroup := router.Group("/api")
	{

		apiGroup.GET("/weather", weatherHandler.GetWeather)

		apiGroup.POST("/subscribe", subscriptionHandler.Subscribe)
		apiGroup.GET("/confirm/:token", subscriptionHandler.ConfirmSubscription)
		apiGroup.GET("/unsubscribe/:token", subscriptionHandler.Unsubscribe)
	}

	log.Println("Router setup complete.")
	return router
}
