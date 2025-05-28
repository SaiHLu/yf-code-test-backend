package portservice

import (
	"context"

	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, request *dto.CreateUserRequest) error
	Find(ctx context.Context, request *dto.QueryUserRequest) ([]*model.UserModel, int64, error)
	GetOneByID(ctx context.Context, id uuid.UUID) (*model.UserModel, error)
	GetOneByEmail(ctx context.Context, email string) (*model.UserModel, error)
	Update(ctx context.Context, id uuid.UUID, request *dto.UpdateUserRequest) error
	DeleteOneByID(ctx context.Context, id uuid.UUID) error
}
