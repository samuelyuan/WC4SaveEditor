package fileio

// CityInfo represents a city with its ID and name
type CityInfo struct {
	ID   uint16
	Name string
}

// CountryInfo represents a country with its ID, name, and color
type CountryInfo struct {
	ID    uint8
	Name  string
	Color string
}

// City codes for World Conqueror 4
// Format: (city_id_hex, name)
var cityCodes = []CityInfo{
	{0x01, "London"},
	{0x03, "Manchester"},
	{0x04, "Dublin"},
	{0x05, "Plymouth"},
	{0x07, "Glasgow"},
	{0x08, "Galway"},
	{0x09, "Stockholm"},
	{0x0b, "Oslo"},
	{0x0d, "Paris"},
	{0x0e, "Bordeaux"},
	{0x0f, "Marseille"},
	{0x10, "Lyon"},
	{0x11, "Brest"},
	{0x15, "Toulouse"},
	{0x17, "Ajaccio"},
	{0x18, "Madrid"},
	{0x19, "Barcelona"},
	{0x1a, "Seville"},
	{0x1c, "A Coruña"},
	{0x1f, "Lisbon"},
	{0x20, "Porto"},
	{0x21, "Brussels"},
	{0x23, "Amsterdam"},
	{0x24, "Berlin"},
	{0x25, "Hamburg"},
	{0x26, "Cologne"},
	{0x27, "Munich"},
	{0x2b, "Zurich"},
	{0x2c, "Rome"},
	{0x2d, "Milan"},
	{0x2e, "Naples"},
	{0x31, "Palermo"},
	{0x32, "Vienna"},
	{0x33, "Warsaw"},
	{0x34, "Poznan"},
	{0x35, "Bialystok"},
	{0x37, "Krakow"},
	{0x38, "Copenhagen"},
	{0x39, "Budapest"},
	{0x3a, "Belgrade"},
	{0x3b, "Sarajevo"},
	{0x3c, "Sofia"},
	{0x3d, "Athens"},
	{0x3e, "Prague"},
	{0x3f, "Bucharest"},
	{0x40, "Cluj-Napoca"},
	{0x41, "Königsberg"},
	{0x42, "Riga"},
	{0x43, "Helsinki"},
	{0x46, "Ankara"},
	{0x47, "Istanbul"},
	{0x4b, "Jerusalem"},
	{0x4c, "Damascus"},
	{0x4d, "Baghdad"},
	{0x4f, "Riyadh"},
	{0x50, "Medina"},
	{0x52, "Dubai"},
	{0x53, "Sanaa"},
	{0x54, "Salalah"},
	{0x55, "Tehran"},
	{0x57, "Mashhad"},
	{0x58, "Almaty"},
	{0x59, "Astana"},
	{0x5c, "Moscow"},
	{0x5d, "Yekaterinburg"},
	{0x5e, "Novosibirsk"},
	{0x5f, "Leningrad"},
	{0x63, "Saratov"},
	{0x69, "Krasnoyarsk"},
	{0x6a, "Vladivostok"},
	{0x6b, "Chita"},
	{0x6c, "Khabarovsk"},
	{0x6d, "Petropavlovsk"},
	{0x70, "Yakutsk"},
	{0x71, "Minsk"},
	{0x73, "Kiev"},
	{0x75, "Donetsk"},
	{0x77, "Stalingrad"},
	{0x7b, "Uliastai"},
	{0x7d, "Ulaanbaatar"},
	{0x7e, "Beijing"},
	{0x7f, "Nanjing"},
	{0x80, "Shanghai"},
	{0x81, "Wuhan"},
	{0x82, "Chongqing"},
	{0x83, "Hong Kong"},
	{0x84, "Guihui"},
	{0x85, "Lanzhou"},
	{0x88, "Urumqi"},
	{0x89, "Lhasa"},
	{0x8b, "Taipei"},
	{0x8d, "Changchun"},
	{0x8e, "Shenyang"},
	{0x90, "New Delhi"},
	{0x91, "Kolkata"},
	{0x92, "Mumbai"},
	{0x93, "Lucknow"},
	{0x97, "Chennai"},
	{0x98, "Mandalay"},
	{0x99, "Yangon"},
	{0x9a, "Bangkok"},
	{0x9c, "Singapore"},
	{0x9d, "Kuala Lumpur"},
	{0x9e, "Hanoi"},
	{0x9f, "Phnom Penh"},
	{0xa0, "Tokyo"},
	{0xa1, "Osaka"},
	{0xa2, "Hiroshima"},
	{0xa3, "Sendai"},
	{0xa4, "Sapporo"},
	{0xa5, "Nagasaki"},
	{0xa6, "Seoul"},
	{0xa7, "Pyongyang"},
	{0xa8, "Busan"},
	{0xa9, "Manila"},
	{0xaa, "Jakarta"},
	{0xab, "Port Moresby"},
	{0xac, "Balikpapan"},
	{0xad, "Medan"},
	{0xae, "Kupang"},
	{0xaf, "Darwin"},
	{0xb0, "Brisbane"},
	{0xb2, "Port Hedland"},
	{0xb4, "Alice Springs"},
	{0xb6, "Casablanca"},
	{0xb7, "Algiers"},
	{0xb8, "Adrar"},
	{0xba, "Tripoli"},
	{0xbc, "Al Jawf"},
	{0xbd, "Cairo"},
	{0xbe, "Alexandria"},
	{0xc0, "Addis Ababa"},
	{0xc1, "Mogadishu"},
	{0xc2, "Khartoum"},
	{0xc3, "Port Sudan"},
	{0xc4, "Juba"},
	{0xc6, "Yaoundé"},
	{0xc7, "N'Djamena"},
	{0xc8, "Bangui"},
	{0xc9, "Dakar"},
	{0xca, "Bamako"},
	{0xcb, "Monrovia"},
	{0xcd, "Niamey"},
	{0xce, "Ottawa"},
	{0xd0, "Toronto"},
	{0xd1, "Edmonton"},
	{0xd3, "Whitehorse"},
	{0xd5, "Vancouver"},
	{0xd6, "Yellowknife"},
	{0xd7, "Winnipeg"},
	{0xdd, "Kuujjuaq"},
	{0xdf, "Washington"},
	{0xe0, "New York"},
	{0xe1, "Miami"},
	{0xe2, "Houston"},
	{0xe3, "Denver"},
	{0xe4, "San Francisco"},
	{0xe5, "Los Angeles"},
	{0xe9, "Dallas"},
	{0xea, "Chicago"},
	{0xee, "Honolulu"},
	{0xf0, "Mexico City"},
	{0xf1, "Monterrey"},
	{0xf3, "Mérida"},
	{0xf4, "Havana"},
	{0xf5, "Santiago"},
	{0xf6, "Port-au-Prince"},
	{0xf7, "Bogotá"},
	{0xf9, "Caracas"},
	{0xfa, "Paramaribo"},
	{0xfb, "Lima"},
	{0xfe, "La Paz"},
	{0xff, "Rio de Janeiro"},
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
	{0x1a, "Switzerland", "ff6984"},
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
var cityIDToName map[uint16]string
var cityNameToID map[string]uint16
var countryIDToName map[uint8]string
var countryIDToColor map[uint8]string

// Initialize the lookup maps
func init() {
	cityIDToName = make(map[uint16]string)
	cityNameToID = make(map[string]uint16)

	for _, city := range cityCodes {
		cityIDToName[city.ID] = city.Name
		cityNameToID[city.Name] = city.ID
	}

	countryIDToName = make(map[uint8]string)
	countryIDToColor = make(map[uint8]string)
	for _, country := range countryCodes {
		countryIDToName[country.ID] = country.Name
		countryIDToColor[country.ID] = country.Color
	}
}

// GetCityName returns the English city name by city ID
func GetCityName(cityID uint16) (string, bool) {
	name, exists := cityIDToName[cityID]
	return name, exists
}

// GetCityID returns the city ID by English city name (case insensitive)
func GetCityID(cityName string) (uint16, bool) {
	// Try exact match first
	if id, exists := cityNameToID[cityName]; exists {
		return id, true
	}

	// Try case insensitive match
	for name, id := range cityNameToID {
		if name == cityName {
			return id, true
		}
	}

	return 0, false
}

// GetAllCities returns all city codes
func GetAllCities() []CityInfo {
	return cityCodes
}

// GetCityCount returns the total number of cities
func GetCityCount() int {
	return len(cityCodes)
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

// CityDisplayInfo represents a city with its display information
type CityDisplayInfo struct {
	Index int
	City  CityData
	Row   int
	Col   int
	Owner byte
}

// ProcessCitiesForDisplay processes cities and returns display information
func ProcessCitiesForDisplay(saveOutput *WC4SaveOutput) []CityDisplayInfo {
	var validCities []CityDisplayInfo
	for i := 0; i < len(saveOutput.Cities); i++ {
		city := saveOutput.Cities[i]
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			continue // Skip invalid coordinates
		}
		owner := saveOutput.UnitOwnerData[row][col]
		validCities = append(validCities, CityDisplayInfo{
			Index: i,
			City:  city,
			Row:   row,
			Col:   col,
			Owner: owner,
		})
	}
	return validCities
}

// GroupCitiesByOwner groups cities by their owner
func GroupCitiesByOwner(cities []CityDisplayInfo) map[byte][]CityDisplayInfo {
	citiesByOwner := make(map[byte][]CityDisplayInfo)
	for _, cityInfo := range cities {
		citiesByOwner[cityInfo.Owner] = append(citiesByOwner[cityInfo.Owner], cityInfo)
	}
	return citiesByOwner
}
