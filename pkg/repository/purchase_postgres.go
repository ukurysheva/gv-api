package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type PurchasePostgres struct {
	db *sqlx.DB
}

func NewPurchasePostgres(db *sqlx.DB) *PurchasePostgres {
	return &PurchasePostgres{db: db}
}

func (r *PurchasePostgres) Create(userId int, purchase gvapi.Purchase) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (user_id, flight_id, cost_rub_amt, class_flg, food_flg, change_dttm ) ` +
		`VALUES ($1, $2, $3, $4, $5, $6) RETURNING purchase_id`
	createpurchaseQuery := fmt.Sprintf(q, purchaseTable)

	row := tx.QueryRow(createpurchaseQuery, userId, purchase.FlightId, purchase.CostRub, purchase.Class, purchase.Food, time.Now())

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *PurchasePostgres) GetAll() ([]gvapi.Purchase, error) {
	var purchases []gvapi.Purchase

	query := fmt.Sprintf("SELECT * FROM %s tl", purchaseTable)
	err := r.db.Select(&purchases, query)

	return purchases, err
}

func (r *PurchasePostgres) GetById(purchaseId int) (gvapi.Purchase, error) {
	var purchase gvapi.Purchase

	query := fmt.Sprintf(`SELECT * FROM %s 
	                      WHERE purchase_id = $1`, purchaseTable)
	if err := r.db.Get(&purchase, query, purchaseId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return purchase, errors.New("Nothing found")
		case nil:
			return purchase, nil
		default:
			return purchase, err
		}
	}

	return purchase, nil
}
