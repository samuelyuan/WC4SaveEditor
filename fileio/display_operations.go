package fileio

import (
	"fmt"
)

// Unit type constants
const (
	UnitTypeLightInfantry      = 1  // Light Infantry
	UnitTypeAssaultInfantry    = 2  // Assault Infantry
	UnitTypeMotorizedInfantry  = 3  // Motorized Infantry
	UnitTypeMechanizedInfantry = 4  // Mechanized Infantry
	UnitTypeArmoredCar         = 6  // Armored Car
	UnitTypeLightTank          = 7  // Light Tank
	UnitTypeMediumTank         = 8  // Medium Tank
	UnitTypeHeavyTank          = 9  // Heavy Tank
	UnitTypeSuperTank          = 10 // Super Tank
	UnitTypeFieldArtillery     = 11 // Field Artillery
	UnitTypeHowitzer           = 12 // Howitzer
	UnitTypeRocketArtillery    = 13 // Rocket Artillery
	UnitTypeSuperArtillery     = 14 // Super Artillery
	UnitTypeSubmarine          = 15 // Submarine
	UnitTypeDestroyer          = 16 // Destroyer
	UnitTypeCruiser            = 17 // Cruiser
	UnitTypeCarrier            = 18 // Carrier
	UnitTypeBunker             = 35 // Bunker
	UnitTypeLandFort           = 36 // Land Fort
	UnitTypeCoastalArtillery   = 37 // Coastal Artillery
	UnitTypeRocketLauncher     = 38 // Rocket Launcher
)

// GetSortedPlayerIDs returns player IDs sorted by their index
func GetSortedPlayerIDs(playerData []CountryData) []int {
	playerIDs := make([]int, len(playerData))
	for i := range playerData {
		playerIDs[i] = i
	}
	return playerIDs
}

// UnitDisplayInfo represents a unit with its display information
type UnitDisplayInfo struct {
	Index int
	Unit  UnitData
	Row   int
	Col   int
	Owner byte
}

// TileInfo represents a tile position
type TileInfo struct {
	Row, Col int
}

// ProcessAllUnitsResult contains both valid units and statistics about skipped units
type ProcessAllUnitsResult struct {
	ValidUnits         []UnitDisplayInfo
	SkippedUnits       map[byte]int // Count of skipped units per owner (invalid owner)
	SkippedCoordinates int          // Count of units skipped due to invalid coordinates
	TotalSkipped       int
}

// ProcessAllUnits processes all units and returns valid units with coordinates and skip statistics
func ProcessAllUnits(saveOutput *WC4SaveOutput) ProcessAllUnitsResult {
	var validUnits []UnitDisplayInfo
	skippedUnits := make(map[byte]int)
	skippedCoordinates := 0
	totalSkipped := 0

	for i := 0; i < len(saveOutput.Units); i++ {
		unit := saveOutput.Units[i]
		row, col, valid := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			fmt.Printf("WARNING: Unit %d: skipping invalid coordinates (CoordinateCode=%d)\n", i, unit.CoordinateCode)
			skippedCoordinates++
			totalSkipped++
			continue
		}
		owner := saveOutput.UnitOwnerData[row][col]

		// Check if owner is valid
		if int(owner) >= len(saveOutput.PlayerData) {
			fmt.Printf("WARNING: Unit %d: invalid owner %d, skipping\n", i, owner)
			skippedUnits[owner]++
			totalSkipped++
			continue
		}

		validUnits = append(validUnits, UnitDisplayInfo{
			Index: i,
			Unit:  unit,
			Row:   row,
			Col:   col,
			Owner: owner,
		})
	}

	return ProcessAllUnitsResult{
		ValidUnits:         validUnits,
		SkippedUnits:       skippedUnits,
		SkippedCoordinates: skippedCoordinates,
		TotalSkipped:       totalSkipped,
	}
}

// GroupUnitsByOwner groups units by their owner
func GroupUnitsByOwner(units []UnitDisplayInfo) map[byte][]UnitDisplayInfo {
	unitsByOwner := make(map[byte][]UnitDisplayInfo)
	for _, unitInfo := range units {
		unitsByOwner[unitInfo.Owner] = append(unitsByOwner[unitInfo.Owner], unitInfo)
	}
	return unitsByOwner
}

