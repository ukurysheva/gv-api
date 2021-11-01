package repository

import "github.com/jmoiron/sqlx"

type AuthClientPostgres struct {
	db *sqlx.DB
}

func NewAuthClientPostgres(db *sqlx.DB) *AuthClientPostgres {
	return &AuthClientPostgres{db: db}
}
