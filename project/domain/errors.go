package domain

import "errors"

var (
	ErrCityNotFound           = errors.New("city not found by external weather API")
	ErrEmailAlreadySubscribed = errors.New("email already subscribed and confirmed")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrTokenInvalidOrExpired  = errors.New("token is invalid, expired, or not found")
	ErrFailedToFetchWeather   = errors.New("failed to fetch weather data from external API")
	ErrEmailSendingFailed     = errors.New("failed to send email")
)