// CountPlayerTiles counts tiles owned by each player
func CountPlayerTiles(unitOwnerData [][]byte) map[byte]int {
	countMap := make(map[byte]int)
	for i := 0; i < len(unitOwnerData); i++ {
		for j := 0; j < len(unitOwnerData[i]); j++ {
			if unitOwnerData[i][j] == TileUnowned {
				continue
			}
			owner := unitOwnerData[i][j]
			existingCount, ok := countMap[owner]
			if !ok {
				countMap[owner] = 1
			} else {
				countMap[owner] = existingCount + 1
			}
		}
	}
	return countMap
}

// GetGeneralName returns a human-readable name for the general ID
func GetGeneralName(generalId uint16) (string, bool) {
	switch generalId {
	case 25182:
		return "Katukov", true
	case 25191:
		return "Kimmel", true
	case 25199:
		return "Pavlov", true
	case 25212:
		return "Abrams", true
	case 25217:
		return "Leslie", true
	case 29001:
		return "Zhang.Z.Z", true
	case 29002:
		return "Sun.L.R", true
	case 29003:
		return "Zhu.D", true
	case 29004:
		return "Peng.D.H", true
	case 29005:
		return "Slim", true
	case 29006:
		return "Mountbatten", true
	case 29009:
		return "Montgomery", true
	case 29012:
		return "Badoglio", true
	case 29013:
		return "Graziani", true
	case 29014:
		return "Chuikov", true
	case 29016:
		return "Bagramyan", true
	case 29017:
		return "Kuznetsov", true
	case 29019:
		return "Konev", true
	case 29020:
		return "Govorov", true
	case 29021:
		return "Rokossovsky", true
	case 29024:
		return "Yamashita", true
	case 29025:
		return "Kuribayashi", true
	case 29027:
		return "Yamamoto", true
	case 29029:
		return "Tito", true
	case 29030:
		return "MacArthur", true
	case 29031:
		return "Eisenhower", true
	case 29032:
		return "Nimitz", true
	case 29035:
		return "Arnold", true
	case 29038:
		return "Crerar", true
	case 29039:
		return "Mannerheim", true
	case 29040:
		return "Tassigny", true
	case 29041:
		return "de Gaulle", true
	case 29042:
		return "Leclerc", true
	case 29043:
		return "Bock", true
	case 29046:
		return "Donitz", true
	case 29048:
		return "Leeb", true
	case 29049:
		return "Guderian", true
	case 29050:
		return "Manstein", true
	case 29052:
		return "Model", true
	case 29053:
		return "Smigly", true
	case 29054:
		return "Blamey", true
	case 29055:
		return "Nasser", true
	case 29060:
		return "Student", true
	case 29066:
		return "Keitel", true
	case 29070:
		return "Paulus", true
	case 29071:
		return "Meyer", true
	case 29072:
		return "Petain", true
	case 29073:
		return "Gamelin", true
	case 29074:
		return "Darlan", true
	case 29075:
		return "Juin", true
	case 29076:
		return "Leopold", true
	case 29077:
		return "Winkelman", true
	case 29078:
		return "Christian", true
	case 29079:
		return "Olav", true
	case 29081:
		return "Voronov", true
	case 29082:
		return "Meretskov", true
	case 29085:
		return "Shaposhnikov", true
	case 29087:
		return "Voroshilov", true
	case 29089:
		return "Antonescu", true
	case 29090:
		return "Dumitrescu", true
	case 29091:
		return "Horthy", true
	case 29092:
		return "Riccardi", true
	case 29093:
		return "Campioni", true
	case 29094:
		return "Cavallero", true
	case 29095:
		return "Balbo", true
	case 29096:
		return "Cunningham", true
	case 29097:
		return "Wavell", true
	case 29098:
		return "Pound", true
	case 29099:
		return "Wingate", true
	case 29100:
		return "Dill", true
	case 29103:
		return "Papagos", true
	case 29104:
		return "Boris", true
	case 29105:
		return "Nagano", true
	case 29106:
		return "Okamura", true
	case 29109:
		return "Hata", true
	case 29110:
		return "Ozawa", true
	case 29111:
		return "Umezu", true
	case 29117:
		return "Inoue", true
	case 29118:
		return "Li.Z.R", true
	case 29119:
		return "Bai.C.X", true
	case 29120:
		return "Xue.Y", true
	case 29121:
		return "Liang.X.C", true
	case 29122:
		return "Du.Y.M", true
	case 29124:
		return "Chen.S.K", true
	case 29125:
		return "Lin.B", true
	case 29129:
		return "Phibun", true
	case 29130:
		return "Thimayya", true
	case 29131:
		return "Crace", true
	case 29132:
		return "Franco", true
	case 29136:
		return "Devers", true
	case 29137:
		return "Eaker", true
	case 29139:
		return "Smith", true
	case 29140:
		return "King", true
	case 29141:
		return "Stilwell", true
	case 29151:
		return "Inonu", true
	case 29152:
		return "Dutra", true
	case 29153:
		return "Camacho", true
	default:
		return "", false
	}
}

