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
		`cd.country_name AS departure_country_name, cl.country_name AS landing_country_name, `+
		`acr.aircraft_model_name AS aircraft_name, `+
		`airline.airline_name AS airline_name, `+
		`fl.ticket_num_economy_class -  
		    COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
				GROUP BY pr.flight_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
				GROUP BY pr.flight_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
				GROUP BY pr.flight_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
				COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
				WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
				GROUP BY pr.flight_id), 0) AS ticket_num_first_class_avail
	FROM %s fl 
	LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
	LEFT JOIN %s cd ON apd.airport_iso_country_id = cd.country_id
	LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
	LEFT JOIN %s cl ON apl.airport_iso_country_id = cl.country_id
	LEFT JOIN %s acr ON fl.aircraft_model_id = acr.aircraft_model_id
	LEFT JOIN %s airline ON fl.airline_id = airline.airline_id ORDER BY fl.flight_id
	`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, countryTable, airportTable, countryTable, aircraftTable, airlineTable)
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
		`apd.airport_name AS departure_airport_name, apl.airport_name AS landing_airport_name, `+
		`cd.country_name AS departure_country_name, cl.country_name AS landing_country_name, `+
		`acr.aircraft_model_name AS aircraft_name, `+
		`airline.airline_name AS airline_name, `+
		`fl.ticket_num_economy_class -  
		     COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy' 
					GROUP BY pr.flight_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
	      	COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy' 
					GROUP BY pr.flight_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business' 
					GROUP BY pr.flight_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
					COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
					WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
					GROUP BY pr.flight_id), 0) AS ticket_num_first_class_avail
		FROM %s fl 
		LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
  	LEFT JOIN %s cd ON apd.airport_iso_country_id = cd.country_id
  	LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
	  LEFT JOIN %s cl ON apl.airport_iso_country_id = cl.country_id
	  LEFT JOIN %s acr ON fl.aircraft_model_id = acr.aircraft_model_id
	  LEFT JOIN %s airline ON fl.airline_id = airline.airline_id
		WHERE fl.flight_id = $1
		`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, countryTable, airportTable, countryTable, aircraftTable, airlineTable)

	if err := r.db.Get(&flight, query, flightId); err != nil {
		switch err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
			return flight, errors.New("Не удалось совершить бронирование данного рейса. Пожалуйста, повторите попытку.")
		case nil:
			return flight, nil
		default:
			return flight, err
		}
	}

	return flight, nil
}

