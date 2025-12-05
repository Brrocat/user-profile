package main

import (
	"github.com/Brrocat/car-sharing-protos/proto/userprofile"
	"github.com/Brrocat/user-profile-service/internal/config"
	"github.com/Brrocat/user-profile-service/internal/handler"
	"github.com/Brrocat/user-profile-service/internal/repository/postgres"
	"github.com/Brrocat/user-profile-service/internal/repository/redis"
	"github.com/Brrocat/user-profile-service/internal/service"
	"github.com/Brrocat/user-profile-service/pkg/validation"
	"google.golang.org/grpc"
	"log"
	"log/slog"
	"net"
	"os"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Setup logger
	logger := setupLogger(cfg.Env)

	// Initialize repositories
	profileRepo, err := postgres.NewProfileRepository(cfg.DatabaseURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer profileRepo.Close()

	cacheRepo, err := redis.NewCacheRepository(cfg.RedisURL)
	if err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		os.Exit(1)
	}
	defer cacheRepo.Close()

	// Initialize utilities
	validator := validation.NewValidator()

	// Initialize service
	profileService := service.NewProfileService(profileRepo, cacheRepo, validator, logger)

	// Initialize gRPC handler
	profileHandler := handler.NewProfileHandler(profileService, logger)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		logger.Error("Failed to listen", "port", cfg.Port, "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	userprofile.RegisterUserProfileServiceServer(grpcServer, profileHandler)

	logger.Info("Starting user profile service", "port", cfg.Port, "env", cfg.Env)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Error("Failed to serve gRPC", "error", err)
		os.Exit(1)
	}
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "development":
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}))
	default:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}
	return logger
}
