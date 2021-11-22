package gvapi

type Country struct {
	Id         int    `json:"countryId" db:"country_id"`
	Code       string `json:"countryCode" db:"country_code"`
	Name       string `json:"countryName" db:"country_name"`
	Continent  string `json:"countryContinent" db:"country_continent"`
	Wiki       string `json:"countryWiki" db:"country_wikipedia_link"`
	ChangeDate string `json:"-" db:"change_dttm"`
}
type Airline struct {
	Id         int    `json:"airlineId" db:"airline_id"`
	Iata       string `json:"airlineIata" db:"airline_iata" binding:"required"`
	Icao       string `json:"airlineIcao" db:"airline_icao" binding:"required"`
	CountryId  int    `json:"airlineCountryId" db:"airline_country_id" binding:"required"`
	Active     string `json:"airlineActive" db:"airline_active_flg" binding:"required"`
	ChangeDate string `json:"-" db:"change_dttm"`
}

type Aircraft struct {
	Id             int    `json:"aircraftId" db:"aircraft_model_id"`
	Iata           string `json:"aircraftIata" db:"aircraft_iata_code" binding:"required"`
	Name           string `json:"aircraftName" db:"aircraft_model_name" binding:"required"`
	Manifacturer   string `json:"aircraftManufacturer" db:"aircraft_model_manufacturer" binding:"required"`
	WingType       string `json:"aircraftWingType" db:"aircraft_model_wing_type" binding:"required"`
	Type           string `json:"aircraftType" db:"aircraft_model_type" binding:"required"`
	Icaic          string `json:"aircraftIcaic" db:"aircraft_icaic_code" binding:"required"`
	EconomyClass   string `json:"economyClass" db:"economy_class_flg" binding:"required"`
	PrEconomyClass string `json:"prEconomyClass" db:"pr_economy_class_flg" binding:"required"`
	BusinessClass  string `json:"businessClass" db:"business_class_flg" binding:"required"`
	FirstClass     string `json:"firstClass" db:"first_class_flg" binding:"required"`
	ChangeDate     string `json:"changeDate" db:"change_dttm"`
}

type Airport struct {
	Id           int    `json:"id" db:"airport_id"`
	Name         string `json:"airportName" db:"airport_name"`
	Type         string `json:"airportType" db:"airport_type"`
	IataCode     string `json:"airportCode" db:"airport_iata_code"`
	CountryId    int    `json:"airportCountryId" db:"airport_iso_country_id"`
	Region       string `json:"airportIsoRegion" db:"airport_iso_region"`
	Municipality string `json:"airportMunicipality" db:"airport_municipality"`
	HomeLink     string `json:"airportHomeLink" db:"airport_home_link"`
	Visa         string `json:"airportVisa" db:"visa_flg"`
	Quarantine   string `json:"airportQuarantine" db:"quarantine_flg"`
	CovidTest    string `json:"airportCovidTest" db:"covid_test_flg"`
	LockDown     string `json:"airportLockDown" db:"lockdown_flg"`
	ChangeDate   string `json:"changeDate" db:"change_dttm"`
}

type User struct {
	Id         int    `json:"-" db:"user_id"`
	Email      string `json:"userEmail" db:"user_email" binding:"required"`
	Password   string `json:"userPassword" db:"user_password" binding:"required"`
	FirstName  string `json:"userFirstName" db:"user_first_name" binding:"required"`
	LastName   string `json:"userLastName" db:"user_last_name" binding:"required"`
	MiddleName string `json:"userMiddleName" db:"user_middle_name"`
	PhoneNum   string `json:"userPhoneNum" db:"user_phone_number" binding:"required"`
	BirthDate  string `json:"birthDate" db:"birth_date" binding:"required"`
	CountryId  int    `json:"userCountryId" db:"user_country_id"`
	CreateDate string `json:"-" db:"create_dttm"`
	ChangeDate string `json:"-" db:"change_dttm"`
}

type Client struct {
	Id              int    `json:"id" db:"client_id"`
	FirstName       string `json:"firstName" db:"client_first_name" binding:"required"`
	LastName        string `json:"lastName" db:"client_last_name" binding:"required"`
	MiddleName      string `json:"middleName" db:"client_middle_name"`
	PhoneNum        string `json:"phoneNum" db:"client_phone_number" binding:"required"`
	BirthDate       string `json:"birthDate" db:"birth_date" binding:"required"`
	PassportNumber  string `json:"passportNumber" db:"client_passport_number" binding:"required"`
	PassportSeries  string `json:"passportSeries" db:"client_passport_series" binding:"required"`
	PassportAddress string `json:"passportAddress" db:"client_passport_address" binding:"required"`
	LivingAddress   string `json:"livingAddress" db:"client_living_address" binding:"required"`
	CountryId       int    `json:"countryId" db:"client_country_id"`
	CreateDate      string `json:"-" db:"create_dttm"`
	ChangeDate      string `json:"-" db:"change_dttm"`
}

type Flight struct {
	Id                   int     `json:"id" db:"flight_id"`
	Name                 string  `json:"flightName" db:"flight_name"`
	AirlineId            int     `json:"airlineId" db:"airline_id"`
	TicketNumEconomy     int     `json:"ticketNumEconomy" db:"ticket_num_economy_class"`
	TicketNumPrEconomy   int     `json:"ticketNumPrEconomy" db:"ticket_num_pr_economy_class"`
	TicketNumBusiness    int     `json:"ticketNumBusiness" db:"ticket_num_business_class"`
	TicketNumFirstClass  int     `json:"ticketNumFirstClass" db:"ticket_num_first_class"`
	CostRubEconomy       float32 `json:"costRubEconomy" db:"cost_economy_class_rub"`
	CostRubPrEconomy     float32 `json:"costRubPrEconomy" db:"cost_pr_economy_class_rub"`
	CostRubBusiness      float32 `json:"costRubBusiness" db:"cost_business_class_rub"`
	CostRubFirstClass    float32 `json:"costRubFirstClass" db:"cost_first_class_rub"`
	AircraftId           int     `json:"aircraftId" db:"aircraft_model_id"`
	AirportDepId         int     `json:"airportDepId" db:"departure_airport_id"`
	AirportLandId        int     `json:"airportLandId" db:"landing_airport_id"`
	DepartureTime        string  `json:"departureTime" db:"departure_time"`
	LandingTime          string  `json:"landingTime" db:"landing_time"`
	MaxLugWeightKg       float32 `json:"maxLuggageWeightKg" db:"max_luggage_weight_kg"`
	CostLugWeightRub     float32 `json:"costLuggageWeightRub" db:"cost_luggage_weight_rub"`
	MaxHandLugWeightKg   float32 `json:"maxHandLuggageWeightKg" db:"max_hand_luggage_weight_kg"`
	CostHandLugWeightRub float32 `json:"costHandLuggageWeightRub" db:"cost_hand_luggage_weight_rub"`
	Wifi                 string  `json:"wifiFlg" db:"wifi_flg"`
	Food                 string  `json:"foodFlg" db:"food_flg"`
	Usb                  string  `json:"usbFlg" db:"usb_flg"`
	ChangeDate           string  `json:"-" db:"change_dttm"`
}
