package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionFrequency string

const (
	FrequencyHourly SubscriptionFrequency = "hourly"
	FrequencyDaily  SubscriptionFrequency = "daily"
)

type Subscription struct {
	ID        uuid.UUID             `gorm:"type:char(36);primary_key;" json:"-"`
	Email     string                `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	City      string                `gorm:"type:varchar(100);not null" json:"city"`
	Frequency SubscriptionFrequency `gorm:"type:varchar(10);not null" json:"frequency"`
	Confirmed bool                  `gorm:"default:false" json:"confirmed"`

	ConfirmToken     *string        `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	UnsubscribeToken *string        `gorm:"type:varchar(64);uniqueIndex" json:"-"`
	CreatedAt        time.Time      `json:"-"`
	UpdatedAt        time.Time      `json:"-"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Subscription) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}

type SubscriptionInput struct {
	Email     string `form:"email" json:"email" binding:"required,email"`
	City      string `form:"city" json:"city" binding:"required,min=2"`
	Frequency string `form:"frequency" json:"frequency" binding:"required,oneof=hourly daily"`
}
