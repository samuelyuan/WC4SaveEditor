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

	// Use table formatting for better display
	tableFormatter := NewTableFormatter()
	headers := []ColumnConfig{
		{"Type", 6, "right"},
		{"Unit Type", 25, "left"},
		{"Count", 8, "right"},
		{"Percentage", 12, "right"},
	}

	var data [][]string
	totalItems := len(items)
	for _, stat := range typeStats {
		percentage := float64(stat.Count) / float64(totalItems) * 100
		row := []string{
			fmt.Sprintf("%d", stat.UnitType),
			stat.Name,
			fmt.Sprintf("%d", stat.Count),
			fmt.Sprintf("%.1f%%", percentage),
		}
		data = append(data, row)
	}

	tableFormatter.PrintTable(headers, data)
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
	tableFormatter := NewTableFormatter()

	fmt.Println("Players:")
	fmt.Println("--------")

	// First, show linear list
	fmt.Println("Linear List:")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)

	// Prepare data for table
	var playerData []PlayerTableData
	for _, i := range sortedPlayerIDs {
		player := saveOutput.PlayerData[i]
		countryName, _ := GetCountryInfo(player.CountryId)
		playerData = append(playerData, PlayerTableData{
			PlayerID:    i,
			CountryName: countryName,
			CountryID:   int(player.CountryId),
			TeamID:      int(player.TeamId),
			UnitsOwned:  countMap[byte(i)],
		})
	}

	tableFormatter.PrintPlayerTable(playerData)

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

		// Prepare data for team table (without Team column)
		var teamData []PlayerTableData
		for _, playerIndex := range playerIndices {
			player := saveOutput.PlayerData[playerIndex]
			countryName, _ := GetCountryInfo(player.CountryId)
			teamData = append(teamData, PlayerTableData{
				PlayerID:    playerIndex,
				CountryName: countryName,
				CountryID:   int(player.CountryId),
				TeamID:      int(player.TeamId),
				UnitsOwned:  countMap[byte(playerIndex)],
			})
		}

		// Use custom table for team display (without Team column)
		headers := []ColumnConfig{
			{"Player", 8, "right"},
			{"Country", 17, "left"},
			{"Country ID", 13, "right"},
			{"Units Owned", 13, "right"},
		}

		var data [][]string
		for _, player := range teamData {
			row := []string{
				fmt.Sprintf("%d", player.PlayerID),
				player.CountryName,
				fmt.Sprintf("%d", player.CountryID),
				fmt.Sprintf("%d", player.UnitsOwned),
			}
			data = append(data, row)
		}

		tableFormatter.PrintTable(headers, data)
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
	tableFormatter := NewTableFormatter()
	for owner := byte(0); owner < byte(len(saveOutput.PlayerData)); owner++ {
		if cities, exists := citiesByOwner[owner]; exists {
			countryName, _ := GetCountryInfo(saveOutput.PlayerData[owner].CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d) owns %d cities:\n", owner, countryName, saveOutput.PlayerData[owner].CountryId, len(cities))

			// Prepare data for table
			var cityData []CityTableData
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				positionStr := fmt.Sprintf("(%d,%d)", cityInfo.Row, cityInfo.Col)

				cityData = append(cityData, CityTableData{
					CityID:   cityInfo.Index,
					CityName: cityName,
					Position: positionStr,
				})
			}

			tableFormatter.PrintCityTable(cityData)
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

			// Prepare data for table
			var unitData []UnitTableData
			for _, unitInfo := range units {
				unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
				healthStr := fmt.Sprintf("%d/%d", unitInfo.Unit.CurrentHealth, unitInfo.Unit.MaxHealth)
				positionStr := fmt.Sprintf("(%d,%d)", unitInfo.Row, unitInfo.Col)

				// Get city name if this is a city unit
				cityName := ""
				if unitInfo.Unit.UnitType == UnitTypeCity {
					cityName = GetCityNameAtPosition(unitInfo.Row, unitInfo.Col, saveOutput)
				}

				unitData = append(unitData, UnitTableData{
					UnitID:   unitInfo.Index,
					UnitType: unitTypeName,
					Level:    int(unitInfo.Unit.Level),
					Health:   healthStr,
					Position: positionStr,
					City:     cityName,
				})
			}

			tableFormatter := NewTableFormatter()
			tableFormatter.PrintUnitTable(unitData)
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

			// Prepare data for table
			var generalData []GeneralTableData
			for _, generalInfo := range generalUnits {
				unitTypeName := GetUnitTypeName(generalInfo.Unit.UnitType)
				generalName, hasGeneralName := GetGeneralName(generalInfo.Unit.GeneralId)
				healthStr := fmt.Sprintf("%d/%d", generalInfo.Unit.CurrentHealth, generalInfo.Unit.MaxHealth)
				positionStr := fmt.Sprintf("(%d,%d)", generalInfo.Row, generalInfo.Col)

				generalInfoStr := fmt.Sprintf("ID %d", generalInfo.Unit.GeneralId)
				if hasGeneralName {
					generalInfoStr = generalName
				}

				generalData = append(generalData, GeneralTableData{
					UnitID:      generalInfo.Index,
					UnitType:    unitTypeName,
					GeneralName: generalInfoStr,
					Level:       int(generalInfo.Unit.Level),
					Health:      healthStr,
					Position:    positionStr,
				})
			}

			tableFormatter := NewTableFormatter()
			tableFormatter.PrintGeneralTable(generalData)
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

// ListLandmines displays all landmines grouped by owner
func ListLandmines(saveOutput *WC4SaveOutput) {
	DisplayHeader(fmt.Sprintf("Landmines (Total: %d)", len(saveOutput.Landmines)))

	if len(saveOutput.Landmines) == 0 {
		fmt.Println("No landmines found.")
		return
	}

	// Process landmines and group by owner
	landminesByOwner := make(map[byte][]LandmineDisplayInfo)
	skippedLandmines := 0

	for i, landmine := range saveOutput.Landmines {
		row, col, valid := ConvertCoordinates(int(landmine.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			fmt.Printf("WARNING: Landmine %d: skipping invalid coordinates (CoordinateCode=%d)\n", i, landmine.CoordinateCode)
			skippedLandmines++
			continue
		}

		// Check if owner is valid
		if int(landmine.Owner) >= len(saveOutput.PlayerData) {
			fmt.Printf("WARNING: Landmine %d: invalid owner %d, skipping\n", i, landmine.Owner)
			skippedLandmines++
			continue
		}

		landminesByOwner[byte(landmine.Owner)] = append(landminesByOwner[byte(landmine.Owner)], LandmineDisplayInfo{
			Index: i,
			Landmine: landmine,
			Row:   row,
			Col:   col,
		})
	}

	// Show skip statistics
	if skippedLandmines > 0 {
		fmt.Printf("Skipped %d landmines due to invalid data\n", skippedLandmines)
	}

	// Display landmines grouped by owner
	tableFormatter := NewTableFormatter()
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		if landmines, exists := landminesByOwner[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("\nPlayer %d (%s - CountryId %d) owns %d landmines:\n",
				ownerID, countryName, player.CountryId, len(landmines))

			// Prepare data for table
			var landmineData []LandmineTableData
			for _, landmineInfo := range landmines {
				positionStr := fmt.Sprintf("(%d,%d)", landmineInfo.Row, landmineInfo.Col)
				ownerName := fmt.Sprintf("Player %d (%s)", ownerID, countryName)
				healthStr := fmt.Sprintf("%d", landmineInfo.Landmine.Health)

				landmineData = append(landmineData, LandmineTableData{
					LandmineID: landmineInfo.Index,
					Position:   positionStr,
					Owner:      ownerName,
					Health:     healthStr,
				})
			}

			tableFormatter.PrintLandmineTable(landmineData)
		}
	}

	// Add landmine summary
	fmt.Println("\nLandmine Summary:")
	fmt.Println("----------------")
	totalValidLandmines := len(saveOutput.Landmines) - skippedLandmines
	fmt.Printf("Total landmines: %d\n", len(saveOutput.Landmines))
	fmt.Printf("Valid landmines: %d\n", totalValidLandmines)
	fmt.Printf("Skipped landmines: %d\n", skippedLandmines)
}

// LandmineDisplayInfo represents a landmine with its display information
type LandmineDisplayInfo struct {
	Index    int
	Landmine LandmineData
	Row      int
	Col      int
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

	tableFormatter := NewTableFormatter()
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	
	// Prepare player summary data
	var playerSummaryData []PlayerSummaryData
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
			playerSummaryData = append(playerSummaryData, PlayerSummaryData{
				PlayerID:      ownerID,
				CountryName:   countryName,
				CountryID:     int(player.CountryId),
				TeamID:        int(player.TeamId),
				CityCount:     len(cities),
				TotalTiles:    totalTilesForPlayer,
				LandPercent:   landPercentage,
				WorldPercent:  worldPercentage,
			})
		}
	}

	// Print player summary table
	tableFormatter.PrintPlayerSummaryTable(playerSummaryData)

	// Print detailed city tables for each player
	for _, ownerID := range sortedPlayerIDs {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			countryName, _ := GetCountryInfo(player.CountryId)

			// Calculate total tiles for this player
			totalTilesForPlayer := 0
			for _, cityInfo := range cities {
				tileCount := result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				totalTilesForPlayer += tileCount
			}

			landPercentage := float64(totalTilesForPlayer) / float64(result.LandTiles) * 100
			worldPercentage := float64(totalTilesForPlayer) / float64(result.TotalTiles) * 100

			// Show detailed player summary with better formatting
			fmt.Printf("\nPlayer %d: %s\n", ownerID, countryName)
			fmt.Printf("  Country ID: %d, Team ID: %d\n", player.CountryId, player.TeamId)
			fmt.Printf("  City Tiles: %d cities, %d total tiles\n", len(cities), totalTilesForPlayer)
			fmt.Printf("  Territory: %.1f%% of land, %.1f%% of world\n", landPercentage, worldPercentage)
			
			// Prepare city data for table
			var cityData []CityTileData
			for _, cityInfo := range cities {
				cityName, hasName := GetCityName(cityInfo.City.CityId)
				if !hasName {
					cityName = "Unknown"
				}
				tileCount := result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
				
				cityData = append(cityData, CityTileData{
					CityIndex:    cityInfo.Index,
					CityName:     cityName,
					CityID:       int(cityInfo.City.CityId),
					Position:     fmt.Sprintf("(%d,%d)", cityInfo.Row, cityInfo.Col),
					CoordinateCode: int(cityInfo.City.CoordinateCode),
					TileCount:    tileCount,
				})
			}
			
			tableFormatter.PrintCityTileTable(cityData)
		}
	}
}