func (r *FlightPostgres) GetByParams(input gvapi.FlightSearchParams) ([]gvapi.Flight, []gvapi.Flight, error) {

	setValuesFlight := make([]string, 0)
	setValuesExtTo := make([]string, 0)
	setValuesExtBack := make([]string, 0)
	argsCountryTo := make([]interface{}, 0)
	argsCountryBack := make([]interface{}, 0)
	argIdCountryTo := 1
	argIdCountryFrom := 1
	flightsTo := []gvapi.Flight{}
	flightsBack := []gvapi.Flight{}
	var err error

	if input.Food != "" && input.Food == "Y" {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND food_flg=$%d", argIdCountryTo))
		fmt.Println(input.Food)
		fmt.Printf("Food: %T\n", input.Food)
		argsCountryTo = append(argsCountryTo, input.Food)
		argsCountryBack = append(argsCountryBack, input.Food)
		argIdCountryTo++
		argIdCountryFrom++
	}

	if input.MaxLugWeightKg != 0 {
		setValuesFlight = append(setValuesFlight, fmt.Sprintf("AND max_luggage_weight_kg >= $%d", argIdCountryTo))
		argsCountryTo = append(argsCountryTo, input.MaxLugWeightKg)
		argsCountryBack = append(argsCountryBack, input.MaxLugWeightKg)
		argIdCountryTo++
		argIdCountryFrom++
	}

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

	if input.DateFrom != "" {

		setValuesExtTo = append(setValuesExtTo, fmt.Sprintf("AND DATE(fl.departure_time) = DATE($%d)", argIdCountryTo))
		argsCountryTo = append(argsCountryTo, input.DateFrom)
		argIdCountryTo++
	}
	if input.CountryIdFrom != 0 {
		setValuesExtTo = append(setValuesExtTo, fmt.Sprintf("AND apd.airport_iso_country_id = $%d", argIdCountryTo))
		argsCountryTo = append(argsCountryTo, input.CountryIdFrom)
		argIdCountryTo++
	}
	if input.CountryIdTo != 0 {
		setValuesExtTo = append(setValuesExtTo, fmt.Sprintf("AND apl.airport_iso_country_id = $%d", argIdCountryTo))
		argsCountryTo = append(argsCountryTo, input.CountryIdTo)
		argIdCountryTo++
	}

	setQueryFlight := strings.Join(setValuesFlight, " ")
	setQueryExt := strings.Join(setValuesExtTo, " ")

	query := fmt.Sprintf(`SELECT * FROM 
		(
			SELECT fl.flight_id, fl.flight_name, fl.airline_id, fl.ticket_num_economy_class, fl.ticket_num_pr_economy_class, `+
		`fl.ticket_num_business_class, fl.ticket_num_first_class, fl.cost_economy_class_rub, fl.cost_pr_economy_class_rub, fl.cost_business_class_rub, `+
		`fl.cost_first_class_rub,fl.aircraft_model_id, fl.departure_airport_id, fl.landing_airport_id, fl.departure_time, fl.landing_time, `+
		`fl.max_luggage_weight_kg, fl.cost_luggage_weight_rub, fl.max_hand_luggage_weight_kg, fl.cost_hand_luggage_weight_rub, fl.wifi_flg, fl.food_flg, `+
		`fl.usb_flg, fl.change_dttm , apd.airport_iso_country_id AS departure_country_id, apl.airport_iso_country_id AS landing_country_id,`+
		`apd.airport_name AS departure_airport_name, apl.airport_name AS landing_airport_name, `+
		`cd.country_name AS departure_country_name, cl.country_name AS landing_country_name, `+
		`fl.ticket_num_economy_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
						GROUP BY pr.flight_id), 0) AS ticket_num_economy_class_avail,`+
		`fl.ticket_num_pr_economy_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
						GROUP BY pr.flight_id), 0) AS ticket_num_pr_economy_class_avail,`+
		`fl.ticket_num_business_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
						GROUP BY pr.flight_id), 0) 	AS ticket_num_business_class_avail,`+
		`fl.ticket_num_first_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
						GROUP BY pr.flight_id), 0) AS ticket_num_first_class_avail
			FROM %s fl
			LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
			LEFT JOIN %s cd ON apd.airport_iso_country_id = cd.country_id
			LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
			LEFT JOIN %s cl ON apl.airport_iso_country_id = cl.country_id
			WHERE TRUE %s
			ORDER BY fl.flight_id
			) q1 
			
			
			WHERE TRUE %s
			
													`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, countryTable, airportTable, countryTable, setQueryExt, setQueryFlight)
	err = r.db.Select(&flightsTo, query, argsCountryTo...)

	if input.BothWays == "Y" {
		if input.DateTo != "" {
			setValuesExtBack = append(setValuesExtBack, fmt.Sprintf("AND DATE(fl.departure_time) = DATE($%d)", argIdCountryFrom))
			argsCountryBack = append(argsCountryBack, input.DateTo)
			argIdCountryFrom++
		}
		if input.CountryIdTo != 0 {
			setValuesExtBack = append(setValuesExtBack, fmt.Sprintf("AND apd.airport_iso_country_id = $%d", argIdCountryFrom))
			argsCountryBack = append(argsCountryBack, input.CountryIdTo)
			argIdCountryFrom++
		}
		if input.CountryIdFrom != 0 {
			setValuesExtBack = append(setValuesExtBack, fmt.Sprintf("AND apl.airport_iso_country_id = $%d", argIdCountryFrom))
			argsCountryBack = append(argsCountryBack, input.CountryIdFrom)
			argIdCountryFrom++
		}

		setQueryFlight := strings.Join(setValuesFlight, " ")
		setQueryExt := strings.Join(setValuesExtBack, " ")

		query := fmt.Sprintf(`SELECT * FROM 
		(
			SELECT fl.flight_id, fl.flight_name, fl.airline_id, fl.ticket_num_economy_class, fl.ticket_num_pr_economy_class, `+
			`fl.ticket_num_business_class, fl.ticket_num_first_class, fl.cost_economy_class_rub, fl.cost_pr_economy_class_rub, fl.cost_business_class_rub, `+
			`fl.cost_first_class_rub,fl.aircraft_model_id, fl.departure_airport_id, fl.landing_airport_id, fl.departure_time, fl.landing_time, `+
			`fl.max_luggage_weight_kg, fl.cost_luggage_weight_rub, fl.max_hand_luggage_weight_kg, fl.cost_hand_luggage_weight_rub, fl.wifi_flg, fl.food_flg, `+
			`fl.usb_flg, fl.change_dttm , apd.airport_iso_country_id AS departure_country_id, apl.airport_iso_country_id AS landing_country_id,`+
			`apd.airport_name AS departure_airport_name, apl.airport_name AS landing_airport_name, `+
			`cd.country_name AS departure_country_name, cl.country_name AS landing_country_name, `+
			`fl.ticket_num_economy_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'economy'
						GROUP BY pr.flight_id), 0) AS ticket_num_economy_class_avail,`+
			`fl.ticket_num_pr_economy_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'pr_economy'
						GROUP BY pr.flight_id), 0) AS ticket_num_pr_economy_class_avail,`+
			`fl.ticket_num_business_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'business'
						GROUP BY pr.flight_id), 0) 	AS ticket_num_business_class_avail,`+
			`fl.ticket_num_first_class -  
						COALESCE((SELECT COUNT(pr.purchase_id) FROM %s pr
						WHERE pr.flight_id = fl.flight_id AND pr.class_flg = 'first_class'
						GROUP BY pr.flight_id), 0) AS ticket_num_first_class_avail
			FROM %s fl
			LEFT JOIN %s apd ON fl.departure_airport_id = apd.airport_id
			LEFT JOIN %s cd ON apd.airport_iso_country_id = cd.country_id
			LEFT JOIN %s apl ON fl.landing_airport_id = apl.airport_id
			LEFT JOIN %s cl ON apl.airport_iso_country_id = cl.country_id
			WHERE TRUE %s ORDER BY fl.flight_id
			) q1 
			
			
			WHERE TRUE %s
			
													`, purchaseTable, purchaseTable, purchaseTable, purchaseTable, flightTable, airportTable, countryTable, airportTable, countryTable, setQueryExt, setQueryFlight)
		err = r.db.Select(&flightsBack, query, argsCountryBack...)
	}

	return flightsTo, flightsBack, err
}
