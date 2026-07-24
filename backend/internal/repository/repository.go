package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/espt0/pomogo/internal/apperrors"
	"github.com/espt0/pomogo/internal/model"

	"github.com/gofrs/uuid/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
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

	if err := r.DB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = $1 AND active = true)", email).Scan(&exists); err != nil {
		return false, fmt.Errorf("erro ao verificar email: %w", err)
	}

	return exists, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)

	if err := r.DB.QueryRow(ctx, "SELECT id, name, email, password_hash, active, created_at, updated_at FROM users WHERE email = $1 AND active = true", email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Active, &user.CreatedAt, &user.UpdatedAt); err != nil {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrEmailAlreadyExists
		}

		return fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return nil
}

func (r *Repository) FindRefreshTokenByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	rT := new(model.RefreshToken)
	err := r.DB.QueryRow(ctx,
		"SELECT rt.id, rt.user_id, rt.token_hash, rt.expires_at, rt.revoked, rt.created_at FROM refresh_token rt JOIN users u ON rt.user_id = u.id WHERE token_hash = $1 AND u.active = true",
		hash,
	).Scan(&rT.ID, &rT.UserID, &rT.TokenHash, &rT.ExpiresAt, &rT.Revoked, &rT.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("erro ao refresh token usuário: %w", err)
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

func (r *Repository) GetUserByID(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	user := new(model.User)

	if err := r.DB.QueryRow(ctx, "SELECT id, name, email, active, created_at, updated_at FROM users WHERE id = $1 AND active = true", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Active, &user.CreatedAt, &user.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}

		return nil, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, name, email *string, userID uuid.UUID) error {
	cmdTag, err := r.DB.Exec(ctx, "UPDATE users SET name = COALESCE($1, name), email = COALESCE($2, email), updated_at = NOW() WHERE id = $3 AND active = true", name, email, userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return apperrors.ErrEmailAlreadyExists
		}

		return fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *Repository) GetHashPasswordByID(ctx context.Context, userID uuid.UUID) (string, error) {
	var hashPassword string

	if err := r.DB.QueryRow(ctx, "SELECT password_hash FROM users WHERE id = $1 AND active = true", userID).Scan(&hashPassword); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", apperrors.ErrUserNotFound
		}

		return "", fmt.Errorf("erro ao buscar o hash: %w", err)
	}

	return hashPassword, nil
}

func (r *Repository) UpdatePassword(ctx context.Context, passwordHash string, userID uuid.UUID) error {
	_, err := r.DB.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE id = $2 AND active = true", passwordHash, userID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar password: %w", err)
	}

	return nil
}

func (r *Repository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	cmdTag, err := r.DB.Exec(ctx, "UPDATE users SET active = false WHERE id = $1 AND active = true", userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *Repository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	_, err := r.DB.Exec(ctx, "UPDATE refresh_token SET revoked = true WHERE user_id = $1", userID)
	if err != nil {
		return fmt.Errorf("erro ao revogar todos os tokens do usuário: %w", err)
	}

	return nil
}

func (r *Repository) CreateTask(ctx context.Context, input *model.Task) error {
	_, err := r.DB.Exec(ctx, "INSERT INTO tasks (id, user_id, title, description, status, active, estimated_pomodoros, completed_pomodoros, created_at, updated_at) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", input.ID, input.UserID, input.Title, input.Description, input.Status, input.Active, input.EstimatedPomodoros, input.CompletedPomodoros, input.CreatedAt, input.UpdatedAt)
	if err != nil {
		return fmt.Errorf("erro ao criar task: %w", err)
	}

	return nil
}

func (r *Repository) GetTasksByID(ctx context.Context, userID uuid.UUID) ([]model.Task, error) {
	rows, err := r.DB.Query(ctx, "SELECT id, user_id, title, description, status, estimated_pomodoros, completed_pomodoros, created_at, updated_at FROM tasks WHERE user_id = $1 AND active = true", userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar tasks: %w", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.EstimatedPomodoros, &task.CompletedPomodoros, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) (*model.Task, error) {
	task := new(model.Task)

	if err := r.DB.QueryRow(ctx, "SELECT id, user_id, title, description, status, estimated_pomodoros, completed_pomodoros, created_at, updated_at FROM tasks WHERE user_id = $1 AND id = $2 AND active = true", userID, taskID).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.EstimatedPomodoros, &task.CompletedPomodoros, &task.CreatedAt, &task.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("erro ao buscar task: %w", err)
	}

	return task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID, input *model.UpdateTaskRequest) error {
	cmdTag, err := r.DB.Exec(ctx, "UPDATE tasks SET title = COALESCE($1, title), description = COALESCE($2, description), status = COALESCE($3, status), estimated_pomodoros = COALESCE($4, estimated_pomodoros), completed_pomodoros = COALESCE($5, completed_pomodoros), updated_at = NOW() WHERE id = $6 AND user_id = $7 AND active = true", input.Title, input.Description, input.Status, input.EstimatedPomodoros, input.CompletedPomodoros, taskID, userID)
	if err != nil {
		return fmt.Errorf("erro ao atualizar task: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, userID uuid.UUID, taskID uuid.UUID) error {
	cmdTag, err := r.DB.Exec(ctx, "UPDATE tasks SET active = false WHERE id = $1 AND user_id = $2 AND active = true", taskID, userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar task: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return apperrors.ErrNotFound
	}
	return nil
}
