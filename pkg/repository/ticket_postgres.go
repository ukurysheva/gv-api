package repository

import "github.com/jmoiron/sqlx"

type TicketPostgres struct {
	db *sqlx.DB
}

func NewTicketPostgres(db *sqlx.DB) *TicketPostgres {
	return &TicketPostgres{db: db}
}
