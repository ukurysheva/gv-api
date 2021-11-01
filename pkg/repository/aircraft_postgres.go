package repository

import "github.com/jmoiron/sqlx"

type AircraftPostgres struct {
	db *sqlx.DB
}

func NewAircraftPostgres(db *sqlx.DB) *AircraftPostgres {
	return &AircraftPostgres{db: db}
}
