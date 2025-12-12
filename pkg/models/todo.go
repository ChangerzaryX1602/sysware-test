package models

import (
	"time"

	"gorm.io/gorm"
)

type Status string

const (
	StatusPending Status = "pending"
	StatusDone    Status = "done"
)

type Todo struct {
	ID          uint           `json:"id" gorm:"primaryKey;autoIncrement;index" swaggerignore:"true"`
	Title       string         `json:"title" gorm:"not null"`
	Description string         `json:"description"`
	Status      Status         `json:"status" gorm:"not null;default:'pending'"`
	UserID      uint           `json:"user_id" gorm:"not null;index"`
	CreatedAt   time.Time      `json:"created_at" gorm:"autoCreateTime" swaggerignore:"true"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"autoUpdateTime" swaggerignore:"true"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index" swaggerignore:"true"`

	User User `json:"-" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" swaggerignore:"true"`
}
