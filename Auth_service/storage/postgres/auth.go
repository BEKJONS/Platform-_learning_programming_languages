package postgres

import (
	"Auth_service/pkg/models"
	"Auth_service/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type authRepo struct {
	db *sqlx.DB
}

func NewAuthStorage(db *sqlx.DB) storage.AuthStorage {
	return &authRepo{
		db: db,
	}
}

// Register a new user
func (a *authRepo) Register(ctx context.Context, req *models.RegisterRequest) (*models.RegisterResponse, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, username, password, role, phone_number)
		VALUES ($1, $2, $3, $4, $5, 'user', $6)
		ON CONFLICT (email) DO NOTHING
		RETURNING user_id, created_at`

	var user models.RegisterResponse
	err := a.db.QueryRowContext(ctx, query, req.FirstName, req.LastName, req.Email, req.Username, req.Password, req.PhoneNumber).
		Scan(&user.UserId, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Username = req.Username
	user.PhoneNumber = req.PhoneNumber

	return &user, nil
}

// Login a user
func (a *authRepo) Login(ctx context.Context, email string) (*models.LoginResponse, error) {
	res := &models.LoginResponse{}
	query := `SELECT user_id, password, role FROM Users WHERE email = $1 and deleted_at=0`

	err := a.db.QueryRowContext(ctx, query, email).Scan(&res.UserId, &res.Password, &res.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return res, nil
}

func (a *authRepo) GetUserByEmail(ctx context.Context, email string) (*models.GetProfileResponse, error) {
	query := `SELECT user_id ,created_at FROM users WHERE email = $1 AND deleted_at=0`

	var user models.GetProfileResponse

	err := a.db.QueryRowContext(ctx, query, email).Scan(&user.Id, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (a *authRepo) UpdatePassword(ctx context.Context, req *models.UpdatePasswordReq) error {
	query := `update users set password=$1 where user_id=$2 and deleted_at=0`

	result, err := a.db.ExecContext(ctx, query, req.Password, req.Id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

// Register a new user
func (a *authRepo) RegisterAdmin(ctx context.Context, pass string) error {
	query := `
		INSERT INTO users (first_name,last_name,username,email, password,role)
		VALUES ($1, $2, $3, $4, $5, 'admin')`

	err := a.db.QueryRowContext(ctx, query, "admin", "adminov", "admin", "admin@admin", pass)
	if err != nil {
		return err.Err()
	}

	return nil
}
