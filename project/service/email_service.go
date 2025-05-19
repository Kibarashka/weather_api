package service

import (
	"fmt"
	"log"
	"weather/project/config"
	"weather/project/domain"
)

type EmailService interface {
	SendConfirmationEmail(subscription *domain.Subscription, token string) error
	SendWeatherUpdateEmail(subscription *domain.Subscription, weather *domain.WeatherResponse) error
}

type emailService struct {
	cfg config.Config
}

func NewEmailService(cfg config.Config) EmailService {

	return &emailService{cfg: cfg}
}

func (s *emailService) SendConfirmationEmail(subscription *domain.Subscription, token string) error {
	if subscription == nil {
		return fmt.Errorf("subscription cannot be nil")
	}
	if token == "" {
		return fmt.Errorf("token cannot be empty")
	}

	confirmationLink := fmt.Sprintf("%s/api/confirm/%s", s.cfg.AppBaseURL, token)

	subject := "Confirm your Weather API Subscription"
	body := fmt.Sprintf("Hello %s,\n\nPlease confirm your subscription for weather updates in %s by clicking the link below:\n%s\n\nIf you did not request this, please ignore this email.\n\nThanks,\nThe Weather API Team",
		subscription.Email, subscription.City, confirmationLink)

	log.Printf("SIMULATING SENDING EMAIL:")
	log.Printf("To: %s", subscription.Email)
	log.Printf("From: %s", "noreply@weatherapp.dev") // s.cfg.EmailFrom if configured
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)

	log.Printf("Successfully simulated sending confirmation email to %s for city %s.", subscription.Email, subscription.City)
	return nil
}

func (s *emailService) SendWeatherUpdateEmail(subscription *domain.Subscription, weather *domain.WeatherResponse) error {
	if subscription == nil || weather == nil {
		return fmt.Errorf("subscription and weather data cannot be nil")
	}

	unsubscribeToken := "some_unsubscribe_token_placeholder" // This should come from subscription.UnsubscribeToken
	if subscription.UnsubscribeToken != nil {
		unsubscribeToken = *subscription.UnsubscribeToken
	}
	unsubscribeLink := fmt.Sprintf("%s/api/unsubscribe/%s", s.cfg.AppBaseURL, unsubscribeToken)

	subject := fmt.Sprintf("Weather Update for %s", subscription.City)
	body := fmt.Sprintf("Hello %s,\n\nHere's your weather update for %s:\nTemperature: %.1f°C\nHumidity: %.0f%%\nDescription: %s\n\nTo stop receiving these updates, click here: %s\n\nThanks,\nThe Weather API Team",
		subscription.Email, subscription.City, weather.Temperature, weather.Humidity, weather.Description, unsubscribeLink)

	log.Printf("SIMULATING SENDING WEATHER UPDATE EMAIL:")
	log.Printf("To: %s", subscription.Email)
	log.Printf("From: %s", "noreply@weatherapp.dev")
	log.Printf("Subject: %s", subject)
	log.Printf("Body:\n%s", body)

	log.Printf("Successfully simulated sending weather update to %s for %s.", subscription.Email, subscription.City)
	return nil
}
