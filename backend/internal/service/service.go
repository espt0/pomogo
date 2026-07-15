package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/espt0/pomogo/internal/apperrors"
	"github.com/espt0/pomogo/internal/auth"
	"github.com/espt0/pomogo/internal/model"
	"github.com/espt0/pomogo/internal/repository"

	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
)

// DI
type Service struct {
	repo *repository.Repository
	gen  uuid.Generator
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo, gen: uuid.DefaultGenerator}
}

func (s *Service) CreateUser(ctx context.Context, input *model.CreateUserRequest) error {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	name := strings.TrimSpace(input.Name)

	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return fmt.Errorf("checando email: %w", err)
	}
	if exists {
		return apperrors.ErrEmailAlreadyExists
	}

	id, err := s.gen.NewV7()
	if err != nil {
		return fmt.Errorf("gerando id: %w", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gerando hash da senha: %w", err)
	}

	user := model.User{
		ID:           id,
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Active:       true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("criando usuário: %w", err)
	}

	return nil
}

func (s *Service) LoginUser(ctx context.Context, input *model.LoginRequest) (string, string, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	password := strings.TrimSpace(input.Password)

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", "", apperrors.ErrInvalidEmailOrPassword
		}

		return "", "", fmt.Errorf("buscando usuario no login: %w", err)
	}

	if !user.Active {
		return "", "", apperrors.ErrInvalidEmailOrPassword
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", apperrors.ErrInvalidEmailOrPassword
	}

	accessToken, err := auth.GenerateAccessToken(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("gerando access token: %w", err)
	}
	refreshToken, hashRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("gerando refresh token: %w", err)
	}

	id, err := s.gen.NewV7()
	if err != nil {
		return "", "", fmt.Errorf("gerando id: %w", err)
	}

	refreshTokenSave := model.RefreshToken{
		ID:        id,
		UserID:    user.ID,
		TokenHash: hashRefreshToken,
		ExpiresAt: time.Now().Add(43200 * time.Minute),
		Revoked:   false,
		CreatedAt: time.Now(),
	}

	err = s.repo.CreateRefreshToken(ctx, refreshTokenSave)
	if err != nil {
		return "", "", fmt.Errorf("salvando refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RotateTokens(ctx context.Context, refreshToken string) (string, string, error) {

	hashRefreshToken := auth.HashRefreshToken(refreshToken)

	rT, err := s.repo.FindActiveRefreshTokenByHash(ctx, hashRefreshToken)
	if err != nil {
		if errors.Is(err, apperrors.ErrInvalidRefreshToken) {
			// Adicionar depois 'RevokeAllUserTokens'
			return "", "", apperrors.ErrInvalidRefreshToken
		}

		return "", "", fmt.Errorf("validando refresh token: %w", err)
	}

	newAccessToken, err := auth.GenerateAccessToken(rT.UserID)
	if err != nil {
		return "", "", fmt.Errorf("gerando access token: %w", err)
	}
	newRefreshToken, newHashRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("gerando refresh token: %w", err)
	}

	err = s.repo.RevokeRefreshToken(ctx, rT.ID)
	if err != nil {
		return "", "", fmt.Errorf("revogando token antigo: %w", err)
	}

	id, err := s.gen.NewV7()
	if err != nil {
		return "", "", fmt.Errorf("gerando id: %w", err)
	}

	newRefreshTokenSave := model.RefreshToken{
		ID:        id,
		UserID:    rT.UserID,
		TokenHash: newHashRefreshToken,
		ExpiresAt: time.Now().Add(43200 * time.Minute),
		Revoked:   false,
		CreatedAt: time.Now(),
	}
	err = s.repo.CreateRefreshToken(ctx, newRefreshTokenSave)
	if err != nil {
		return "", "", fmt.Errorf("salvando refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) LogoutUser(ctx context.Context, refreshToken string) error {

	hashRefreshToken := auth.HashRefreshToken(refreshToken)

	err := s.repo.RevokeRefreshTokenByHash(ctx, hashRefreshToken)
	if err != nil {
		return fmt.Errorf("revogando token no logout: %w", err)
	}

	return nil
}
