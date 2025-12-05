package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/Brrocat/user-profile-service/internal/models"
	"github.com/Brrocat/user-profile-service/internal/repository/postgres"
	"github.com/Brrocat/user-profile-service/internal/repository/redis"
	"github.com/Brrocat/user-profile-service/pkg/validation"
	"log/slog"
)

var (
	ErrProfileNotFound      = errors.New("profile not found")
	ErrProfileAlreadyExists = errors.New("profile already exists")
	ErrInvalidData          = errors.New("invalid data")
)

type ProfileService struct {
	profileRepo *postgres.ProfileRepository
	cacheRepo   *redis.CacheRepository
	validator   *validation.Validator
	logger      *slog.Logger
}

func NewProfileService(
	profileRepo *postgres.ProfileRepository,
	cacheRepo *redis.CacheRepository,
	validator *validation.Validator,
	logger *slog.Logger,
) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		cacheRepo:   cacheRepo,
		validator:   validator,
		logger:      logger,
	}
}

func (s *ProfileService) GetUserProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	s.logger.Debug("Getting user profile", "user_id", userID)

	// Try to get from cache first
	cachedProfile, err := s.cacheRepo.GetCachedProfile(ctx, userID)
	if err != nil {
		s.logger.Warn("Failed to get profile from cache", "user_id", userID, "error", err)
		// Continue to database lookup
	}

	if cachedProfile != nil {
		s.logger.Debug("Profile found in cache", "user_id", userID)
		return cachedProfile, nil
	}

	// Get from database
	profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get profile from database", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	if profile == nil {
		s.logger.Debug("Profile not found", "user_id", userID)
		return nil, ErrProfileNotFound
	}

	// Cache the profile for future requests
	if err := s.cacheRepo.CacheProfile(ctx, profile); err != nil {
		s.logger.Warn("Failed to cache profile", "user_id", userID, "error", err)
		// Non-critical error, continue
	}

	s.logger.Debug("Profile retrieved from database", "user_id", userID)
	return profile, nil
}

func (s *ProfileService) CreateUserProfile(ctx context.Context, req *models.CreateProfileRequest) (*models.UserProfile, error) {
	s.logger.Debug("Creating user profile", "user_id", req.UserID)

	// Validate input
	if err := s.validator.ValidateStruct(req); err != nil {
		validationErrors := s.validator.FormatValidationErrors(err)
		s.logger.Warn("Validation failed for create profile", "user_id", req.UserID, "errors", validationErrors)
		return nil, fmt.Errorf("%w: %v", ErrInvalidData, validationErrors)
	}

	// Check if profile already exists
	existingProfile, err := s.profileRepo.GetProfileByUserID(ctx, req.UserID)
	if err != nil {
		s.logger.Error("Failed to check existing profile", "user_id", req.UserID, "error", err)
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}

	if existingProfile != nil {
		s.logger.Warn("Profile already exists", "user_id", req.UserID)
		return nil, ErrProfileAlreadyExists
	}

	// Create profile
	profile, err := s.profileRepo.CreateProfile(ctx, req)
	if err != nil {
		s.logger.Error("Failed to create profile", "user_id", req.UserID, "error", err)
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	// Cache the new profile
	if err := s.cacheRepo.CacheProfile(ctx, profile); err != nil {
		s.logger.Warn("Failed to cache new profile", "user_id", req.UserID, "error", err)
		// Non-critical error, continue
	}

	s.logger.Info("Profile created successfully", "user_id", req.UserID, "profile_id", profile.ID)
	return profile, nil
}

func (s *ProfileService) UpdateUserProfile(ctx context.Context, userID string, req *models.UpdateProfileRequest) (*models.UserProfile, error) {
	s.logger.Debug("Updating user profile", "user_id", userID)

	// Validate input
	if err := s.validator.ValidateStruct(req); err != nil {
		validationErrors := s.validator.FormatValidationErrors(err)
		s.logger.Warn("Validation failed for update profile", "user_id", userID, "errors", validationErrors)
		return nil, fmt.Errorf("%w: %v", ErrInvalidData, validationErrors)
	}

	// Check if profile exists
	existingProfile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get existing profile for update", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to get existing profile: %w", err)
	}

	if existingProfile == nil {
		s.logger.Warn("Profile not found for update", "user_id", userID)
		return nil, ErrProfileNotFound
	}

	// Update profile
	updatedProfile, err := s.profileRepo.UpdateProfile(ctx, userID, req)
	if err != nil {
		s.logger.Error("Failed to update profile", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	if updatedProfile == nil {
		s.logger.Error("Profile update returned nil", "user_id", userID)
		return nil, fmt.Errorf("profile update failed")
	}

	// Update cache
	if err := s.cacheRepo.CacheProfile(ctx, updatedProfile); err != nil {
		s.logger.Warn("Failed to update cached profile", "user_id", userID, "error", err)
		// Non-critical error, continue
	}

	s.logger.Info("Profile updated successfully", "user_id", userID, "profile_id", updatedProfile.ID)
	return updatedProfile, nil
}

func (s *ProfileService) DeleteUserProfile(ctx context.Context, userID string) error {
	s.logger.Debug("Deleting user profile", "user_id", userID)

	// Check if profile exists
	existingProfile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to get existing profile for deletion", "user_id", userID, "error", err)
		return fmt.Errorf("failed to get existing profile: %w", err)
	}

	if existingProfile == nil {
		s.logger.Warn("Profile not found for deletion", "user_id", userID)
		return ErrProfileNotFound
	}

	// Delete from database
	err = s.profileRepo.DeleteProfile(ctx, userID)
	if err != nil {
		s.logger.Error("Failed to delete profile", "user_id", userID, "error", err)
		return fmt.Errorf("failed to delete profile: %w", err)
	}

	// Delete from cache
	if err := s.cacheRepo.DeleteCachedProfile(ctx, userID); err != nil {
		s.logger.Warn("Failed to delete cached profile", "user_id", userID, "error", err)
		// Non-critical error, continue
	}

	s.logger.Info("Profile deleted successfully", "user_id", userID)
	return nil
}

func (s *ProfileService) GetMultipleProfiles(ctx context.Context, userIDs []string) ([]*models.UserProfile, error) {
	s.logger.Debug("Getting multiple profiles", "user_ids", userIDs)

	profiles := make([]*models.UserProfile, 0, len(userIDs))
	missingFromCache := make([]string, 0)

	// Try to get from cache first
	for _, userID := range userIDs {
		cachedProfile, err := s.cacheRepo.GetCachedProfile(ctx, userID)
		if err != nil {
			s.logger.Warn("Failed to get profile from cache", "user_id", userID, "error", err)
			missingFromCache = append(missingFromCache, userID)
			continue
		}

		if cachedProfile != nil {
			profiles = append(profiles, cachedProfile)
		} else {
			missingFromCache = append(missingFromCache, userID)
		}
	}

	// If all profiles were in cache, return them
	if len(missingFromCache) == 0 {
		return profiles, nil
	}

	// Get missing profiles from database
	for _, userID := range missingFromCache {
		profile, err := s.profileRepo.GetProfileByUserID(ctx, userID)
		if err != nil {
			s.logger.Error("Failed to get profile from database", "user_id", userID, "error", err)
			continue
		}

		if profile != nil {
			profiles = append(profiles, profile)
			// Cache the profile for future requests
			if err := s.cacheRepo.CacheProfile(ctx, profile); err != nil {
				s.logger.Warn("Failed to cache profile", "user_id", userID, "error", err)
			}
		}
	}

	return profiles, nil
}
