package portrepository

import (
	"context"

	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.UserModel) error
	Find(ctx context.Context, request *dto.QueryUserRequest) ([]*model.UserModel, int64, error)
	GetOneBy(ctx context.Context, column, value string) (*model.UserModel, error)
	Update(ctx context.Context, user *model.UserModel) error
	DeleteOneBy(ctx context.Context, column, value string) error
}
