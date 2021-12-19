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

	query := fmt.Sprintf(`SELECT fl.flight_id, fl.flight_name, fl.airline_id, fl.ticket_num_economy_class, fl.ticket_num_pr_economy_class, `+
		`fl.ticket_num_business_class, fl.ticket_num_first_class, fl.cost_economy_class_rub, fl.cost_pr_economy_class_rub, fl.cost_business_class_rub, `+
		`fl.cost_first_class_rub,fl.aircraft_model_id, fl.departure_airport_id, fl.landing_airport_id, fl.departure_time, fl.landing_time, `+
		`fl.max_luggage_weight_kg, fl.cost_luggage_weight_rub, fl.max_hand_luggage_weight_kg, fl.cost_hand_luggage_weight_rub, fl.wifi_flg, fl.food_flg, `+
		`fl.usb_flg, fl.change_dttm , apd.airport_iso_country_id AS departure_country_id, apl.airport_iso_country_id AS landing_country_id, `+
		`apd.airport_name AS departure_airport_name, apl.airport_name AS landing_airport_name, `+
		`fl.ticket_num_economy_class -  
		    COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy' AND pr.payed = 1
				GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy' AND pr.payed = 1
				GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business' AND pr.payed = 1
				GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class' AND pr.payed = 1
				GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
	FROM %s fl 
	LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
	LEFT JOIN %s cd ON apd.airport_iso_country_id = cd.country_id
	LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
	LEFT JOIN %s cl ON apl.airport_iso_country_id = cl.country_id
	`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, countryTable, airportTable, countryTable)
	err := r.db.Select(&flights, query)

	return flights, err
}

func (r *FlightPostgres) GetById(flightId int) (gvapi.Flight, error) {
	var flight gvapi.Flight

	query := fmt.Sprintf(`SELECT fl.flight_id, fl.flight_name, fl.airline_id, fl.ticket_num_economy_class, fl.ticket_num_pr_economy_class, `+
		`fl.ticket_num_business_class, fl.ticket_num_first_class, fl.cost_economy_class_rub, fl.cost_pr_economy_class_rub, fl.cost_business_class_rub, `+
		`fl.cost_first_class_rub,fl.aircraft_model_id, fl.departure_airport_id, fl.landing_airport_id, fl.departure_time, fl.landing_time, `+
		`fl.max_luggage_weight_kg, fl.cost_luggage_weight_rub, fl.max_hand_luggage_weight_kg, fl.cost_hand_luggage_weight_rub, fl.wifi_flg, fl.food_flg, `+
		`fl.usb_flg, fl.change_dttm , apd.airport_iso_country_id AS departure_country_id, apl.airport_iso_country_id AS landing_country_id,`+
		`fl.ticket_num_economy_class -  
		     COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
		FROM %s fl 
		LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
  	LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
		WHERE fl.flight_id = $1
												`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, airportTable)
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

	setValuesFlight := make([]string, 0)
	setValuesExt := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	var flights []gvapi.Flight

	if input.Food != "" {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND food_flg=$%d", argId))
		fmt.Println(input.Food)
		fmt.Printf("Food: %T\n", input.Food)
		args = append(args, input.Food)
		argId++
	}

	if input.MaxLugWeightKg != 0 {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND max_luggage_weight_kg >= $%d", argId))
		args = append(args, input.MaxLugWeightKg)
		argId++
	}
	if input.DateFrom != "" {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND DATE(departure_time) = DATE($%d)", argId))
		args = append(args, input.DateFrom)
		argId++
	}
	if input.DateTo != "" {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND DATE(landing_time) = DATE($%d)", argId))
		args = append(args, input.DateTo)
		argId++
	}

	if input.CountryIdFrom != 0 {
		setValuesExt = append(setValuesExt, fmt.Sprintf("AND apd.airport_iso_country_id = $%d", argId))
		args = append(args, input.CountryIdFrom)
		argId++
	}
	if input.CountryIdTo != 0 {
		setValuesExt = append(setValuesExt, fmt.Sprintf("AND apl.airport_iso_country_id = $%d", argId))
		args = append(args, input.CountryIdTo)
		argId++
	}
	// if input.BothWays == "Y" {
	// 	setValuesExt = append(setValuesExt, fmt.Sprintf("AND apl.airport_iso_country_id = $%d", argId))
	// 	args = append(args, input.CountryIdTo)
	// 	argId++
	// }

	if input.Class != "" {
		switch input.Class {
		case "economy":
			setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND ticket_num_economy_class_avail > 0"))
		case "pr_economy":
			setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND ticket_num_pr_economy_class_avail > 0"))
		case "business":
			setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND ticket_num_business_class_avail > 0"))
		case "first":
			setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND ticket_num_first_class_avail > 0"))
		}
	}

	setQueryFlight := strings.Join(setValuesFlight, " ")
	setQueryExt := strings.Join(setValuesExt, " ")

	query := fmt.Sprintf(`SELECT * FROM 
	(
		SELECT fl.flight_id, fl.flight_name, fl.airline_id, fl.ticket_num_economy_class, fl.ticket_num_pr_economy_class, `+
		`fl.ticket_num_business_class, fl.ticket_num_first_class, fl.cost_economy_class_rub, fl.cost_pr_economy_class_rub, fl.cost_business_class_rub, `+
		`fl.cost_first_class_rub,fl.aircraft_model_id, fl.departure_airport_id, fl.landing_airport_id, fl.departure_time, fl.landing_time, `+
		`fl.max_luggage_weight_kg, fl.cost_luggage_weight_rub, fl.max_hand_luggage_weight_kg, fl.cost_hand_luggage_weight_rub, fl.wifi_flg, fl.food_flg, `+
		`fl.usb_flg, fl.change_dttm , apd.airport_iso_country_id AS departure_country_id, apl.airport_iso_country_id AS landing_country_id,`+
		`fl.ticket_num_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class' AND pr.payed = 1
					GROUP BY pr.purchase_id), 0) AS ticket_num_first_class_avail
		FROM %s fl
		LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
		LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
		WHERE TRUE %s
		) q1 
		
		
		WHERE TRUE %s
		
												`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, airportTable, setQueryExt, setQueryFlight)
	err := r.db.Select(&flights, query, args...)

	return flights, err
}
