package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/espt0/pomogo/internal/apperrors"
	"github.com/espt0/pomogo/internal/model"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(DB *pgxpool.Pool) *Repository {
	return &Repository{DB: DB}
}

func (r *Repository) EmailExists(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar email: %w", err)
	}

	return exists, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)

	err := r.DB.QueryRow(ctx, "SELECT id, name, email, password_hash, active, created_at, updated_at FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Active, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *Repository) CreateRefreshToken(ctx context.Context, input model.RefreshToken) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO refresh_token (id, user_id, token_hash, expires_at, revoked, created_at) VALUES($1, $2, $3, $4, $5, $6)", input.ID, input.UserID, input.TokenHash, input.ExpiresAt, input.Revoked, input.CreatedAt)
	if err != nil {
		return fmt.Errorf("erro ao salvar refresh token: %w", err)
	}

	return nil
}

func (r *Repository) CreateUser(ctx context.Context, user model.User) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO users (id, name, email, password_hash, active, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7)", user.ID, user.Name, user.Email, user.PasswordHash, user.Active, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}

func (r *Repository) FindActiveRefreshTokenByHash(ctx context.Context, hashRefreshToken string) (*model.RefreshToken, error) {
	rT := new(model.RefreshToken)

	err := r.DB.QueryRow(ctx, "SELECT id, user_id, token_hash, expires_at, revoked, created_at FROM refresh_token WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()", hashRefreshToken).Scan(&rT.ID, &rT.UserID, &rT.TokenHash, &rT.ExpiresAt, &rT.Revoked, &rT.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrInvalidRefreshToken
		}

		return nil, fmt.Errorf("erro ao buscar refresh token: %w", err)
	}

	return rT, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, refreshTokenID uuid.UUID) error {

	_, err := r.DB.Exec(ctx, "UPDATE refresh_token SET revoked = true WHERE id = $1", refreshTokenID)
	if err != nil {
		return fmt.Errorf("erro ao revogar refresh token: %w", err)
	}

	return nil
}

func (r *Repository) RevokeRefreshTokenByHash(ctx context.Context, tokenHash string) error {
	_, err := r.DB.Exec(ctx, "UPDATE refresh_token SET revoked = true WHERE token_hash = $1", tokenHash)
	if err != nil {
		return fmt.Errorf("erro ao revogar refresh token por hash: %w", err)
	}
	return nil
}
