package handler

import (
	"context"
	"github.com/Brrocat/car-sharing-protos/proto/userprofile"
	"github.com/Brrocat/user-profile-service/internal/models"
	"github.com/Brrocat/user-profile-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
)

type ProfileHandler struct {
	userprofile.UnimplementedUserProfileServiceServer
	profileService *service.ProfileService
	logger         *slog.Logger
}

func NewProfileHandler(profileService *service.ProfileService, logger *slog.Logger) *ProfileHandler {
	return &ProfileHandler{
		profileService: profileService,
		logger:         logger,
	}
}

func (h *ProfileHandler) GetUserProfile(ctx context.Context, req *userprofile.GetUserProfileRequest) (*userprofile.GetUserProfileResponse, error) {
	h.logger.Debug("GetUserProfile request received", "user_id", req.UserId)

	profile, err := h.profileService.GetUserProfile(ctx, req.UserId)
	if err != nil {
		h.logger.Warn("GetUserProfile failed", "user_id", req.UserId, "error", err)

		switch err {
		case service.ErrProfileNotFound:
			return nil, status.Error(codes.NotFound, "profile not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	h.logger.Debug("GetUserProfile successful", "user_id", req.UserId)

	return &userprofile.GetUserProfileResponse{
		Profile: &userprofile.UserProfile{
			UserId:      profile.UserID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Phone:       profile.Phone,
			DateOfBirth: profile.DateOfBirth,
			AvatarUrl:   profile.AvatarURL,
			CreatedAt:   timestamppb.New(profile.CreatedAt),
			UpdatedAt:   timestamppb.New(profile.UpdatedAt),
		},
	}, nil
}

func (h *ProfileHandler) CreateUserProfile(ctx context.Context, req *userprofile.CreateUserProfileRequest) (*userprofile.CreateUserProfileResponse, error) {
	h.logger.Debug("CreateUserProfile request received", "user_id", req.UserId)

	createReq := &models.CreateProfileRequest{
		UserID:      req.UserId,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Phone:       req.Phone,
		DateOfBirth: req.DateOfBirth,
	}

	profile, err := h.profileService.CreateUserProfile(ctx, createReq)
	if err != nil {
		h.logger.Warn("CreateUserProfile failed", "user_id", req.UserId, "error", err)

		switch err {
		case service.ErrProfileAlreadyExists:
			return nil, status.Error(codes.AlreadyExists, "profile already exists")
		case service.ErrInvalidData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	h.logger.Info("CreateUserProfile successful", "user_id", req.UserId, "profile_id", profile.ID)

	return &userprofile.CreateUserProfileResponse{
		Profile: &userprofile.UserProfile{
			UserId:      profile.UserID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Phone:       profile.Phone,
			DateOfBirth: profile.DateOfBirth,
			AvatarUrl:   profile.AvatarURL,
			CreatedAt:   timestamppb.New(profile.CreatedAt),
			UpdatedAt:   timestamppb.New(profile.UpdatedAt),
		},
	}, nil
}

func (h *ProfileHandler) UpdateUserProfile(ctx context.Context, req *userprofile.UpdateUserProfileRequest) (*userprofile.UpdateUserProfileResponse, error) {
	h.logger.Debug("UpdateUserProfile request received", "user_id", req.UserId)

	updateReq := &models.UpdateProfileRequest{
		UserID:         req.UserId,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Phone:          req.Phone,
		DateOfBirth:    req.DateOfBirth,
		AvatarURL:      req.AvatarUrl,
		Address:        req.Address,
		City:           req.City,
		Country:        req.Country,
		PostalCode:     req.PostalCode,
		DrivingLicense: req.DrivingLicense,
	}

	profile, err := h.profileService.UpdateUserProfile(ctx, req.UserId, updateReq)
	if err != nil {
		h.logger.Warn("UpdateUserProfile failed", "user_id", req.UserId, "error", err)

		switch err {
		case service.ErrProfileNotFound:
			return nil, status.Error(codes.NotFound, "profile not found")
		case service.ErrInvalidData:
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	h.logger.Info("UpdateUserProfile successful", "user_id", req.UserId, "profile_id", profile.ID)

	return &userprofile.UpdateUserProfileResponse{
		Profile: &userprofile.UserProfile{
			UserId:      profile.UserID,
			FirstName:   profile.FirstName,
			LastName:    profile.LastName,
			Phone:       profile.Phone,
			DateOfBirth: profile.DateOfBirth,
			AvatarUrl:   profile.AvatarURL,
			CreatedAt:   timestamppb.New(profile.CreatedAt),
			UpdatedAt:   timestamppb.New(profile.UpdatedAt),
		},
	}, nil
}

func (h *ProfileHandler) DeleteUserProfile(ctx context.Context, req *userprofile.DeleteUserProfileRequest) (*userprofile.DeleteUserProfileResponse, error) {
	h.logger.Debug("DeleteUserProfile request received", "user_id", req.UserId)

	err := h.profileService.DeleteUserProfile(ctx, req.UserId)
	if err != nil {
		h.logger.Warn("DeleteUserProfile failed", "user_id", req.UserId, "error", err)

		switch err {
		case service.ErrProfileNotFound:
			return nil, status.Error(codes.NotFound, "profile not found")
		default:
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	h.logger.Info("DeleteUserProfile successful", "user_id", req.UserId)

	return &userprofile.DeleteUserProfileResponse{
		Success: true,
	}, nil
}
