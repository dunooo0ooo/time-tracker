package models

import "time"

type User struct {
	ID             uint   `gorm:"primaryKey"`
	PassportNumber string `gorm:"unique;not null"`
	Surname        string
	Name           string
	Patronymic     string
	Address        string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Task struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint `gorm:"not null"`
	Description string
	StartTime   time.Time
	EndTime     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
