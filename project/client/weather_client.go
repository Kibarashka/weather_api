package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
	"weather/project/config"
	"weather/project/domain"
)

const weatherAPIURL = "http://api.weatherapi.com/v1/current.json"

type WeatherAPIClient struct {
	apiKey     string
	httpClient *http.Client
}

func NewWeatherAPIClient(cfg config.Config) *WeatherAPIClient {
	if cfg.WeatherAPIKey == "" {
		log.Println("Warning: WeatherAPIKey is not set in config. Weather functionality will be disabled.")
	}
	return &WeatherAPIClient{
		apiKey: cfg.WeatherAPIKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *WeatherAPIClient) GetCurrentWeather(city string) (*domain.WeatherResponse, error) {
	if c.apiKey == "" {
		log.Println("WeatherAPIClient: API key not configured.")
		return nil, fmt.Errorf("weather API key is not configured")
	}

	params := url.Values{}
	params.Add("key", c.apiKey)
	params.Add("q", city)

	fullURL := fmt.Sprintf("%s?%s", weatherAPIURL, params.Encode())
	log.Printf("Fetching weather from: %s", fullURL)

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("client.GetCurrentWeather: error creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client.GetCurrentWeather: error performing request to WeatherAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusBadRequest { // WeatherAPI can return 400 for bad city
		return nil, domain.ErrCityNotFound
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("WeatherAPI request for city %s failed with status %s", city, resp.Status)
		return nil, fmt.Errorf("client.GetCurrentWeather: WeatherAPI request failed with status %s", resp.Status)
	}

	var apiResp domain.ExternalWeatherAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("client.GetCurrentWeather: error decoding WeatherAPI response: %w", err)
	}

	weather := &domain.WeatherResponse{
		Temperature: apiResp.Current.TempC,
		Humidity:    float64(apiResp.Current.Humidity),
		Description: apiResp.Current.Condition.Text,
	}

	return weather, nil
}
