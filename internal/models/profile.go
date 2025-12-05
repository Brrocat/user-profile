package models

import (
	"time"
)

type UserProfile struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	DateOfBirth    string    `json:"date_of_birth"` // YYYY-MM-DD
	AvatarURL      string    `json:"avatar_url"`
	Address        string    `json:"address"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
	PostalCode     string    `json:"postal_code"`
	DrivingLicense string    `json:"driving_license"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateProfileRequest struct {
	UserID      string `json:"user_id" validate:"required"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Phone       string `json:"phone"`
	DateOfBirth string `json:"date_of_birth"`
}

type UpdateProfileRequest struct {
	UserID         string `json:"user_id" validate:"required"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Phone          string `json:"phone"`
	DateOfBirth    string `json:"date_of_birth"`
	AvatarURL      string `json:"avatar_url"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`
	PostalCode     string `json:"postal_code"`
	DrivingLicense string `json:"driving_license"`
}