// GetUnitTypeName returns a human-readable name for the unit type
func GetUnitTypeName(unitType uint8) string {
	switch unitType {
	case UnitTypeLightInfantry:
		return "Light Infantry"
	case UnitTypeAssaultInfantry:
		return "Assault Infantry"
	case UnitTypeMotorizedInfantry:
		return "Motorized Infantry"
	case UnitTypeMechanizedInfantry:
		return "Mechanized Infantry"
	case UnitTypeArmoredCar:
		return "Armored Car"
	case UnitTypeLightTank:
		return "Light Tank"
	case UnitTypeMediumTank:
		return "Medium Tank"
	case UnitTypeHeavyTank:
		return "Heavy Tank"
	case UnitTypeSuperTank:
		return "Super Tank"
	case UnitTypeFieldArtillery:
		return "Field Artillery"
	case UnitTypeHowitzer:
		return "Howitzer"
	case UnitTypeRocketArtillery:
		return "Rocket Artillery"
	case UnitTypeSuperArtillery:
		return "Super Artillery"
	case UnitTypeSubmarine:
		return "Submarine"
	case UnitTypeDestroyer:
		return "Destroyer"
	case UnitTypeCruiser:
		return "Cruiser"
	case UnitTypeCarrier:
		return "Carrier"
	case UnitTypeBunker:
		return "Bunker"
	case UnitTypeLandFort:
		return "Land Fort"
	case UnitTypeCoastalArtillery:
		return "Coastal Artillery"
	case UnitTypeRocketLauncher:
		return "Rocket Launcher"
	case 39: // UnitTypeCity - defined in writer.go
		return "City"
	default:
		return fmt.Sprintf("Unknown Type %d", unitType)
	}
}

// ListPlayers displays all players with their information
func ListPlayers(saveOutput *WC4SaveOutput) {
	countMap := CountPlayerTiles(saveOutput.UnitOwnerData)

	fmt.Println("Players:")
	fmt.Println("--------")

	// First, show linear list
	fmt.Println("Linear List:")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, i := range sortedPlayerIDs {
		player := saveOutput.PlayerData[i]
		countryName, countryInfo := GetCountryInfo(player.CountryId)
		fmt.Printf("Player %d: %s (%s), TeamId %d, units owned: %d\n", i, countryName, countryInfo, player.TeamId, countMap[byte(i)])
	}

	// Then, group players by TeamId
	fmt.Println("\nGrouped by Team:")
	playersByTeam := make(map[uint32][]int)
	for i := 0; i < len(saveOutput.PlayerData); i++ {
		teamId := saveOutput.PlayerData[i].TeamId
		playersByTeam[teamId] = append(playersByTeam[teamId], i)
	}

	// Display players grouped by team
	for teamId, playerIndices := range playersByTeam {
		fmt.Printf("\nTeam %d (%d players):\n", teamId, len(playerIndices))
		for _, playerIndex := range playerIndices {
			player := saveOutput.PlayerData[playerIndex]
			countryName, countryInfo := GetCountryInfo(player.CountryId)
			fmt.Printf("  Player %d: %s (%s), units owned: %d\n",
				playerIndex, countryName, countryInfo, countMap[byte(playerIndex)])
		}
	}
}

