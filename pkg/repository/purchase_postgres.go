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
	q := `INSERT INTO %s (user_id, flight_id, class_flg, food_flg, payed, payed_dttm, change_dttm) ` +
		`VALUES ($1, $2,  $3, $4, 0, $5, $6) RETURNING purchase_id`
	createpurchaseQuery := fmt.Sprintf(q, purchaseTable)
	fmt.Println(createpurchaseQuery)
	t := new(time.Time)
	fmt.Println(t)
	row := tx.QueryRow(createpurchaseQuery, userId, purchase.FlightId, purchase.Class, purchase.Food,
		t, time.Now())

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *PurchasePostgres) GetAll() ([]gvapi.Purchase, error) {
	purchases := []gvapi.Purchase{}

	query := fmt.Sprintf("SELECT *  FROM %s tl", purchaseTable)
	err := r.db.Select(&purchases, query)

	return purchases, err
}

func (r *PurchasePostgres) GetByUserId(userId int) ([]gvapi.Purchase, error) {
	purchases := []gvapi.Purchase{}

	query := fmt.Sprintf("SELECT * FROM %s tl WHERE user_id = $1 AND payed = 1", purchaseTable)

	if err := r.db.Select(&purchases, query, userId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return purchases, errors.New("Не найдено купленных билетов для данного пользователя.")
		case nil:
			return purchases, nil
		default:
			return purchases, err
		}
	}

	return purchases, nil
}

func (r *PurchasePostgres) GetBasketByUserId(userId int) ([]gvapi.Purchase, error) {

	purchases := []gvapi.Purchase{}
	query := fmt.Sprintf(`SELECT *, 
	   CASE 
			 WHEN date_part('minute', current_timestamp - purchase_dttm) > 15 THEN 0
			 ELSE 15 - date_part('minute', current_timestamp - purchase_dttm)
		 END time_left `+
		`FROM %s tl WHERE user_id = $1 AND payed = 0`, purchaseTable)

	if err := r.db.Select(&purchases, query, userId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return purchases, errors.New("Не найдено забронированных билетов")
		case nil:
			return purchases, nil
		default:
			return purchases, err
		}
	}

	for _, v := range purchases {
		fmt.Println(v)
	}
	return purchases, nil
}

func (r *PurchasePostgres) GetById(purchaseId int) (gvapi.Purchase, error) {
	var purchase gvapi.Purchase

	query := fmt.Sprintf(`SELECT *, 
	CASE 
		WHEN date_part('minute', current_timestamp - purchase_dttm) > 15 THEN 0
		ELSE 15 - date_part('minute', current_timestamp - purchase_dttm)
	END time_left `+
		`FROM %s  WHERE purchase_id = $1`, purchaseTable)
	if err := r.db.Get(&purchase, query, purchaseId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return purchase, errors.New("Ошибка: Не найдено билетов с переданными параметрами.")
		case nil:
			return purchase, nil
		default:
			return purchase, err
		}
	}

	return purchase, nil
}

func (r *PurchasePostgres) Update(input gvapi.PurchasePayInput) error {

	query := fmt.Sprintf(`UPDATE %s SET payed = 1, pay_method = $1, payed_dttm = $2, change_dttm = $3 WHERE purchase_id = $4`, purchaseTable)
	_, err := r.db.Exec(query, input.PayMethod, time.Now(), time.Now(), input.PurchaseId)
	return err
}