// DisplayTerritoryControl shows player territory control summary
func DisplayTerritoryControl(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	fmt.Println("\nPlayer Territory Control:")
	fmt.Println("------------------------")
	
	tableFormatter := NewTableFormatter()
	totalOwnedTiles := 0
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)

	// Prepare territory control data
	var territoryData []TerritoryControlData
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
			
			territoryData = append(territoryData, TerritoryControlData{
				PlayerID:      ownerID,
				CountryName:   countryName,
				CountryID:     int(player.CountryId),
				TileCount:     playerTiles,
				LandPercent:   landPercentage,
				WorldPercent:  worldPercentage,
			})
		}
	}

	// Print individual player territory control table
	tableFormatter.PrintTerritoryControlTable(territoryData)

	// Print team territory control summary
	fmt.Println("\nTeam Territory Control:")
	fmt.Println("----------------------")
	DisplayTeamTerritoryControl(saveOutput, citiesByOwner, result)

	// Print total summary
	landPercentage := float64(totalOwnedTiles) / float64(result.LandTiles) * 100
	worldPercentage := float64(totalOwnedTiles) / float64(result.TotalTiles) * 100
	fmt.Printf("\nTotal controlled: %d tiles (%.1f%% of land, %.1f%% of world)\n",
		totalOwnedTiles, landPercentage, worldPercentage)
}