// ListCities displays cities in both linear and grouped format
func ListCities(saveOutput *WC4SaveOutput) {
	fmt.Printf("Cities (Map: %dx%d):\n", len(saveOutput.UnitOwnerData), len(saveOutput.UnitOwnerData[0]))
	fmt.Println("-------------------")

	// Process cities using the fileio package
	validCities := ProcessCitiesForDisplay(saveOutput)

	// First, show linear list
	fmt.Println("Linear List:")
	for _, cityInfo := range validCities {
		cityName, hasName := GetCityName(cityInfo.City.CityId)
		if !hasName {
			cityName = "Unknown"
		}
		fmt.Printf("City %d: %s (ID:0x%02x), Owner=%d, CityID=%d, Position=(%d,%d)\n",
			cityInfo.Index, cityName, cityInfo.City.CityId, cityInfo.Owner, cityInfo.City.CityId, cityInfo.Row, cityInfo.Col)
	}

	// Then, group cities by owner
	fmt.Println("\nGrouped by Owner:")
	citiesByOwner := GroupCitiesByOwner(validCities)

	// Display cities grouped by owner
	for owner := byte(0); owner < byte(len(saveOutput.PlayerData)); owner++ {
		if cities, exists := citiesByOwner[owner]; exists {
			countryName, countryInfo := GetCountryInfo(saveOutput.PlayerData[owner].CountryId)
			fmt.Printf("\nPlayer %d (%s - %s) owns %d cities:\n", owner, countryName, countryInfo, len(cities))
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				fmt.Printf("  City %d: %s (ID:0x%02x), Position=(%d,%d)\n",
					cityInfo.Index, cityName, cityInfo.City.CityId, cityInfo.Row, cityInfo.Col)
			}
		}
	}
}

// ListUnits displays units in both linear and grouped format
func ListUnits(saveOutput *WC4SaveOutput) {
	fmt.Printf("Units (Total: %d):\n", len(saveOutput.Units))
	fmt.Println("-----------------")

	// Process all units using common function
	result := ProcessAllUnits(saveOutput)
	validUnits := result.ValidUnits

	// Show skip statistics
	if result.TotalSkipped > 0 {
		fmt.Printf("Skipped %d units due to invalid data\n", result.TotalSkipped)
		if result.SkippedCoordinates > 0 {
			fmt.Printf("  %d units skipped due to invalid coordinates\n", result.SkippedCoordinates)
		}
		if len(result.SkippedUnits) > 0 {
			fmt.Println("  Skipped units by owner:")
			for owner, count := range result.SkippedUnits {
				fmt.Printf("    Owner %d: %d units skipped\n", owner, count)
			}
		}
		fmt.Println()
	}

	// First, show linear list
	fmt.Println("Linear List:")
	for _, unitInfo := range validUnits {
		fmt.Printf("Unit %d (owner: %d): %+v\n", unitInfo.Index, unitInfo.Owner, unitInfo.Unit)
	}

	// Then, group units by owner
	fmt.Println("\nGrouped by Owner:")
	unitsByOwner := GroupUnitsByOwner(validUnits)

	// Display units grouped by owner using sorted player IDs
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if units, exists := unitsByOwner[byte(ownerID)]; exists {
			skippedCount := result.SkippedUnits[byte(ownerID)]
			countryName, countryInfo := GetCountryInfo(saveOutput.PlayerData[ownerID].CountryId)
			fmt.Printf("\nPlayer %d (%s - %s) owns %d units", ownerID, countryName, countryInfo, len(units))
			if skippedCount > 0 {
				fmt.Printf(" (skipped %d invalid units)", skippedCount)
			}
			fmt.Println(":")
			for _, unitInfo := range units {
				unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
				unitDescription := fmt.Sprintf("Unit %d: Type=%d (%s), Level=%d, Health=%d/%d, Position=(%d,%d)",
					unitInfo.Index, unitInfo.Unit.UnitType, unitTypeName, unitInfo.Unit.Level,
					unitInfo.Unit.CurrentHealth, unitInfo.Unit.MaxHealth, unitInfo.Row, unitInfo.Col)

				// If this is a city unit, try to get the city name
				if unitInfo.Unit.UnitType == UnitTypeCity {
					cityName := GetCityNameAtPosition(unitInfo.Row, unitInfo.Col, saveOutput)
					unitDescription += fmt.Sprintf(" [City: %s]", cityName)
				}

				fmt.Printf("  %s\n", unitDescription)
			}
		}
	}

	// Add unit type analysis
	fmt.Println("\nUnit Type Analysis:")
	fmt.Println("------------------")
	unitTypeCounts := make(map[uint8]int)
	for _, unitInfo := range validUnits {
		unitTypeCounts[unitInfo.Unit.UnitType]++
	}

	// Sort unit types by count (descending)
	type UnitTypeCount struct {
		UnitType uint8
		Count    int
		Name     string
	}
	var unitTypeStats []UnitTypeCount
	for unitType, count := range unitTypeCounts {
		unitTypeStats = append(unitTypeStats, UnitTypeCount{
			UnitType: unitType,
			Count:    count,
			Name:     GetUnitTypeName(unitType),
		})
	}

	// Simple bubble sort by count (descending)
	for i := 0; i < len(unitTypeStats); i++ {
		for j := i + 1; j < len(unitTypeStats); j++ {
			if unitTypeStats[i].Count < unitTypeStats[j].Count {
				unitTypeStats[i], unitTypeStats[j] = unitTypeStats[j], unitTypeStats[i]
			}
		}
	}

	totalUnits := len(validUnits)
	for _, stat := range unitTypeStats {
		percentage := float64(stat.Count) / float64(totalUnits) * 100
		fmt.Printf("Type %d (%s): %d units (%.1f%%)\n",
			stat.UnitType, stat.Name, stat.Count, percentage)
	}
}

