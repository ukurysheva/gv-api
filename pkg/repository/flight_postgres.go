package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	gvapi "github.com/ukurysheva/gv-api"
)

type FlightPostgres struct {
	db *sqlx.DB
}

func NewFlightPostgres(db *sqlx.DB) *FlightPostgres {
	return &FlightPostgres{db: db}
}

func (r *FlightPostgres) Create(userId int, flight gvapi.Flight) (int, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	q := `INSERT INTO %s (flight_name, airline_id, ticket_num_economy_class, ticket_num_pr_economy_class, ` +
		`ticket_num_business_class, ticket_num_first_class, cost_economy_class_rub, cost_pr_economy_class_rub, cost_business_class_rub, ` +
		`cost_first_class_rub,aircraft_model_id, departure_airport_id, landing_airport_id, departure_time, landing_time, ` +
		`max_luggage_weight_kg, cost_luggage_weight_rub, max_hand_luggage_weight_kg, cost_hand_luggage_weight_rub, wifi_flg, food_flg, ` +
		`usb_flg, change_dttm ) ` +
		`VALUES ($1, $2, $3, $4, $5, $6, $7, $8,  $9, $10, $11, $12, $13, $14, $15, $16, $17,$18, $19, $20, $21, $22, $23) RETURNING flight_id`
	createFlightQuery := fmt.Sprintf(q, flightTable)

	row := tx.QueryRow(createFlightQuery, flight.Name, flight.AirlineId, flight.TicketNumEconomy, flight.TicketNumPrEconomy,
		flight.TicketNumBusiness, flight.TicketNumFirstClass, flight.CostRubEconomy, flight.CostRubPrEconomy, flight.CostRubBusiness,
		flight.CostRubFirstClass, flight.AircraftId, flight.AirportDepId, flight.AirportLandId, flight.DepartureTime, flight.LandingTime,
		flight.MaxLugWeightKg, flight.CostLugWeightRub, flight.MaxHandLugWeightKg, flight.CostHandLugWeightRub, flight.Wifi, flight.Food,
		flight.Usb, time.Now())

	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *FlightPostgres) GetAll() ([]gvapi.Flight, error) {
	var flights []gvapi.Flight

	query := fmt.Sprintf(`SELECT fl.*,`+
		`fl.ticket_num_economy_class -  
		    COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
				GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
				GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
				GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
				GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
	FROM %s fl
											`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable)
	err := r.db.Select(&flights, query)

	return flights, err
}

func (r *FlightPostgres) GetById(flightId int) (gvapi.Flight, error) {
	var flight gvapi.Flight

	query := fmt.Sprintf(`SELECT fl.*,`+
		`fl.ticket_num_economy_class -  
		     COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
					GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
					GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
					GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
					GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
		FROM %s fl WHERE fl.flight_id = $1
												`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable)
	if err := r.db.Get(&flight, query, flightId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return flight, errors.New("Nothing found")
		case nil:
			return flight, nil
		default:
			return flight, err
		}
	}

	return flight, nil
}

func (r *FlightPostgres) GetByParams(input gvapi.FlightSearchParams) ([]gvapi.Flight, error) {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	var flights []gvapi.Flight

	if input.Food != "" {
		setValues = append(setValues, fmt.Sprintf("AND food_flg=$%d", argId))
		fmt.Println(input.Food)
		fmt.Printf("Food: %T\n", input.Food)
		args = append(args, input.Food)
		argId++
	}

	if input.MaxLugWeightKg != 0 {
		setValues = append(setValues, fmt.Sprintf("AND max_luggage_weight_kg >= $%d", argId))
		args = append(args, input.MaxLugWeightKg)
		argId++
	}

	if input.Class != "" {
		switch input.Class {
		case "economy":
			setValues = append(setValues, fmt.Sprintf("AND ticket_num_economy_class_avail > 0"))
		case "pr_economy":
			setValues = append(setValues, fmt.Sprintf("AND ticket_num_pr_economy_class_avail > 0"))
		case "business":
			setValues = append(setValues, fmt.Sprintf("AND ticket_num_business_class_avail > 0"))
		case "first":
			setValues = append(setValues, fmt.Sprintf("AND ticket_num_first_class_avail > 0"))
		}
	}

	setQuery := strings.Join(setValues, " ")

	query := fmt.Sprintf(`SELECT * FROM 
	(
		SELECT fl.*,`+
		`fl.ticket_num_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
					GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
					GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
					GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
					GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
		FROM %s fl
		) q1 
		WHERE TRUE %s
												`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, setQuery)

	err := r.db.Select(&flights, query, args...)

	return flights, err
}
