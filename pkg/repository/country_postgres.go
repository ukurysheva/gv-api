package repository

import "github.com/jmoiron/sqlx"

type CountryPostgres struct {
	db *sqlx.DB
}

func NewCountryPostgres(db *sqlx.DB) *CountryPostgres {
	return &CountryPostgres{db: db}
}