// ListGenerals displays all units that have generals assigned, grouped by player
func ListGenerals(saveOutput *WC4SaveOutput) {
	fmt.Println("Generals:")
	fmt.Println("---------")

	// First, collect all generals with their information
	var generals []UnitDisplayInfo
	skippedGenerals := 0

	for i := 0; i < len(saveOutput.Units); i++ {
		unit := saveOutput.Units[i]
		row, col, valid := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			fmt.Printf("WARNING: Unit %d: invalid coordinates for general check (CoordinateCode=%d)\n", i, unit.CoordinateCode)
			skippedGenerals++
			continue
		}
		if unit.GeneralId > 0 {
			owner := saveOutput.UnitOwnerData[row][col]

			// Check if owner is valid
			if int(owner) >= len(saveOutput.PlayerData) {
				fmt.Printf("WARNING: General unit %d: invalid owner %d, skipping\n", i, owner)
				skippedGenerals++
				continue
			}

			generals = append(generals, UnitDisplayInfo{
				Index: i,
				Unit:  unit,
				Row:   row,
				Col:   col,
				Owner: owner,
			})
		}
	}

	if len(generals) == 0 {
		fmt.Println("No generals found.")
		if skippedGenerals > 0 {
			fmt.Printf("(Skipped %d units due to invalid data)\n", skippedGenerals)
		}
		return
	}

	// Show skip statistics
	if skippedGenerals > 0 {
		fmt.Printf("Skipped %d units due to invalid data\n", skippedGenerals)
		fmt.Println()
	}

	// First, show linear list
	fmt.Println("Linear List:")
	for _, generalInfo := range generals {
		unitTypeName := GetUnitTypeName(generalInfo.Unit.UnitType)
		countryName, _ := GetCountryInfo(saveOutput.PlayerData[generalInfo.Owner].CountryId)
		generalName, hasGeneralName := GetGeneralName(generalInfo.Unit.GeneralId)
		
		generalInfoStr := fmt.Sprintf("GeneralId=%d", generalInfo.Unit.GeneralId)
		if hasGeneralName {
			generalInfoStr = fmt.Sprintf("GeneralId=%d (%s)", generalInfo.Unit.GeneralId, generalName)
		}
		
		fmt.Printf("General (unit %d, owner %d - %s): Type=%d (%s), %s, Level=%d, Health=%d/%d, Position=(%d,%d)\n", 
			generalInfo.Index, generalInfo.Owner, countryName, generalInfo.Unit.UnitType, unitTypeName, 
			generalInfoStr, generalInfo.Unit.Level, generalInfo.Unit.CurrentHealth, 
			generalInfo.Unit.MaxHealth, generalInfo.Row, generalInfo.Col)
	}

	// Then, group generals by owner
	fmt.Println("\nGrouped by Owner:")
	generalsByOwner := GroupUnitsByOwner(generals)

	// Display generals grouped by owner using sorted player IDs
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if generalUnits, exists := generalsByOwner[byte(ownerID)]; exists {
			countryName, countryInfo := GetCountryInfo(saveOutput.PlayerData[ownerID].CountryId)
			fmt.Printf("\nPlayer %d (%s - %s) has %d generals:\n",
				ownerID, countryName, countryInfo, len(generalUnits))

			for _, generalInfo := range generalUnits {
				unitTypeName := GetUnitTypeName(generalInfo.Unit.UnitType)
				generalName, hasGeneralName := GetGeneralName(generalInfo.Unit.GeneralId)
				
				generalInfoStr := fmt.Sprintf("GeneralId=%d", generalInfo.Unit.GeneralId)
				if hasGeneralName {
					generalInfoStr = fmt.Sprintf("GeneralId=%d (%s)", generalInfo.Unit.GeneralId, generalName)
				}
				
				fmt.Printf("  General (unit %d): Type=%d (%s), %s, Level=%d, Health=%d/%d, Position=(%d,%d)\n",
					generalInfo.Index, generalInfo.Unit.UnitType, unitTypeName, generalInfoStr,
					generalInfo.Unit.Level, generalInfo.Unit.CurrentHealth, generalInfo.Unit.MaxHealth, 
					generalInfo.Row, generalInfo.Col)
			}
		}
	}

	// Add general type analysis
	fmt.Println("\nGeneral Unit Type Analysis:")
	fmt.Println("---------------------------")
	generalTypeCounts := make(map[uint8]int)
	for _, generalInfo := range generals {
		generalTypeCounts[generalInfo.Unit.UnitType]++
	}

	// Sort general unit types by count (descending)
	type GeneralTypeCount struct {
		UnitType uint8
		Count    int
		Name     string
	}
	var generalTypeStats []GeneralTypeCount
	for unitType, count := range generalTypeCounts {
		generalTypeStats = append(generalTypeStats, GeneralTypeCount{
			UnitType: unitType,
			Count:    count,
			Name:     GetUnitTypeName(unitType),
		})
	}

	// Simple bubble sort by count (descending)
	for i := 0; i < len(generalTypeStats); i++ {
		for j := i + 1; j < len(generalTypeStats); j++ {
			if generalTypeStats[i].Count < generalTypeStats[j].Count {
				generalTypeStats[i], generalTypeStats[j] = generalTypeStats[j], generalTypeStats[i]
			}
		}
	}

	totalGenerals := len(generals)
	for _, stat := range generalTypeStats {
		percentage := float64(stat.Count) / float64(totalGenerals) * 100
		fmt.Printf("Type %d (%s): %d generals (%.1f%%)\n",
			stat.UnitType, stat.Name, stat.Count, percentage)
	}

	fmt.Printf("\nTotal generals found: %d\n", totalGenerals)
}

