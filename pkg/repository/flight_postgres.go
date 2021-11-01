package repository

import "github.com/jmoiron/sqlx"

type FlightPostgres struct {
	db *sqlx.DB
}

func NewFlightPostgres(db *sqlx.DB) *FlightPostgres {
	return &FlightPostgres{db: db}
}
