package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type CountryPostgres struct {
	db *sqlx.DB
}

func NewCountryPostgres(db *sqlx.DB) *CountryPostgres {
	return &CountryPostgres{db: db}
}

func (r *CountryPostgres) Create(userId int, country gvapi.Country) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (country_code, country_name, country_continent, country_wikipedia_link, change_dttm) ` +
		`VALUES ($1, $2, $3, $4, $5) RETURNING country_id`
	createListQuery := fmt.Sprintf(q, countryTable)
	row := tx.QueryRow(createListQuery, country.Code, country.Name, country.Continent, country.Wiki, time.Now())
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *CountryPostgres) GetAll() ([]gvapi.Country, error) {
	var countries []gvapi.Country

	query := fmt.Sprintf("SELECT * FROM %s tl", countryTable)
	err := r.db.Select(&countries, query)

	return countries, err
}

func (r *CountryPostgres) GetById(countryId int) (gvapi.Country, error) {
	var country gvapi.Country

	query := fmt.Sprintf(`SELECT country_id, country_code, country_name, country_wikipedia_link, country_continent FROM %s 
	                      WHERE country_id = $1`, countryTable)

	if err := r.db.Get(&country, query, countryId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return country, errors.New("No country with such id found")
		case nil:
			return country, nil
		default:
			return country, err
		}
	}
	return country, nil
}