// DisplayTeamTerritoryControl shows team territory control summary
func DisplayTeamTerritoryControl(saveOutput *WC4SaveOutput, citiesByOwner map[byte][]CityDisplayInfo, result TileAnalysisResult) {
	tableFormatter := NewTableFormatter()
	
	// Group players by team and calculate team totals
	teamData := make(map[uint32]*TeamTerritoryData)
	
	for _, ownerID := range GetSortedPlayerIDs(saveOutput.PlayerData) {
		if cities, exists := citiesByOwner[byte(ownerID)]; exists {
			player := saveOutput.PlayerData[ownerID]
			teamId := player.TeamId
			
			// Calculate player tiles
			playerTiles := 0
			for _, cityInfo := range cities {
				playerTiles += result.CoordinateCodeCounts[cityInfo.City.CoordinateCode]
			}
			
			// Initialize team data if not exists
			if teamData[teamId] == nil {
				teamData[teamId] = &TeamTerritoryData{
					TeamID:      int(teamId),
					PlayerCount: 0,
					TileCount:   0,
					Players:     []string{},
				}
			}
			
			// Add player data to team
			countryName, _ := GetCountryInfo(player.CountryId)
			teamData[teamId].PlayerCount++
			teamData[teamId].TileCount += playerTiles
			teamData[teamId].Players = append(teamData[teamId].Players, fmt.Sprintf("%d (%s)", ownerID, countryName))
		}
	}
	
	// Convert to slice and calculate percentages
	var teamTerritoryData []TeamTerritoryData
	for _, team := range teamData {
		landPercentage := float64(team.TileCount) / float64(result.LandTiles) * 100
		worldPercentage := float64(team.TileCount) / float64(result.TotalTiles) * 100
		team.LandPercent = landPercentage
		team.WorldPercent = worldPercentage
		teamTerritoryData = append(teamTerritoryData, *team)
	}
	
	// Sort teams by tile count (descending) for better readability
	sort.Slice(teamTerritoryData, func(i, j int) bool {
		return teamTerritoryData[i].TileCount > teamTerritoryData[j].TileCount
	})
	
	// Print team territory control table
	tableFormatter.PrintTeamTerritoryControlTable(teamTerritoryData)
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
