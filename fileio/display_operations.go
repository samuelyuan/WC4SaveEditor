package fileio

import (
	"fmt"
	"sort"
	"strings"
)

// GetSortedUnitTypes returns unit types sorted in ascending order
func GetSortedUnitTypes(unitsByType map[byte]int) []int {
	var unitTypes []int
	for unitType := range unitsByType {
		unitTypes = append(unitTypes, int(unitType))
	}
	sort.Ints(unitTypes)
	return unitTypes
}

// GetSortedPlayerIDs returns player IDs sorted by their index
func GetSortedPlayerIDs(playerData []CountryData) []int {
	playerIDs := make([]int, len(playerData))
	for i := range playerData {
		playerIDs[i] = i
	}
	return playerIDs
}

// DisplayHeader prints a formatted header with title and separator
func DisplayHeader(title string) {
	fmt.Println(title)
	fmt.Println(strings.Repeat("-", len(title)))
}

// DisplayPlayerInfo returns formatted player information string
func DisplayPlayerInfo(playerID int, player CountryData, count int) string {
	countryName, _ := GetCountryInfo(player.CountryId)
	return fmt.Sprintf("Player %d (%s - CountryId %d) owns %d units", playerID, countryName, player.CountryId, count)
}

// DisplaySkipStatistics shows statistics about skipped items
func DisplaySkipStatistics(skipped int, skippedByOwner map[byte]int, skippedCoordinates int) {
	if skipped > 0 {
		fmt.Printf("Skipped %d units due to invalid data\n", skipped)
		if skippedCoordinates > 0 {
			fmt.Printf("  %d units skipped due to invalid coordinates\n", skippedCoordinates)
		}
		if len(skippedByOwner) > 0 {
			fmt.Println("  Skipped units by owner:")
			for owner, count := range skippedByOwner {
				fmt.Printf("    Owner %d: %d units skipped\n", owner, count)
			}
		}
		fmt.Println()
	}
}

// SortByCount sorts a slice of items by their count in descending order
func SortByCount[T any](items []T, getCount func(T) int) {
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if getCount(items[i]) < getCount(items[j]) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
}

// DisplayTypeAnalysis shows unit type analysis with percentages
func DisplayTypeAnalysis[T any](items []T, getType func(T) uint8, getCount func(T) int, itemName string) {
	typeCounts := make(map[uint8]int)
	for _, item := range items {
		typeCounts[getType(item)]++
	}

	type TypeCount struct {
		UnitType uint8
		Count    int
		Name     string
	}
	var typeStats []TypeCount
	for unitType, count := range typeCounts {
		typeStats = append(typeStats, TypeCount{
			UnitType: unitType,
			Count:    count,
			Name:     GetUnitTypeName(unitType),
		})
	}

	SortByCount(typeStats, func(tc TypeCount) int { return tc.Count })

	totalItems := len(items)
	for _, stat := range typeStats {
		percentage := float64(stat.Count) / float64(totalItems) * 100
		fmt.Printf("Type %d (%s): %d %s (%.1f%%)\n",
			stat.UnitType, stat.Name, stat.Count, itemName, percentage)
	}
}

