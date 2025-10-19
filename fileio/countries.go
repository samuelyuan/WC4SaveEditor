package fileio

// CountryInfo represents a country with its ID, name, and color
type CountryInfo struct {
	ID    uint8
	Name  string
	Color string
}

// Country codes for World Conqueror 4
// Format: (country_id_hex, name, color_hex)
var countryCodes = []CountryInfo{
	{0x01, "UK", "fafa96"},
	{0x02, "France", "96e6ff"},
	{0x03, "Germany", "a0a0a0"},
	{0x04, "Germany", "a0a0a0"},
	{0x05, "Soviet Union", "f0aaaa"},
	{0x06, "USA", "b4ffdc"},
	{0x07, "Italy", "dcbee6"},
	{0x08, "ROC", "64c8fa"},
	{0x09, "PRC", "ff8787"},
	{0x0a, "Japan", "e6c85a"},
	{0x0b, "Finland", "bedcc8"},
	{0x0c, "Poland", "fadcc8"},
	{0x0d, "Yugoslavia", "d2bea0"},
	{0x0e, "Canada", "96b48c"},
	{0x0f, "Australia", "8296ff"},
	{0x10, "Norway", "e1a0e1"},
	{0x11, "Sweden", "87a5d2"},
	{0x12, "Denmark", "9682b4"},
	{0x13, "Netherlands", "e6aa8c"},
	{0x14, "Belgium", "a0c896"},
	{0x15, "Spain", "e696dc"},
	{0x16, "Portugal", "fadc96"},
	{0x17, "Hungary", "facd8c"},
	{0x18, "Romania", "ffeba5"},
	{0x19, "Bulgaria", "649678"},
	{0x1a, "Switzerland", "9aa588"},
	{0x1b, "Greece", "cdffff"},
	{0x1c, "Turkey", "a59678"},
	{0x1d, "Saudi Arabia", "009a00"},
	{0x1e, "Iraq", "b49696"},
	{0x1f, "Iran", "78b4c8"},
	{0x20, "India", "82b4fa"},
	{0x21, "Thailand", "b4e664"},
	{0x22, "Mongolia", "e19664"},
	{0x23, "North Korea", "ffa082"},
	{0x24, "South Korea", "f5ebd7"},
	{0x25, "Mexico", "b48750"},
	{0x26, "Cuba", "dca5be"},
	{0x27, "Colombia", "64C8C8"},
	{0x28, "Brazil", "46be64"},
	{0x29, "Bolivia", "64AA50"},
	{0x2a, "Venezuela", "82966E"},
	{0x2b, "Peru", "785A64"},
	{0x2c, "Chile", "648CC8"},
	{0x2d, "Argentina", "7ECEF4"},
	{0x2e, "Egypt", "d2b45a"},
	{0x2f, "Liberia", "1d2739"},
	{0x30, "Mysterious Forces", "7878b4"},
}

// Maps for faster lookups
var countryIDToName map[uint8]string
var countryIDToColor map[uint8]string

// Initialize the lookup maps
func init() {
	countryIDToName = make(map[uint8]string)
	countryIDToColor = make(map[uint8]string)
	for _, country := range countryCodes {
		countryIDToName[country.ID] = country.Name
		countryIDToColor[country.ID] = country.Color
	}
}

// GetCountryName returns the country name by country ID
func GetCountryName(countryID uint8) (string, bool) {
	name, exists := countryIDToName[countryID]
	return name, exists
}

// GetCountryColor returns the country color by country ID
func GetCountryColor(countryID uint8) (string, bool) {
	color, exists := countryIDToColor[countryID]
	return color, exists
}

// GetAllCountries returns all country codes
func GetAllCountries() []CountryInfo {
	return countryCodes
}

// GetCountryCount returns the total number of countries
func GetCountryCount() int {
	return len(countryCodes)
}

// GetCountryInfo returns both name and color for a country ID
func GetCountryInfo(countryID uint32) (string, string) {
	name, _ := GetCountryName(uint8(countryID))
	color, _ := GetCountryColor(uint8(countryID))
	return name, color
}
