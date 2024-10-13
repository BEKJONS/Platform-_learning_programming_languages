package storage

import (
	"Auth_service/pkg/models"
	"context"
)

type AuthStorage interface {
	Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(ctx context.Context, email string) (*models.LoginResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error)

	UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error
	RegisterAdmin(ctx context.Context, pass string) error
}
