package repository

import "github.com/jmoiron/sqlx"

type AirportPostgres struct {
	db *sqlx.DB
}

func NewAirportPostgres(db *sqlx.DB) *AirportPostgres {
	return &AirportPostgres{db: db}
}
