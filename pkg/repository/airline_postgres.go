package repository

import "github.com/jmoiron/sqlx"

type AirlinePostgres struct {
	db *sqlx.DB
}

func NewAirlinePostgres(db *sqlx.DB) *AirlinePostgres {
	return &AirlinePostgres{db: db}
}
