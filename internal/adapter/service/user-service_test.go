package service

import (
	"codetest/internal/adapter/api/dto"
	"codetest/internal/model"
	"codetest/mocks/repository"
	"context"
	"errors"
	"testing"

	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockUserRepo)

	ctx := context.Background()

	tests := []struct {
		name          string
		request       *dto.CreateUserRequest
		setupMock     func()
		expectedError bool
		err           error
	}{
		{
			name: "successful user creation",
			request: &dto.CreateUserRequest{
				Name:            "John Doe",
				Email:           "john@doe.com",
				Password:        "password",
				ConfirmPassword: "password",
			},
			setupMock: func() {
				mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(nil)
			},
			expectedError: false,
			err:           nil,
		},
		{
			name: "user creation success and password must be hashed",
			request: &dto.CreateUserRequest{
				Name:            "Jane Doe",
				Email:           "jane@doe.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func() {
				mockUserRepo.EXPECT().Create(ctx, gomock.Any()).DoAndReturn(func(ctx context.Context, user *model.UserModel) error {
					// Check if Password is plain text
					if user.Password == "password123" {
						t.Errorf("Password should be hashed, but got plain text: %s", user.Password)
					}

					// Compare the hashed password
					err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
					if err != nil {
						t.Errorf("Password hash verification failed: %v", err)
					}

					return nil
				})
			},
			expectedError: false,
			err:           nil,
		},
		{
			name: "user creation fails when email already exists",
			request: &dto.CreateUserRequest{
				Name:            "Duplicate User",
				Email:           "gg@doe.com",
				Password:        "password123",
				ConfirmPassword: "password123",
			},
			setupMock: func() {
				mockUserRepo.EXPECT().Create(ctx, gomock.Any()).Return(errors.New("user with this email already exists"))
			},
			expectedError: true,
			err:           errors.New("user with this email already exists"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()

			err := userService.Create(ctx, tt.request)

			if tt.expectedError {
				t.Logf("[%s]: Received expected error: %v", tt.name, err)

				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
