package service

import (
	"errors"
	"fmt"
	"log"
	"time"
	"weather/project/domain"
	"weather/project/repository"
)

type SubscriptionService interface {
	Subscribe(input domain.SubscriptionInput) (*domain.Subscription, error)
	ConfirmSubscription(token string) error
	UnsubscribeByToken(token string) error
}

type subscriptionService struct {
	repo         repository.SubscriptionRepository
	tokenService TokenService
	emailService EmailService
}

func NewSubscriptionService(
	repo repository.SubscriptionRepository,
	tokenService TokenService,
	emailService EmailService,
) SubscriptionService {
	return &subscriptionService{
		repo:         repo,
		tokenService: tokenService,
		emailService: emailService,
	}
}

func (s *subscriptionService) Subscribe(input domain.SubscriptionInput) (*domain.Subscription, error) {
	existingSub, err := s.repo.FindByEmail(input.Email)

	if err != nil && !errors.Is(err, domain.ErrSubscriptionNotFound) {
		log.Printf("Error finding subscription by email %s: %v", input.Email, err)
		return nil, fmt.Errorf("failed to check for existing subscription: %w", err)
	}

	if existingSub != nil {
		if existingSub.Confirmed {
			log.Printf("Attempt to subscribe with already confirmed email: %s", input.Email)
			return nil, domain.ErrEmailAlreadySubscribed
		}

		log.Printf("Email %s exists but not confirmed. Updating and re-sending confirmation.", input.Email)

		confirmToken, tokenErr := s.tokenService.GenerateToken(32)
		if tokenErr != nil {
			log.Printf("Error generating new confirmation token for %s: %v", input.Email, tokenErr)
			return nil, fmt.Errorf("failed to generate confirmation token: %w", tokenErr)
		}

		existingSub.City = input.City
		existingSub.Frequency = domain.SubscriptionFrequency(input.Frequency)
		existingSub.ConfirmToken = &confirmToken
		existingSub.UpdatedAt = time.Now()

		if updateErr := s.repo.Update(existingSub); updateErr != nil {
			log.Printf("Error updating existing unconfirmed subscription for %s: %v", input.Email, updateErr)
			return nil, fmt.Errorf("failed to update subscription: %w", updateErr)
		}

		go s.sendConfirmationEmailAsync(existingSub, confirmToken)
		return existingSub, nil
	}

	confirmToken, err := s.tokenService.GenerateToken(32)
	if err != nil {
		log.Printf("Error generating confirmation token for new sub %s: %v", input.Email, err)
		return nil, fmt.Errorf("failed to generate confirmation token: %w", err)
	}

	newSub := &domain.Subscription{

		Email:        input.Email,
		City:         input.City,
		Frequency:    domain.SubscriptionFrequency(input.Frequency),
		Confirmed:    false,
		ConfirmToken: &confirmToken,
	}

	if err := s.repo.Create(newSub); err != nil {

		log.Printf("Error creating new subscription for %s: %v", input.Email, err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	go s.sendConfirmationEmailAsync(newSub, confirmToken)

	log.Printf("New subscription initiated for %s, city %s. Confirmation pending.", newSub.Email, newSub.City)
	return newSub, nil
}

func (s *subscriptionService) sendConfirmationEmailAsync(sub *domain.Subscription, token string) {
	if s.emailService != nil {
		err := s.emailService.SendConfirmationEmail(sub, token)
		if err != nil {
			log.Printf("Async sendConfirmationEmail: Failed to send email to %s: %v", sub.Email, err)
		}
	} else {
		log.Printf("Async sendConfirmationEmail: EmailService is nil. Email to %s not sent.", sub.Email)
	}
}

func (s *subscriptionService) ConfirmSubscription(token string) error {
	if token == "" {
		return domain.ErrTokenInvalidOrExpired
	}
	sub, err := s.repo.FindByConfirmToken(token)
	if err != nil {

		log.Printf("Error finding subscription by confirm token %s: %v", token, err)
		return err
	}

	if sub.Confirmed {
		log.Printf("Subscription for email %s already confirmed.", sub.Email)
		return nil
	}

	sub.Confirmed = true
	sub.ConfirmToken = nil
	sub.UpdatedAt = time.Now()

	unsubscribeToken, tokenErr := s.tokenService.GenerateToken(32)
	if tokenErr != nil {

		log.Printf("Error generating unsubscribe token for %s after confirmation: %v", sub.Email, tokenErr)
	} else {
		sub.UnsubscribeToken = &unsubscribeToken
	}

	if err := s.repo.Update(sub); err != nil {
		log.Printf("Error updating subscription to confirmed for %s: %v", sub.Email, err)
		return fmt.Errorf("failed to confirm subscription in DB: %w", err)
	}

	log.Printf("Subscription for %s confirmed successfully.", sub.Email)
	return nil
}

func (s *subscriptionService) UnsubscribeByToken(token string) error {
	if token == "" {
		return domain.ErrTokenInvalidOrExpired
	}
	sub, err := s.repo.FindByUnsubscribeToken(token)
	if err != nil {
		log.Printf("Error finding subscription by unsubscribe token %s: %v", token, err)
		return err
	}

	if err := s.repo.Delete(sub.ID); err != nil {
		log.Printf("Error deleting (unsubscribing) subscription ID %s for email %s: %v", sub.ID, sub.Email, err)
		return fmt.Errorf("failed to unsubscribe: %w", err)
	}

	log.Printf("Email %s (ID: %s) unsubscribed successfully using token.", sub.Email, sub.ID)

	return nil
}
