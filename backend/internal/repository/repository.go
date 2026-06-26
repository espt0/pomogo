package repository

import "github.com/jackc/pgx/v5/pgxpool"

type Repository struct {
	DB *pgxpool.Pool
}

func NewRepository(DB *pgxpool.Pool) *Repository {
	return &Repository{DB: DB}
}
