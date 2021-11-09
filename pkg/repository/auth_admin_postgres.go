package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type AuthAdminPostgres struct {
	db *sqlx.DB
}

func NewAuthAdminPostgres(db *sqlx.DB) *AuthAdminPostgres {
	return &AuthAdminPostgres{db: db}
}

func (auth *AuthAdminPostgres) CreateAdminUser(adminUser gvapi.AdminUser) (int, error) {
	var id int
	// id = 1
	query := fmt.Sprintf("INSERT INTO %s (admin_email, admin_password, privileges_level, change_dttm) values ($1, $2, $3, $4) RETURNING admin_id", usersAdminTable)

	row := auth.db.QueryRow(query, adminUser.Email, adminUser.Password, "1", time.Now())
	if err := row.Scan(&id); err != nil {
		// fmt.Println()
		return 0, err
	}

	return id, nil
}
func (auth *AuthAdminPostgres) GetUserAdmin(email, password string) (gvapi.AdminUser, error) {
	var adminUser gvapi.AdminUser
	// TODO - get user from db
	query := fmt.Sprintf("SELECT admin_id FROM %s WHERE admin_email=$1 AND admin_password=$2", usersAdminTable)
	err := auth.db.Get(&adminUser, query, email, password)
	fmt.Println(adminUser)
	return adminUser, err
}
