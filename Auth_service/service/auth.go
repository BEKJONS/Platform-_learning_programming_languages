package service

import (
	"Auth_service/pkg/config"
	"Auth_service/pkg/hashing"
	"Auth_service/pkg/models"
	"Auth_service/pkg/token"
	"Auth_service/storage"
	"context"
	"fmt"
	_ "github.com/badoux/checkmail"
	"github.com/pkg/errors"
	"log/slog"
)

type AuthService interface {
	Register(ctx context.Context, in models.RegisterRequest) (*models.RegisterResponse, error)
	Login(ctx context.Context, in models.LoginRequest) (*models.Token, error)
	GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error)
	UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error
	RegisterAdmin(ctx context.Context) error
}

func NewAuthService(st storage.AuthStorage, logger *slog.Logger) AuthService {
	return &authService{st, logger}
}

type authService struct {
	st  storage.AuthStorage
	log *slog.Logger
}

func (a *authService) Register(ctx context.Context, in models.RegisterRequest) (*models.RegisterResponse, error) {
	hash, err := hashing.HashPassword(in.Password)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return nil, err
	}

	in.Password = hash

	res, err := a.st.Register(ctx, &in)
	if err != nil {
		a.log.Error("Failed to register user", "error", err)
		return nil, err
	}

	return res, nil
}

func (a *authService) RegisterAdmin(ctx context.Context) error {

	hash, err := hashing.HashPassword(config.Load().ADMIN_PASSWORD)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return err
	}

	err = a.st.RegisterAdmin(ctx, hash)
	if err != nil {
		a.log.Error("Failed to register admin", "error", err)
		return err
	}

	return nil
}

func (a *authService) Login(ctx context.Context, in models.LoginRequest) (*models.Token, error) {
	res, err := a.st.Login(ctx, in.Email)
	if err != nil {
		a.log.Error("Failed to login", "error", err)
		return nil, err
	}

	check := hashing.CheckPasswordHash(res.Password, in.Password)
	if !check {
		a.log.Warn("Invalid password")
		return nil, errors.New("Invalid password")
	}

	accessToken, err := token.GenerateAccessToken(res.UserId, res.Role, in.Email)
	if err != nil {
		a.log.Error("Failed to generate access token", "error", err)
		return nil, err
	}

	refreshToken, err := token.GenerateRefreshToken(res.UserId, res.Role, in.Email)
	if err != nil {
		a.log.Error("Failed to generate refresh token", "error", err)
		return nil, err
	}

	tk := &models.Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		UserID:       res.UserId,
	}

	return tk, nil
}
func (a *authService) GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error) {
	a.log.Info("Getting user user by email")
	res, err := a.st.GetUserByEmail(ctx, email)
	if err != nil {
		a.log.Error(err.Error())
		return nil, errors.Wrap(err, "failed to get user by email")
	}
	return res, nil
}

func (a *authService) UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error {
	hash, err := hashing.HashPassword(req.Password)
	if err != nil {
		a.log.Error("Failed to hash password", "error", err)
		return err
	}

	req.Password = hash

	err = a.st.UpdatePassword(ctx, req)
	if err != nil {
		a.log.Error(fmt.Sprintf("Error update pasword: %v", err))
		return err
	}
	return nil
}
