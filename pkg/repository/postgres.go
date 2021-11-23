package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersAdminTable = "dbo.t_admins"
	usersTable      = "dbo.t_users"
	countryTable    = "dbo.t_countries"
	airportTable    = "dbo.t_airports"
	aircraftTable   = "dbo.t_aircraft_models"
	airlineTable    = "dbo.t_airlines"
	flightTable     = "dbo.t_flights"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
        cfg.Username,
        cfg.Password,
        cfg.Host,
        cfg.Port,
        cfg.DBName)
	db, err := sqlx.Open("postgres", url)
	    
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
