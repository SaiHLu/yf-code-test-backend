package gorm

import (
	"context"

	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	portrepository "codetest/internal/port/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) portrepository.UserRepository {
	return &userRepository{
		DB: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user *model.UserModel) error {
	return u.DB.WithContext(ctx).Create(user).Error
}

func (u *userRepository) Find(ctx context.Context, request *dto.QueryUserRequest) ([]*model.UserModel, int64, error) {
	var (
		users []*model.UserModel
		total int64
	)

	query := u.DB.WithContext(ctx).Model(&model.UserModel{})

	if request.Search != "" {
		query = query.Where("name LIKE ?", "%"+request.Search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Limit(request.PageSize).
		Offset((request.Page - 1) * request.PageSize).
		Order("created_at DESC").
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (u *userRepository) GetOneBy(ctx context.Context, column string, value string) (*model.UserModel, error) {
	var user model.UserModel
	if err := u.DB.WithContext(ctx).Where(column+" = ?", value).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *model.UserModel) error {
	return u.DB.WithContext(ctx).Where("id = ?", user.ID).Updates(user).Error
}

func (u *userRepository) DeleteOneBy(ctx context.Context, column string, value string) error {
	return u.DB.WithContext(ctx).Where(column+" = ?", value).Delete(&model.UserModel{}).Error
}
