package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/Brrocat/user-profile-service/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(databaseURL string) (*ProfileRepository, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ProfileRepository{db: pool}, nil
}

func (r *ProfileRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *ProfileRepository) CreateProfile(ctx context.Context, profile *models.CreateProfileRequest) (*models.UserProfile, error) {
	query := `
		INSERT INTO user_profiles (user_id, first_name, last_name, phone, date_of_birth)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	var id string
	var createdAt, updatedAt time.Time

	err := r.db.QueryRow(ctx, query,
		profile.UserID,
		profile.FirstName,
		profile.LastName,
		profile.Phone,
		profile.DateOfBirth,
	).Scan(&id, &createdAt, &updatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create user profile: %w", err)
	}

	return &models.UserProfile{
		ID:          id,
		UserID:      profile.UserID,
		FirstName:   profile.FirstName,
		LastName:    profile.LastName,
		Phone:       profile.Phone,
		DateOfBirth: profile.DateOfBirth,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func (r *ProfileRepository) GetProfileByID(ctx context.Context, id string) (*models.UserProfile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, date_of_birth, 
		       avatar_url, address, city, country, postal_code, driving_license,
		       created_at, updated_at
		FROM user_profiles
		WHERE id = $1
	`

	var profile models.UserProfile
	err := r.db.QueryRow(ctx, query, id).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.DateOfBirth,
		&profile.AvatarURL,
		&profile.Address,
		&profile.City,
		&profile.Country,
		&profile.PostalCode,
		&profile.DrivingLicense,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get profile by ID: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID string) (*models.UserProfile, error) {
	query := `
		SELECT id, user_id, first_name, last_name, phone, date_of_birth, 
		       avatar_url, address, city, country, postal_code, driving_license,
		       created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile models.UserProfile
	err := r.db.QueryRow(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.DateOfBirth,
		&profile.AvatarURL,
		&profile.Address,
		&profile.City,
		&profile.Country,
		&profile.PostalCode,
		&profile.DrivingLicense,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get profile by user ID: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID string, updates *models.UpdateProfileRequest) (*models.UserProfile, error) {
	query := `
		UPDATE user_profiles 
		SET first_name = COALESCE($1, first_name),
		    last_name = COALESCE($2, last_name),
		    phone = COALESCE($3, phone),
		    date_of_birth = COALESCE($4, date_of_birth),
		    avatar_url = COALESCE($5, avatar_url),
		    address = COALESCE($6, address),
		    city = COALESCE($7, city),
		    country = COALESCE($8, country),
		    postal_code = COALESCE($9, postal_code),
		    driving_license = COALESCE($10, driving_license),
		    updated_at = NOW()
		WHERE user_id = $11
		RETURNING id, user_id, first_name, last_name, phone, date_of_birth, 
		          avatar_url, address, city, country, postal_code, driving_license,
		          created_at, updated_at
	`

	var profile models.UserProfile
	err := r.db.QueryRow(ctx, query,
		updates.FirstName,
		updates.LastName,
		updates.Phone,
		updates.DateOfBirth,
		updates.AvatarURL,
		updates.Address,
		updates.City,
		updates.Country,
		updates.PostalCode,
		updates.DrivingLicense,
		userID,
	).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.FirstName,
		&profile.LastName,
		&profile.Phone,
		&profile.DateOfBirth,
		&profile.AvatarURL,
		&profile.Address,
		&profile.City,
		&profile.Country,
		&profile.PostalCode,
		&profile.DrivingLicense,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return &profile, nil
}

func (r *ProfileRepository) DeleteProfile(ctx context.Context, userID string) error {
	query := "DELETE FROM user_profiles WHERE user_id = $1"
	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("profile not found for user ID: %s", userID)
	}

	return nil
}
