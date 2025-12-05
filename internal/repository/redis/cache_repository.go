package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Brrocat/user-profile-service/internal/models"
	"github.com/redis/go-redis/v9"
	"time"
)

type CacheRepository struct {
	client *redis.Client
	ttl    time.Duration
}

func NewCacheRepository(redisURL string) (*CacheRepository, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opts)

	// Test connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &CacheRepository{
		client: client,
		ttl:    1 * time.Hour, // default TTL
	}, nil
}

func (r *CacheRepository) Close() {
	if r.client != nil {
		r.client.Close()
	}
}

func (r *CacheRepository) SetTTL(ttl time.Duration) {
	r.ttl = ttl
}

func (r *CacheRepository) CacheProfile(ctx context.Context, profile *models.UserProfile) error {
	key := fmt.Sprintf("user_profile:%s", profile.UserID)

	profileJSON, err := json.Marshal(profile)
	if err != nil {
		return fmt.Errorf("failed to marshal profile: %w", err)
	}

	err = r.client.Set(ctx, key, profileJSON, r.ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to cache profile: %w", err)
	}

	return nil
}

func (r *CacheRepository) GetCachedProfile(ctx context.Context, userID string) (*models.UserProfile, error) {
	key := fmt.Sprintf("user_profile:%s", userID)

	profileJSON, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cached profile: %w", err)
	}

	var profile models.UserProfile
	err = json.Unmarshal([]byte(profileJSON), &profile)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile: %w", err)
	}

	return &profile, nil
}

func (r *CacheRepository) DeleteCachedProfile(ctx context.Context, userID string) error {
	key := fmt.Sprintf("user_profile:%s", userID)
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete cached profile: %w", err)
	}
	return nil
}

func (r *CacheRepository) CacheProfileList(ctx context.Context, userIDs []string, profiles []*models.UserProfile) error {
	pipeline := r.client.Pipeline()

	for i, userID := range userIDs {
		if i < len(profiles) && profiles[i] != nil {
			key := fmt.Sprintf("user_profile:%s", userID)
			profileJSON, err := json.Marshal(profiles[i])
			if err != nil {
				return fmt.Errorf("failed to marshal profile: %w", err)
			}
			pipeline.Set(ctx, key, profileJSON, r.ttl)
		}
	}

	_, err := pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to cache profile list: %w", err)
	}

	return nil
}