// DisplayCoordinatesGrid prints coordinates in a grid format (10 per row)
func DisplayCoordinatesGrid(coordinates []string) {
	countPerRow := 10

	for i, coord := range coordinates {
		if i%countPerRow == 0 {
			fmt.Printf("  ")
		}

		fmt.Printf("%s", coord)

		if (i+1)%countPerRow == 0 {
			fmt.Println()
		} else {
			fmt.Printf(" ")
		}
	}

	if len(coordinates)%countPerRow != 0 {
		fmt.Println()
	}
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

// TileAnalysisResult contains the results of tile analysis
type TileAnalysisResult struct {
	CoordinateCodeCounts map[uint16]int
	TotalTiles           int
	OceanCount           int
	LandTiles            int
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
		countryName, _ := GetCountryInfo(player.CountryId)
		fmt.Printf("Player %d: %s (CountryId %d), TeamId %d, units owned: %d\n", i, countryName, player.CountryId, player.TeamId, countMap[byte(i)])
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
			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("  Player %d: %s (CountryId %d), units owned: %d\n",
				playerIndex, countryName, player.CountryId, countMap[byte(playerIndex)])
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
		fmt.Printf("CityInfo[%d]: %+v\n", cityInfo.Index, cityInfo)
	}

	// Then, group cities by owner
	fmt.Println("\nGrouped by Owner:")
	citiesByOwner := GroupCitiesByOwner(validCities)

	// Display cities grouped by owner
	for owner := byte(0); owner < byte(len(saveOutput.PlayerData)); owner++ {
		if cities, exists := citiesByOwner[owner]; exists {
			countryName, _ := GetCountryInfo(saveOutput.PlayerData[owner].CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d) owns %d cities:\n", owner, countryName, saveOutput.PlayerData[owner].CountryId, len(cities))
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
	DisplayHeader(fmt.Sprintf("Units (Total: %d)", len(saveOutput.Units)))

	// Process all units using common function
	result := ProcessAllUnits(saveOutput)
	validUnits := result.ValidUnits

	// Show skip statistics
	DisplaySkipStatistics(result.TotalSkipped, result.SkippedUnits, result.SkippedCoordinates)

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
			playerInfo := DisplayPlayerInfo(ownerID, saveOutput.PlayerData[ownerID], len(units))
			if skippedCount > 0 {
				playerInfo += fmt.Sprintf(" (skipped %d invalid units)", skippedCount)
			}
			fmt.Printf("\n%s:\n", playerInfo)

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
	DisplayTypeAnalysis(validUnits,
		func(u UnitDisplayInfo) uint8 { return u.Unit.UnitType },
		func(u UnitDisplayInfo) int { return 1 },
		"units")
}

// ListGenerals displays all units that have generals assigned, grouped by player
func ListGenerals(saveOutput *WC4SaveOutput) {
	DisplayHeader("Generals")

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
			countryName, _ := GetCountryInfo(saveOutput.PlayerData[ownerID].CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d) has %d generals:\n",
				ownerID, countryName, saveOutput.PlayerData[ownerID].CountryId, len(generalUnits))

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
	DisplayTypeAnalysis(generals,
		func(g UnitDisplayInfo) uint8 { return g.Unit.UnitType },
		func(g UnitDisplayInfo) int { return 1 },
		"generals")

	fmt.Printf("\nTotal generals found: %d\n", len(generals))
}

// ListUnitsByMap displays units by map analysis
func ListUnitsByMap(saveOutput *WC4SaveOutput) {
	fmt.Println("Units by Owner Analysis:")
	fmt.Println("----------------------------")

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

	// Display units by owner using sorted player IDs
	fmt.Println("Units by Owner:")

	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if count, exists := ownerCounts[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			percentage := float64(count) / float64(totalUnits) * 100

			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d, TeamId %d) owns %d units (%.1f%% of total)\n",
				ownerID, countryName, player.CountryId, player.TeamId, count, percentage)

			// Aggregate coordinates for this player
			coordinates := make([]string, 0)
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

						if hasUnit {
							coordinates = append(coordinates, fmt.Sprintf("(%d,%d)", i, j))
						} else {
							coordinates = append(coordinates, fmt.Sprintf("[%d,%d]", i, j))
							missingFromUnitList++
						}
					}
				}
			}

			// Print coordinates in grid format
			DisplayCoordinatesGrid(coordinates)

			// Show summary of missing units
			if missingFromUnitList > 0 {
				fmt.Printf("  [Brackets] indicate %d tiles missing from unit list\n", missingFromUnitList)
			}
		}
	}
}

// AnalyzeCityTiles analyzes the city tiles array and returns statistics
func AnalyzeCityTiles(saveOutput *WC4SaveOutput) TileAnalysisResult {
	coordinateCodeCounts := make(map[uint16]int)
	totalTiles := 0

	for i := 0; i < len(saveOutput.CityTiles); i++ {
		for j := 0; j < len(saveOutput.CityTiles[i]); j++ {
			totalTiles++
			coordinateCode := saveOutput.CityTiles[i][j]
			coordinateCodeCounts[coordinateCode]++
		}
	}

	oceanCount := coordinateCodeCounts[65535]
	landTiles := totalTiles - oceanCount

	return TileAnalysisResult{
		CoordinateCodeCounts: coordinateCodeCounts,
		TotalTiles:           totalTiles,
		OceanCount:           oceanCount,
		LandTiles:            landTiles,
	}
}

