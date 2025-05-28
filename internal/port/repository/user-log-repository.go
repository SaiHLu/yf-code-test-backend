package portrepository

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	"context"
)

type UserLogRepository interface {
	Create(ctx context.Context, userLog *model.UserLogModel) error
	Find(ctx context.Context, request *dto.QueryUserLogRequest) ([]*model.UserLogModel, int64, error)
}
