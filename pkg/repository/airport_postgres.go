package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type AirportPostgres struct {
	db *sqlx.DB
}

func NewAirportPostgres(db *sqlx.DB) *AirportPostgres {
	return &AirportPostgres{db: db}
}
func (r *AirportPostgres) Create(userId int, airport gvapi.Airport) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (airport_type, airport_name, airport_iso_country_id, airport_iso_region, airport_municipality, ` +
		`airport_iata_code, airport_home_link, visa_flg, quarantine_flg, covid_test_flg, lockdown_flg,change_dttm ) ` +
		`VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  $9, $10, $11, $12) RETURNING airport_id`
	createAirportQuery := fmt.Sprintf(q, airportTable)

	row := tx.QueryRow(createAirportQuery, airport.Type, airport.Name, airport.CountryId, airport.Region, airport.Municipality,
		airport.IataCode, airport.HomeLink, airport.Visa, airport.Quarantine, airport.CovidTest, airport.LockDown, time.Now())

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AirportPostgres) GetAll() ([]gvapi.Airport, error) {
	var airports []gvapi.Airport

	query := fmt.Sprintf(`SELECT airport_id, airport_type, airport_name, airport_iso_country_id, `+
		`COALESCE(airport_iata_code, '') as airport_iata_code, COALESCE(airport_home_link, '') as airport_home_link,`+
		`COALESCE(airport_iso_region, '') as airport_iso_region, COALESCE(airport_municipality, '') as airport_municipality, `+
		`visa_flg, quarantine_flg, covid_test_flg, lockdown_flg,change_dttm FROM %s`, airportTable)
	err := r.db.Select(&airports, query)

	return airports, err
}

func (r *AirportPostgres) GetById(airportId int) (gvapi.Airport, error) {
	var airport gvapi.Airport

	query := fmt.Sprintf(`SELECT airport_id, airport_type, airport_name, airport_iso_country_id, airport_iso_region, airport_municipality, `+
		`COALESCE(airport_iata_code, '') as airport_iata_code, COALESCE(airport_home_link, '') as airport_home_link, visa_flg, quarantine_flg, `+
		`covid_test_flg, lockdown_flg,change_dttm FROM %s 
	                      WHERE airport_id = $1`, airportTable)
	if err := r.db.Get(&airport, query, airportId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return airport, errors.New("Nothing found")
		case nil:
			return airport, nil
		default:
			return airport, err
		}
	}

	return airport, nil
}
