package repository

import (
	"errors"
	"weather/project/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(sub *domain.Subscription) error
	FindByEmail(email string) (*domain.Subscription, error)
	FindByConfirmToken(token string) (*domain.Subscription, error)
	FindByUnsubscribeToken(token string) (*domain.Subscription, error)
	Update(sub *domain.Subscription) error
	Delete(id uuid.UUID) error
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(sub *domain.Subscription) error {
	return r.db.Create(sub).Error
}

func (r *subscriptionRepository) FindByEmail(email string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.Where("email = ?", email).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSubscriptionNotFound
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) FindByConfirmToken(token string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.Where("confirm_token = ?", token).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTokenInvalidOrExpired
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) FindByUnsubscribeToken(token string) (*domain.Subscription, error) {
	var sub domain.Subscription
	err := r.db.Where("unsubscribe_token = ?", token).First(&sub).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrTokenInvalidOrExpired
		}
		return nil, err
	}
	return &sub, nil
}

func (r *subscriptionRepository) Update(sub *domain.Subscription) error {

	if sub.ID == uuid.Nil {
		return errors.New("cannot update subscription without ID")
	}
	return r.db.Save(sub).Error
}

func (r *subscriptionRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&domain.Subscription{}, "id = ?", id).Error
}
