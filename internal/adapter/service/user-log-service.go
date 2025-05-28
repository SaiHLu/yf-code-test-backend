package service

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	portrepository "codetest/internal/port/repository"
	portservice "codetest/internal/port/service"
	"context"
)

type userLogService struct {
	userLogRepository portrepository.UserLogRepository
}

func NewUserLogService(userLogRepository portrepository.UserLogRepository) portservice.UserLogService {
	return &userLogService{
		userLogRepository: userLogRepository,
	}
}

// Create implements portservice.UserLogService.
func (u *userLogService) Create(ctx context.Context, userLog *model.UserLogModel) error {
	return u.userLogRepository.Create(ctx, userLog)
}

// Find implements portservice.UserLogService.
func (u *userLogService) Find(ctx context.Context, request *dto.QueryUserLogRequest) ([]*model.UserLogModel, int64, error) {
	userLogs, total, err := u.userLogRepository.Find(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return userLogs, total, nil
}