// ListUnitsByMap displays units by map analysis
func ListUnitsByMap(saveOutput *WC4SaveOutput, playerID ...int) {
	// Check if we're filtering by a specific player
	filterByPlayer := len(playerID) > 0 && playerID[0] >= 0

	if filterByPlayer {
		fmt.Printf("Units owned by Player %d:\n", playerID[0])
		fmt.Println("------------------------------")
	} else {
		fmt.Println("Units by Owner Analysis:")
		fmt.Println("----------------------------")
	}

	// Analyze UnitOwnerData map directly
	ownerCounts := make(map[byte]int)
	totalUnits := 0

	// Count tiles by owner
	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			owner := saveOutput.UnitOwnerData[i][j]

			// Check if owner is valid
			if int(owner) >= len(saveOutput.PlayerData) {
				continue
			}

			ownerCounts[owner]++
			totalUnits++
		}
	}

	// Show statistics
	fmt.Printf("Total units: %d\n", totalUnits)
	fmt.Println()

	// Display tiles by owner using sorted player IDs
	if filterByPlayer {
		fmt.Println("Units:")
	} else {
		fmt.Println("Units by Owner:")
	}

	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if count, exists := ownerCounts[byte(ownerID)]; exists {
			// If filtering by player, only show that player
			if filterByPlayer && ownerID != playerID[0] {
				continue
			}

			player := saveOutput.PlayerData[ownerID]
			percentage := float64(count) / float64(totalUnits) * 100

			if filterByPlayer {
				fmt.Printf("\nPlayer %d owns %d units (%.1f%% of total)\n",
					ownerID, count, percentage)
			} else {
				countryName, countryInfo := GetCountryInfo(player.CountryId)
				fmt.Printf("\nPlayer %d (%s - %s, TeamId %d) owns %d units (%.1f%% of total)\n",
					ownerID, countryName, countryInfo, player.TeamId, count, percentage)
			}

			// Show units in a grid format (10 per row)
			countPerRow := 10
			count := 0
			missingFromUnitList := 0

			for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
				for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
					if saveOutput.UnitOwnerData[i][j] == byte(ownerID) {
						// Derive coordinate code from (i, j)
						coordinateCode := ConvertToCoordinateCode(i, j, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))

						// Check if any unit has this coordinate code
						hasUnit := false
						for _, unit := range saveOutput.Units {
							if int(unit.CoordinateCode) == coordinateCode {
								hasUnit = true
								break
							}
						}

						if count%countPerRow == 0 {
							fmt.Printf("  ")
						}

						if hasUnit {
							fmt.Printf("(%d,%d)", i, j)
						} else {
							fmt.Printf("[%d,%d]", i, j) // Brackets indicate missing from unit list
							missingFromUnitList++
						}

						if (count+1)%countPerRow == 0 {
							fmt.Println()
						} else {
							fmt.Printf(" ")
						}
						count++
					}
				}
			}
			if count%countPerRow != 0 {
				fmt.Println()
			}

			// Show summary of missing units
			if missingFromUnitList > 0 {
				fmt.Printf("  [Brackets] indicate %d tiles missing from unit list\n", missingFromUnitList)
			}
		}
	}
}

