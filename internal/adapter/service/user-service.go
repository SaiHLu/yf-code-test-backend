package service

import (
	"context"

	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	portrepository "codetest/internal/port/repository"
	portservice "codetest/internal/port/service"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository portrepository.UserRepository
}

func NewUserService(userRepository portrepository.UserRepository) portservice.UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (u *userService) Create(ctx context.Context, request *dto.CreateUserRequest) error {
	user := &model.UserModel{
		Name:  request.Name,
		Email: request.Email,
	}

	passBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(passBytes)

	return u.userRepository.Create(ctx, user)
}

func (u *userService) Find(ctx context.Context, request *dto.QueryUserRequest) ([]*model.UserModel, int64, error) {
	users, total, err := u.userRepository.Find(ctx, request)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *userService) GetOneByID(ctx context.Context, id uuid.UUID) (*model.UserModel, error) {
	return u.userRepository.GetOneBy(ctx, "id", id.String())
}

func (u *userService) GetOneByEmail(ctx context.Context, email string) (*model.UserModel, error) {
	return u.userRepository.GetOneBy(ctx, "email", email)
}

func (u *userService) DeleteOneByID(ctx context.Context, id uuid.UUID) error {
	return u.userRepository.DeleteOneBy(ctx, "id", id.String())
}

func (u *userService) Update(ctx context.Context, id uuid.UUID, request *dto.UpdateUserRequest) error {
	user := &model.UserModel{
		ID:    id,
		Name:  request.Name,
		Email: request.Email,
	}

	if len(request.Password) > 0 {
		passBytes, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(passBytes)
	}

	return u.userRepository.Update(ctx, user)
}
