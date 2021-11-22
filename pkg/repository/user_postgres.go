package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) CreateUser(user gvapi.User) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	query := fmt.Sprintf("INSERT INTO %s (user_email, user_password, user_first_name, user_last_name, user_phone_number, birth_date)"+
		" VALUES ($1, $2, $3, $4, $5, $6) RETURNING user_id", usersTable)

	row := r.db.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.PhoneNum, user.BirthDate)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}
func (r *UserPostgres) GetUser(email, password string) (gvapi.User, error) {
	var user gvapi.User
	// TODO - get user from db
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE user_email=$1 AND user_password=$2", usersTable)

	if err := r.db.Get(&user, query, email, password); err != nil {
		switch err {
		case sql.ErrNoRows:
			return user, errors.New("Wrong email and password")
		case nil:
			return user, nil
		default:
			return user, err
		}
	}

	return user, nil
}

func (r *UserPostgres) GetProfile(userId int) (gvapi.User, error) {

	var user gvapi.User
	query := fmt.Sprintf("SELECT user_email, user_first_name, user_last_name, user_phone_number, birth_date, "+
		"COALESCE(user_middle_name, '') as user_middle_name, COALESCE(user_country_id, 0) as user_country_id FROM %s WHERE user_id=$1", usersTable)
	fmt.Println("GetProfile postgres")

	if err := r.db.Get(&user, query, userId); err != nil {
		switch err {
		case sql.ErrNoRows:
			return user, errors.New("Wrong id")
		case nil:
			return user, nil
		default:
			return user, err
		}
	}

	// query := fmt.Sprintf("SELECT * FROM %s tl", countryTable)
	// err := r.db.Select(&countries, query)
	fmt.Println(user)
	return user, nil
}

func (r *UserPostgres) Update(userId int, input gvapi.UpdateUserInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("user_email=$%d", argId))
		args = append(args, *input.Email)
		argId++
	}

	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("user_password=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	if input.FirstName != nil {
		setValues = append(setValues, fmt.Sprintf("user_first_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}

	if input.LastName != nil {
		setValues = append(setValues, fmt.Sprintf("user_last_name=$%d", argId))
		args = append(args, *input.FirstName)
		argId++
	}

	if input.MiddleName != nil {
		setValues = append(setValues, fmt.Sprintf("user_middle_name=$%d", argId))
		args = append(args, *input.MiddleName)
		argId++
	}

	if input.PhoneNum != nil {
		setValues = append(setValues, fmt.Sprintf("user_phone_number=$%d", argId))
		args = append(args, *input.PhoneNum)
		argId++
	}

	if input.BirthDate != nil {
		setValues = append(setValues, fmt.Sprintf("birth_date=$%d", argId))
		args = append(args, *input.BirthDate)
		argId++
	}

	if input.CountryId != nil {
		setValues = append(setValues, fmt.Sprintf("user_country_id=$%d", argId))
		args = append(args, *input.CountryId)
		argId++
	}
	fmt.Println(input)
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE %s SET %s WHERE user_id = $%d`, usersTable, setQuery, argId)
	args = append(args, userId)
	fmt.Println(query)
	_, err := r.db.Exec(query, args...)
	return err
}
