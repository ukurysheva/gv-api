package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type AircraftPostgres struct {
	db *sqlx.DB
}

func NewAircraftPostgres(db *sqlx.DB) *AircraftPostgres {
	return &AircraftPostgres{db: db}
}

func (r *AircraftPostgres) Create(userId int, aircraft gvapi.Aircraft) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (aircraft_iata_code, aircraft_model_name, aircraft_model_manufacturer, aircraft_model_type` +
		`, aircraft_icaic_code, aircraft_model_wing_type,  economy_class_flg,  pr_economy_class_flg, business_class_flg, first_class_flg, change_dttm)` +
		`VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING aircraft_model_id`

	createListQuery := fmt.Sprintf(q, aircraftTable)
	row := tx.QueryRow(createListQuery, aircraft.Iata, aircraft.Name, aircraft.Manifacturer, aircraft.Type,
		aircraft.Icaic, aircraft.WingType, aircraft.EconomyClass, aircraft.PrEconomyClass, aircraft.BusinessClass, aircraft.FirstClass, time.Now())
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *AircraftPostgres) GetAll() ([]gvapi.Aircraft, error) {
	var countries []gvapi.Aircraft

	query := fmt.Sprintf(`SELECT aircraft_model_id, aircraft_iata_code, aircraft_model_name, aircraft_model_type`+
		`, aircraft_icaic_code, aircraft_model_wing_type, economy_class_flg, pr_economy_class_flg, business_class_flg, first_class_flg `+
		`, COALESCE(aircraft_model_manufacturer, '') as aircraft_model_manufacturer `+
		` FROM %s	`, aircraftTable)
	err := r.db.Select(&countries, query)

	return countries, err
}

func (r *AircraftPostgres) GetById(aircraftId int) (gvapi.Aircraft, error) {
	var aircraft gvapi.Aircraft

	query := fmt.Sprintf(`SELECT aircraft_model_id, aircraft_iata_code, aircraft_model_name, aircraft_model_type`+
		`, aircraft_icaic_code, aircraft_model_wing_type, economy_class_flg, pr_economy_class_flg, business_class_flg, first_class_flg `+
		`, COALESCE(aircraft_model_manufacturer, '') as aircraft_model_manufacturer `+
		` FROM %s WHERE aircraft_model_id = $1`, aircraftTable)

	if err := r.db.Get(&aircraft, query, aircraftId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return aircraft, errors.New("No aircraft with such id found")
		case nil:
			return aircraft, nil
		default:
			return aircraft, err
		}
	}
	return aircraft, nil
}
