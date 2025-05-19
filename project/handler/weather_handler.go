package handler

import (
	"errors"
	"log"
	"net/http"
	"weather/project/domain"
	"weather/project/service"

	"github.com/gin-gonic/gin"
)

type WeatherHandler struct {
	weatherService service.WeatherService
}

func NewWeatherHandler(ws service.WeatherService) *WeatherHandler {
	return &WeatherHandler{weatherService: ws}
}

func (h *WeatherHandler) GetWeather(c *gin.Context) {
	city := c.Query("city")
	if city == "" {
		log.Println("GetWeather handler: city parameter is missing")
		c.JSON(http.StatusBadRequest, gin.H{"error": "City parameter is required"})
		return
	}

	weather, err := h.weatherService.GetWeatherForCity(city)
	if err != nil {
		log.Printf("GetWeather handler: error from weatherService for city %s: %v", city, err)
		if errors.Is(err, domain.ErrCityNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": domain.ErrCityNotFound.Error()})
			return
		}
		if errors.Is(err, domain.ErrFailedToFetchWeather) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve weather information at this time"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred"})
		return
	}

	c.JSON(http.StatusOK, weather)
}
