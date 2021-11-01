package repository

import "github.com/jmoiron/sqlx"

type AuthAdminPostgres struct {
	db *sqlx.DB
}

func NewAuthAdminPostgres(db *sqlx.DB) *AuthAdminPostgres {
	return &AuthAdminPostgres{db: db}
}
