package v1

type BlockedMSISDN struct {
	IP          string  `json:"ip"`
	CountryName string  `json:"country_name"`
	CountryCode string  `json:"country_code"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int   `json:"metro_code"`
}
