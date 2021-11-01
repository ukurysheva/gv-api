package gvapi

type Country struct {
	Id         int    `json:"id" db:"country_id"`
	Code       string `json:"countryCode" db:"country_code"`
	Name       string `json:"countryName" db:"country_name"`
	Continent  string `json:"countryContinent" db:"country_continent"`
	Wiki       string `json:"countryWiki" db:"country_wikipedia_link"`
	ChangeDate string `json:"changeDate" db:"change_dttm"`
}
type Airline struct {
	Id         int    `json:"id" db:"airline_id"`
	Iata       string `json:"airlineIata" db:"airline_iata"`
	Icao       string `json:"airlineIcao" db:"airline_icao"`
	CountryId  string `json:"airlineCountryId" db:"airline_country_id"`
	Active     int    `json:"airlineActive" db:"airline_active_flg"`
	ChangeDate string `json:"changeDate" db:"change_dttm"`
}

type Aircraft struct {
	Id            int    `json:"id" db:"aircraft_id"`
	Iata          string `json:"aircraftIata" db:"aircraft_model_iata_code"`
	Name          string `json:"aircraftName" db:"aircraft_model_name"`
	Manifacturer  string `json:"aircraftManifacturer" db:"aircraft_model_manifacturer"`
	Type          string `json:"aircraftType" db:"aircraft_model_type"`
	Icaic         string `json:"aircraftIcaic" db:"aircraft_model_icaic_code"`
	EconomyClass  int    `json:"economyClass" db:"pr_economy_class_flg"`
	BusinessClass int    `json:"businessClass" db:"business_class_flg"`
	FirstClass    int    `json:"firstClass" db:"first_class_flg"`
	ChangeDate    string `json:"changeDate" db:"change_dttm"`
}

type Airport struct {
	Id        int    `json:"id" db:"airport_id"`
	Name      string `json:"airportName" db:"airport_name"`
	Type      string `json:"airportType" db:"airport_type"`
	Code      string `json:"airportCode" db:"airport_code"`
	CountryId string `json:"airportCountryId" db:"airport_iso_country_id"`
}
