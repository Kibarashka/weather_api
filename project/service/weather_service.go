package service

import (
	"errors"
	"log"
	"weather/project/client"
	"weather/project/domain"
)

type WeatherService interface {
	GetWeatherForCity(city string) (*domain.WeatherResponse, error)
}

type weatherService struct {
	weatherAPIClient *client.WeatherAPIClient
}

func NewWeatherService(apiClient *client.WeatherAPIClient) WeatherService {
	return &weatherService{
		weatherAPIClient: apiClient,
	}
}

func (s *weatherService) GetWeatherForCity(city string) (*domain.WeatherResponse, error) {
	if city == "" {
		return nil, domain.ErrCityNotFound
	}
	if s.weatherAPIClient == nil {
		log.Println("WeatherService: weatherAPIClient is nil")
		return nil, errors.New("weather service is not properly initialized")
	}

	log.Printf("Fetching weather for city: %s", city)
	weather, err := s.weatherAPIClient.GetCurrentWeather(city)
	if err != nil {
		log.Printf("Error fetching weather for city %s from API client: %v", city, err)
		if errors.Is(err, domain.ErrCityNotFound) {
			return nil, domain.ErrCityNotFound
		}
		return nil, domain.ErrFailedToFetchWeather
	}

	log.Printf("Successfully fetched weather for %s: %+v", city, weather)
	return weather, nil
}
