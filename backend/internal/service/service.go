package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
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

	user, err := s.repo.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return "", "", apperrors.ErrInvalidEmailOrPassword
		}

		return "", "", fmt.Errorf("buscando usuario no login: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
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

	if err = s.repo.CreateRefreshToken(ctx, refreshTokenSave); err != nil {
		return "", "", fmt.Errorf("salvando refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (s *Service) RotateTokens(ctx context.Context, refreshToken string) (string, string, error) {

	hashRefreshToken := auth.HashRefreshToken(refreshToken)

	rT, err := s.repo.FindRefreshTokenByHash(ctx, hashRefreshToken)
	if err != nil {
		return "", "", apperrors.ErrInvalidRefreshToken
	}
	if rT.Revoked || rT.ExpiresAt.Before(time.Now()) {
		if revokeErr := s.repo.RevokeAllUserTokens(ctx, rT.UserID); revokeErr != nil {
			slog.Error("falha ao revogar tokens em possível reuso de refresh token",
				"userID", rT.UserID, "err", revokeErr)
		}

		return "", "", apperrors.ErrInvalidRefreshToken
	}

	newAccessToken, err := auth.GenerateAccessToken(rT.UserID)
	if err != nil {
		return "", "", fmt.Errorf("gerando access token: %w", err)
	}
	newRefreshToken, newHashRefreshToken, err := auth.GenerateRefreshToken()
	if err != nil {
		return "", "", fmt.Errorf("gerando refresh token: %w", err)
	}

	if err = s.repo.RevokeRefreshToken(ctx, rT.ID); err != nil {
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
	if err = s.repo.CreateRefreshToken(ctx, newRefreshTokenSave); err != nil {
		return "", "", fmt.Errorf("salvando refresh token: %w", err)
	}

	return newAccessToken, newRefreshToken, nil
}

func (s *Service) LogoutUser(ctx context.Context, refreshToken string) error {

	hashRefreshToken := auth.HashRefreshToken(refreshToken)

	if err := s.repo.RevokeRefreshTokenByHash(ctx, hashRefreshToken); err != nil {
		return fmt.Errorf("revogando token no logout: %w", err)
	}

	return nil
}

func (s *Service) GetByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {

	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, fmt.Errorf("buscando usuario no banco: %w", err)
	}

	return user, nil
}

func (s *Service) UpdateUser(ctx context.Context, userID uuid.UUID, input *model.UpdateUserRequest) error {
	var newEmail *string
	var newName *string

	if input.Email != nil {
		e := strings.ToLower(strings.TrimSpace(*input.Email))
		newEmail = &e
	}
	if input.Name != nil {
		n := strings.TrimSpace(*input.Name)
		newName = &n
	}

	if err := s.repo.UpdateUser(ctx, newName, newEmail, userID); err != nil {
		return fmt.Errorf("atualizando usuário: %w", err)
	}

	return nil
}

func (s *Service) UpdatePasswordUser(ctx context.Context, userID uuid.UUID, input *model.UpdatePasswordUserRequest) error {

	oldHash, err := s.repo.GetHashPasswordByID(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.ErrUserNotFound
		}

		return fmt.Errorf("buscando hash no banco: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldHash), []byte(input.OldPassword)); err != nil {
		return apperrors.ErrInvalidEmailOrPassword
	}

	newPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return fmt.Errorf("gerando hash da senha: %w", err)
	}

	if err = s.repo.UpdatePassword(ctx, string(newPassword), userID); err != nil {
		return fmt.Errorf("atualizando password: %w", err)
	}

	if err = s.RevokeAllUserSessions(ctx, userID); err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := s.RevokeAllUserSessions(ctx, userID); err != nil {
		return err
	}

	if err := s.repo.DeleteUser(ctx, userID); err != nil {
		return fmt.Errorf("deletando usuário: %w", err)
	}

	return nil
}

func (s *Service) RevokeAllUserSessions(ctx context.Context, userID uuid.UUID) error {
	if err := s.repo.RevokeAllUserTokens(ctx, userID); err != nil {
		return fmt.Errorf("revogando sessões: %w", err)
	}
	return nil
}

func (s *Service) CreateTask(ctx context.Context, userID uuid.UUID, input *model.CreateTaskRequest) (*model.Task, error) {
	id, err := s.gen.NewV7()
	if err != nil {
		return nil, fmt.Errorf("gerando id: %w", err)
	}

	task := &model.Task{
		ID:                 id,
		UserID:             userID,
		Title:              input.Title,
		Description:        input.Description,
		Status:             "pending",
		Active:             true,
		EstimatedPomodoros: input.EstimatedPomodoros,
		CompletedPomodoros: 0,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("criando tarefa: %w", err)
	}

	return task, nil
}

func (s *Service) GetTasks(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	tasks, err := s.repo.GetTasksByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("buscando tarefas no banco: %w", err)
	}

	if tasks == nil {
		return []model.Task{}, nil
	}

	return tasks, nil
}

func (s *Service) GetTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*model.Task, error) {
	task, err := s.repo.GetTaskByID(ctx, userID, taskID)
	if err != nil {
		return nil, fmt.Errorf("buscando tarefa no banco: %w", err)
	}

	return task, nil
}

func (s *Service) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, input *model.UpdateTaskRequest) error {

	if err := s.repo.UpdateTask(ctx, userID, taskID, input); err != nil {
		return fmt.Errorf("atualizando tarefa: %w", err)
	}

	return nil
}

func (s *Service) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {

	if err := s.repo.DeleteTask(ctx, userID, taskID); err != nil {
		return fmt.Errorf("deletando tarefa: %w", err)
	}

	return nil
}
