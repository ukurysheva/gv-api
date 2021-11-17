package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type AirlinePostgres struct {
	db *sqlx.DB
}

func NewAirlinePostgres(db *sqlx.DB) *AirlinePostgres {
	return &AirlinePostgres{db: db}
}

func (r *AirlinePostgres) Create(userId int, airline gvapi.Airline) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (airline_iata, airline_icao, airline_country_id, airline_active_flg,change_dttm ) ` +
		`VALUES ($1, $2, $3, $4, $5) RETURNING airline_id`
	createairlineQuery := fmt.Sprintf(q, airlineTable)

	row := tx.QueryRow(createairlineQuery, airline.Iata, airline.Icao, airline.CountryId, airline.Active, time.Now())

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AirlinePostgres) GetAll() ([]gvapi.Airline, error) {
	var airlines []gvapi.Airline

	query := fmt.Sprintf("SELECT * FROM %s tl", airlineTable)
	err := r.db.Select(&airlines, query)

	return airlines, err
}

func (r *AirlinePostgres) GetById(airlineId int) (gvapi.Airline, error) {
	var airline gvapi.Airline

	query := fmt.Sprintf(`SELECT * FROM %s 
	                      WHERE airline_id = $1`, airlineTable)
	if err := r.db.Get(&airline, query, airlineId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return airline, errors.New("Nothing found")
		case nil:
			return airline, nil
		default:
			return airline, err
		}
	}

	return airline, nil
}