// ListTilesByOwner displays city tiles by owner analysis
func ListTilesByOwner(saveOutput *WC4SaveOutput, playerID ...int) {
	// Check if we're filtering by a specific player
	filterByPlayer := len(playerID) > 0 && playerID[0] >= 0

	if filterByPlayer {
		fmt.Printf("City tiles owned by Player %d:\n", playerID[0])
		fmt.Println("------------------------------")
	} else {
		fmt.Println("City Tile Ownership Analysis:")
		fmt.Println("----------------------------")
	}

	// Analyze City Tiles 2D array - count coordinate codes
	fmt.Println("\nCity Tiles Array Analysis:")
	fmt.Println("--------------------------")
	coordinateCodeCounts := make(map[uint16]int)
	totalTiles := 0

	for i := 0; i < len(saveOutput.CityTiles); i++ {
		for j := 0; j < len(saveOutput.CityTiles[i]); j++ {
			totalTiles++
			coordinateCode := saveOutput.CityTiles[i][j]
			coordinateCodeCounts[coordinateCode]++
		}
	}

	fmt.Printf("Total tiles analyzed: %d\n", totalTiles)
	fmt.Printf("Unique coordinate codes found: %d\n", len(coordinateCodeCounts))

	// Find coordinate codes that don't match any known city
	knownCoordinateCodes := make(map[uint16]bool)
	for _, city := range saveOutput.Cities {
		knownCoordinateCodes[city.CoordinateCode] = true
	}

	// Show orphaned coordinate codes (not matching any known city)
	orphanedCodes := 0
	orphanedTiles := 0
	fmt.Println("\nOrphaned Coordinate Codes (not matching any known city):")
	fmt.Println("-------------------------------------------------------")
	for coordCode, count := range coordinateCodeCounts {
		if coordCode != 65535 && !knownCoordinateCodes[coordCode] {
			orphanedCodes++
			orphanedTiles += count
			fmt.Printf("Coordinate Code %d: %d tiles (%.1f%%)\n",
				coordCode, count, float64(count)/float64(totalTiles)*100)
		}
	}

	if orphanedCodes == 0 {
		fmt.Println("No orphaned coordinate codes found.")
	}

	oceanCount := coordinateCodeCounts[65535]
	landTiles := totalTiles - oceanCount
	ownedTiles := landTiles - orphanedTiles

	fmt.Printf("\nSummary: %d land tiles (%.1f%%), %d ocean tiles (%.1f%%)\n",
		landTiles, float64(landTiles)/float64(totalTiles)*100,
		oceanCount, float64(oceanCount)/float64(totalTiles)*100)
	fmt.Printf("Land breakdown: %d owned tiles (%.1f%% of land), %d orphaned tiles (%.1f%% of land)\n",
		ownedTiles, float64(ownedTiles)/float64(landTiles)*100,
		orphanedTiles, float64(orphanedTiles)/float64(landTiles)*100)

	// Process all cities and group by owner, combining with coordinate code counts
	citiesByOwner := make(map[byte][]CityDisplayInfo)
	totalCities := 0
	validCities := 0
	skippedCities := 0

	for i := 0; i < len(saveOutput.Cities); i++ {
		city := saveOutput.Cities[i]
		totalCities++

		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			fmt.Printf("WARNING: City %d: skipping invalid coordinates (CoordinateCode=%d)\n", i, city.CoordinateCode)
			skippedCities++
			continue
		}

		owner := saveOutput.UnitOwnerData[row][col]
		if int(owner) >= len(saveOutput.PlayerData) {
			fmt.Printf("WARNING: City %d: invalid owner %d, skipping\n", i, owner)
			skippedCities++
			continue
		}

		// If filtering by player, only include cities owned by that player
		if filterByPlayer && owner != byte(playerID[0]) {
			continue
		}

		validCities++
		cityDisplayInfo := CityDisplayInfo{
			Index: i,
			City:  city,
			Row:   row,
			Col:   col,
		}
		citiesByOwner[owner] = append(citiesByOwner[owner], cityDisplayInfo)
	}

	if !filterByPlayer {
		fmt.Printf("Total cities: %d\n", totalCities)
		fmt.Printf("Valid cities: %d\n", validCities)
		fmt.Printf("Skipped cities: %d\n", skippedCities)
	}

	// Display cities grouped by owner with tile counts
	if filterByPlayer {
		fmt.Println("\nCity Tiles:")
	} else {
		fmt.Println("\nCity Tiles by Owner:")
	}

	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			// If filtering by player, only show that player
			if filterByPlayer && ownerID != playerID[0] {
				continue
			}

			player := saveOutput.PlayerData[ownerID]

			// Calculate total tiles for this player
			totalTilesForPlayer := 0
			for _, cityInfo := range cities {
				tileCount := coordinateCodeCounts[cityInfo.City.CoordinateCode]
				totalTilesForPlayer += tileCount
			}

			landPercentage := float64(totalTilesForPlayer) / float64(landTiles) * 100
			worldPercentage := float64(totalTilesForPlayer) / float64(totalTiles) * 100

			if filterByPlayer {
				fmt.Printf("\nPlayer %d owns %d city tiles (total of %d tiles, %.1f%% of land, %.1f%% of world):\n",
					ownerID, len(cities), totalTilesForPlayer, landPercentage, worldPercentage)
			} else {
				countryName, countryInfo := GetCountryInfo(player.CountryId)
				fmt.Printf("\nPlayer %d (%s - %s, TeamId %d) owns %d city tiles (total of %d tiles, %.1f%% of land, %.1f%% of world):\n",
					ownerID, countryName, countryInfo, player.TeamId, len(cities), totalTilesForPlayer, landPercentage, worldPercentage)
			}

			// Show cities with their names, positions, coordinate codes, and tile counts
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				tileCount := coordinateCodeCounts[cityInfo.City.CoordinateCode]
				fmt.Printf("  City %d: %s (ID:0x%02x) at (%d,%d) coordCode:%d has %d tiles\n",
					cityInfo.Index, cityName, cityInfo.City.CityId, cityInfo.Row, cityInfo.Col, cityInfo.City.CoordinateCode, tileCount)
			}
		}
	}

	// Show player percentage distribution
	if !filterByPlayer {
		fmt.Println("\nPlayer Territory Control:")
		fmt.Println("------------------------")
		totalOwnedTiles := 0
		for _, ownerID := range sortedPlayerIDs {
			if cities, exists := citiesByOwner[byte(ownerID)]; exists {
				playerTiles := 0
				for _, cityInfo := range cities {
					playerTiles += coordinateCodeCounts[cityInfo.City.CoordinateCode]
				}
				totalOwnedTiles += playerTiles
				player := saveOutput.PlayerData[ownerID]
				landPercentage := float64(playerTiles) / float64(landTiles) * 100
				worldPercentage := float64(playerTiles) / float64(totalTiles) * 100
				countryName, countryInfo := GetCountryInfo(player.CountryId)
				fmt.Printf("Player %d (%s - %s): %d tiles (%.1f%% of land, %.1f%% of world)\n",
					ownerID, countryName, countryInfo, playerTiles, landPercentage, worldPercentage)
			}
		}
		landPercentage := float64(totalOwnedTiles) / float64(landTiles) * 100
		worldPercentage := float64(totalOwnedTiles) / float64(totalTiles) * 100
		fmt.Printf("Total controlled: %d tiles (%.1f%% of land, %.1f%% of world)\n",
			totalOwnedTiles, landPercentage, worldPercentage)
	}
}