// DisplayTileAnalysis shows the tile analysis results
func DisplayTileAnalysis(result TileAnalysisResult, saveOutput *WC4SaveOutput) {
	fmt.Println("\nCity Tiles Array Analysis:")
	fmt.Println("--------------------------")
	fmt.Printf("Total tiles analyzed: %d\n", result.TotalTiles)
	fmt.Printf("Unique coordinate codes found: %d\n", len(result.CoordinateCodeCounts))

	fmt.Printf("\nSummary: %d land tiles (%.1f%%), %d ocean tiles (%.1f%%)\n",
		result.LandTiles, float64(result.LandTiles)/float64(result.TotalTiles)*100,
		result.OceanCount, float64(result.OceanCount)/float64(result.TotalTiles)*100)
}

// ProcessCitiesForTiles processes cities and groups them by owner for tile analysis
func ProcessCitiesForTiles(saveOutput *WC4SaveOutput) (map[byte][]CityDisplayInfo, int, int, int) {
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

		validCities++
		cityDisplayInfo := CityDisplayInfo{
			Index: i,
			City:  city,
			Row:   row,
			Col:   col,
		}
		citiesByOwner[owner] = append(citiesByOwner[owner], cityDisplayInfo)
	}

	return citiesByOwner, totalCities, validCities, skippedCities
}

// DisplayCitiesByOwner shows cities grouped by owner with tile counts
func DisplayCitiesByOwner(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	fmt.Println("\nCity Tiles by Owner:")

	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]

			// Calculate total tiles for this player
			totalTilesForPlayer := 0
			for _, cityInfo := range cities {
				tileCount := result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				totalTilesForPlayer += tileCount
			}

			landPercentage := float64(totalTilesForPlayer) / float64(result.LandTiles) * 100
			worldPercentage := float64(totalTilesForPlayer) / float64(result.TotalTiles) * 100

			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d, TeamId %d) owns %d city tiles (total of %d tiles, %.1f%% of land, %.1f%% of world):\n",
				ownerID, countryName, player.CountryId, player.TeamId, len(cities), totalTilesForPlayer, landPercentage, worldPercentage)

			// Show cities with their names, positions, coordinate codes, and tile counts
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				tileCount := result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				fmt.Printf("  City %d: %s (ID:0x%02x) at (%d,%d) coordCode:%d has %d tiles\n",
					cityInfo.Index, cityName, cityInfo.City.CityId, cityInfo.Row, cityInfo.Col, cityInfo.City.CoordinateCode, tileCount)
			}
		}
	}
}

// DisplayTerritoryControl shows player territory control summary
func DisplayTerritoryControl(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	fmt.Println("\nPlayer Territory Control:")
	fmt.Println("------------------------")
	totalOwnedTiles := 0
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)

	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			playerTiles := 0
			for _, cityInfo := range cities {
				playerTiles += result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
			}
			totalOwnedTiles += playerTiles
			player := saveOutput.PlayerData[ownerID]
			landPercentage := float64(playerTiles) / float64(result.LandTiles) * 100
			worldPercentage := float64(playerTiles) / float64(result.TotalTiles) * 100
			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("Player %d (%s - CountryId %d): %d tiles (%.1f%% of land, %.1f%% of world)\n",
				ownerID, countryName, player.CountryId, playerTiles, landPercentage, worldPercentage)
		}
	}

	landPercentage := float64(totalOwnedTiles) / float64(result.LandTiles) * 100
	worldPercentage := float64(totalOwnedTiles) / float64(result.TotalTiles) * 100
	fmt.Printf("Total controlled: %d tiles (%.1f%% of land, %.1f%% of world)\n",
		totalOwnedTiles, landPercentage, worldPercentage)
}

// ListTilesByOwner displays city tiles by owner analysis
func ListTilesByOwner(saveOutput *WC4SaveOutput) {
	DisplayHeader("City Tile Ownership Analysis")

	// Analyze city tiles
	result := AnalyzeCityTiles(saveOutput)
	DisplayTileAnalysis(result, saveOutput)

	// Process cities and group by owner
	citiesByOwner, totalCities, validCities, skippedCities := ProcessCitiesForTiles(saveOutput)

	fmt.Printf("Total cities: %d\n", totalCities)
	fmt.Printf("Valid cities: %d\n", validCities)
	fmt.Printf("Skipped cities: %d\n", skippedCities)

	// Display cities grouped by owner
	DisplayCitiesByOwner(saveOutput, citiesByOwner, result)

	// Show territory control summary
	DisplayTerritoryControl(saveOutput, citiesByOwner, result)
}
