package fileio

import (
	"fmt"
	"log"
	"os"
)

// Tile ownership constants
const (
	TileUnowned = 255 // Value indicating unowned/neutral tile
	MainPlayer  = 0   // Main player ID (player 0)
)

// ConversionStats holds aggregated statistics for conversion operations
type ConversionStats struct {
	TotalChanged    int          // Total number of tiles/units changed
	ChangesByPlayer map[byte]int // Changes grouped by original player
}

// NewConversionStats creates a new ConversionStats struct with initialized maps
func NewConversionStats() *ConversionStats {
	return &ConversionStats{
		ChangesByPlayer: make(map[byte]int),
	}
}

// AddChange records a change in the statistics
func (cs *ConversionStats) AddChange(oldPlayer, newPlayer byte, x, y int) {
	cs.TotalChanged++
	cs.ChangesByPlayer[oldPlayer]++
}

// PrintSummary prints a formatted summary of the conversion statistics
func (cs *ConversionStats) PrintSummary(operation string, saveOutput *WC4SaveOutput) {
	fmt.Printf("\n=== %s Summary ===\n", operation)
	fmt.Printf("Total changes: %d\n", cs.TotalChanged)

	if len(cs.ChangesByPlayer) > 0 {
		fmt.Println("\nChanges by original player:")
		// Get sorted player IDs for consistent ordering
		sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
		for _, playerID := range sortedPlayerIDs {
			if count, exists := cs.ChangesByPlayer[byte(playerID)]; exists {
				player := saveOutput.PlayerData[playerID]
				countryName, _ := GetCountryInfo(player.CountryId)
				fmt.Printf("  Player %d (%s - CountryId %d): %d changes\n", playerID, countryName, player.CountryId, count)
			}
		}
	}

	fmt.Println("========================\n")
}

func WriteUnitOwnerToFile(inputFilename string, value int, targetX int, targetY int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	offsetStartOriginalBlockKey := buildSingleUnitOwnerStartKey(targetX, targetY)
	offset, ok := fileOffsetMap[offsetStartOriginalBlockKey]
	if !ok {
		log.Fatal(fmt.Sprintf("Error: Unable to find start of data block with key %v. Command not run.", offsetStartOriginalBlockKey))
	}

	if _, err := inputFile.WriteAt([]byte{uint8(value)}, int64(offset)); err != nil {
		log.Fatal("Failed to write uint8 to file:", err)
	}
}

func WriteAllUnitOwnersToFile(inputFilename string, tileDataOverwrite [][]byte) {
	byteData := make([]byte, 0)
	for i := 0; i < len(tileDataOverwrite); i++ {
		byteData = append(byteData, tileDataOverwrite[i]...)
	}

	WriteAndShiftData(inputFilename, buildUnitOwnerStartKey(), buildUnitOwnerEndKey(), byteData)
}

func SetPlayerMaxCurrency(filename string, playerID int, currencyValue int) {
	offset := fileOffsetMap[BuildPlayerStartKey(playerID)]
	fmt.Printf("Player %d base offset: 0x%04X\n", playerID, offset)

	// Currency values are at offsets +8, +12, +16 from player start
	currencyOffsets := []int{8, 12, 16}
	currencyNames := []string{"Currency 1", "Currency 2", "Currency 3"}

	fmt.Println("Current currency values:")
	for i, currencyOffset := range currencyOffsets {
		fullOffset := offset + currencyOffset
		oldValue := ReadUint16AtFileOffset(filename, fullOffset)
		fmt.Printf("  %s: %d at offset 0x%04X\n", currencyNames[i], oldValue, fullOffset)
	}

	fmt.Println("Setting new currency values:")
	for i, currencyOffset := range currencyOffsets {
		fullOffset := offset + currencyOffset
		oldValue := ReadUint16AtFileOffset(filename, fullOffset)
		WriteUint16AtFileOffset(filename, fullOffset, currencyValue)
		fmt.Printf("  %s: %d -> %d at offset 0x%04X\n", currencyNames[i], oldValue, currencyValue, fullOffset)
	}

	fmt.Printf("Set max currency to %d for player %d\n", currencyValue, playerID)
}

func SetPlayerMaxCityTech(filename string, saveOutput *WC4SaveOutput, playerID int, techLevel int) {
	fmt.Printf("Setting max city tech to level %d for player %d\n", techLevel, playerID)
	fmt.Println("Processing cities...")

	citiesModified := 0
	for i := 0; i < len(saveOutput.Cities); i++ {
		city := saveOutput.Cities[i]
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			fmt.Printf("WARNING: City %d: skipping invalid coordinates (CoordinateCode=%d)\n", i, city.CoordinateCode)
			continue
		}

		owner := saveOutput.UnitOwnerData[row][col]
		if owner != uint8(playerID) {
			continue
		}

		offset := fileOffsetMap[BuildCityStartKey(i)]
		cityName, hasName := GetCityName(city.CityId)
		if !hasName {
			cityName = "Unknown"
		}
		fmt.Printf("City %d: %s (ID:0x%02x) base offset: 0x%04X (owner: %d)\n", i, cityName, city.CityId, offset, owner)

		// Tech levels are at offsets +24 to +29 (6 bytes)
		techOffsets := []int{24, 25, 26, 27, 28, 29}
		techNames := []string{"Tech 1", "Tech 2", "Tech 3", "Tech 4", "Tech 5", "Tech 6"}

		fmt.Println("  Current tech levels:")
		for j, techOffset := range techOffsets {
			fullOffset := offset + techOffset
			oldValue := ReadUint8AtFileOffset(filename, fullOffset)
			fmt.Printf("    %s: %d at offset 0x%04X\n", techNames[j], oldValue, fullOffset)
		}

		fmt.Println("  Setting new tech levels:")
		for j, techOffset := range techOffsets {
			fullOffset := offset + techOffset
			oldValue := ReadUint8AtFileOffset(filename, fullOffset)
			WriteUint8AtFileOffset(filename, fullOffset, techLevel)
			fmt.Printf("    %s: %d -> %d at offset 0x%04X\n", techNames[j], oldValue, techLevel, fullOffset)
		}

		citiesModified++
	}

	fmt.Printf("Set max city tech to level %d for %d cities owned by player %d\n", techLevel, citiesModified, playerID)
}

func RestoreAlliesHealth(filename string, saveOutput *WC4SaveOutput, playerID int) {
	playerTeamId := saveOutput.PlayerData[playerID].TeamId
	fmt.Printf("Restoring health for allies of player %d (TeamId: %d)\n", playerID, playerTeamId)
	fmt.Println("Processing units...")

	unitsRestored := 0
	unitsByPlayer := make(map[byte][]UnitDisplayInfo) // Track unit info by player

	// Process all units using common function
	result := ProcessAllUnits(saveOutput)
	validUnits := result.ValidUnits

	for _, unitInfo := range validUnits {
		owner := unitInfo.Owner
		if int(owner) >= len(saveOutput.PlayerData) {
			fmt.Printf("WARNING: Invalid owner %d, skipping unit %d\n", owner, unitInfo.Index)
			continue
		}
		if saveOutput.PlayerData[owner].TeamId != playerTeamId {
			continue
		}

		offset := fileOffsetMap[BuildUnitHealthKey(unitInfo.Index)]
		oldHealth := ReadUint16AtFileOffset(filename, offset)

		fmt.Printf("Unit %d: Restoring health %d -> %d at offset 0x%04X (owner: %d)\n",
			unitInfo.Index, oldHealth, unitInfo.Unit.MaxHealth, offset, owner)

		WriteUint16AtFileOffset(filename, offset, int(unitInfo.Unit.MaxHealth))
		unitsByPlayer[owner] = append(unitsByPlayer[owner], unitInfo)
		unitsRestored++
	}

	// Print breakdown by player
	fmt.Println("\nBreakdown by Player:")
	fmt.Println("-------------------")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		unitInfos, exists := unitsByPlayer[byte(ownerID)]
		if !exists || len(unitInfos) == 0 {
			continue
		}
		skippedCount := result.SkippedUnits[byte(ownerID)]
		fmt.Printf("Player %d: %d units restored", ownerID, len(unitInfos))
		if skippedCount > 0 {
			fmt.Printf(" (skipped %d invalid units)", skippedCount)
		}
		fmt.Println()
		for _, unitInfo := range unitInfos {
			unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
			cityInfo := ""
			if unitInfo.Unit.UnitType == UnitTypeCity {
				// Try to find the city name by looking for a city at the same position
				cityName := GetCityNameAtPosition(unitInfo.Row, unitInfo.Col, saveOutput)
				cityInfo = fmt.Sprintf(" (%s)", cityName)
			}
			fmt.Printf("  Unit %d: Type=%d (%s)%s, Level=%d, Health=%d->%d, Position=(%d,%d)\n",
				unitInfo.Index, unitInfo.Unit.UnitType, unitTypeName, cityInfo, unitInfo.Unit.Level, unitInfo.Unit.CurrentHealth, unitInfo.Unit.MaxHealth, unitInfo.Row, unitInfo.Col)
		}
	}

	fmt.Printf("\nRestored allies health. Changed %d units to have max health.\n", unitsRestored)
}

func WeakenEnemies(filename string, saveOutput *WC4SaveOutput, playerID int) {
	playerTeamId := saveOutput.PlayerData[playerID].TeamId
	fmt.Printf("Weakening enemies of player %d (TeamId: %d)\n", playerID, playerTeamId)
	fmt.Println("Processing units...")

	unitsWeakened := 0
	unitsByPlayer := make(map[byte][]UnitDisplayInfo) // Track unit info by player
	unitsByType := make(map[uint8]int)                // Track unit types

	// Process all units using common function
	result := ProcessAllUnits(saveOutput)
	validUnits := result.ValidUnits

	for _, unitInfo := range validUnits {
		owner := unitInfo.Owner
		if int(owner) >= len(saveOutput.PlayerData) {
			fmt.Printf("WARNING: Invalid owner %d, skipping unit %d\n", owner, unitInfo.Index)
			continue
		}

		if saveOutput.PlayerData[owner].TeamId == playerTeamId {
			continue // Skip allies
		}

		offset := fileOffsetMap[BuildUnitHealthKey(unitInfo.Index)]
		oldHealth := ReadUint16AtFileOffset(filename, offset)

		var newHealth int
		var unitTypeName string
		if unitInfo.Unit.UnitType == UnitTypeCity {
			newHealth = 0
			unitTypeName = "city"
		} else {
			newHealth = 1
			unitTypeName = "unit"
		}

		fmt.Printf("Enemy %s %d: Reducing health %d -> %d at offset 0x%04X (owner: %d)\n",
			unitTypeName, unitInfo.Index, oldHealth, newHealth, offset, owner)

		WriteUint16AtFileOffset(filename, offset, newHealth)
		unitsByPlayer[owner] = append(unitsByPlayer[owner], unitInfo)
		unitsByType[unitInfo.Unit.UnitType]++
		unitsWeakened++
	}

	// Print breakdown by player
	fmt.Println("\nBreakdown by Player:")
	fmt.Println("-------------------")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, ownerID := range sortedPlayerIDs {
		unitInfos, exists := unitsByPlayer[byte(ownerID)]
		if !exists || len(unitInfos) == 0 {
			continue
		}
		skippedCount := result.SkippedUnits[byte(ownerID)]
		fmt.Printf("Player %d: %d units weakened", ownerID, len(unitInfos))
		if skippedCount > 0 {
			fmt.Printf(" (skipped %d invalid units)", skippedCount)
		}
		fmt.Println()
		for _, unitInfo := range unitInfos {
			unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
			cityInfo := ""
			if unitInfo.Unit.UnitType == UnitTypeCity {
				// Try to find the city name by looking for a city at the same position
				cityName := GetCityNameAtPosition(unitInfo.Row, unitInfo.Col, saveOutput)
				cityInfo = fmt.Sprintf(" (%s)", cityName)
			}
			var newHealth int
			if unitInfo.Unit.UnitType == UnitTypeCity {
				newHealth = 0
			} else {
				newHealth = 1
			}
			fmt.Printf("  Unit %d: Type=%d (%s)%s, Level=%d, Health=%d->%d, Position=(%d,%d)\n",
				unitInfo.Index, unitInfo.Unit.UnitType, unitTypeName, cityInfo, unitInfo.Unit.Level, unitInfo.Unit.CurrentHealth, newHealth, unitInfo.Row, unitInfo.Col)
		}
	}

	// Print breakdown by unit type
	fmt.Println("\nBreakdown by Unit Type:")
	fmt.Println("----------------------")

	// Sort unit types in numeric ascending order
	sortedUnitTypes := GetSortedUnitTypes(unitsByType)

	for _, unitType := range sortedUnitTypes {
		count := unitsByType[byte(unitType)]
		unitTypeName := GetUnitTypeName(uint8(unitType))
		fmt.Printf("Type %d (%s): %d units\n", unitType, unitTypeName, count)
	}

	fmt.Printf("\nWeakened enemies. Changed %d units to have low health.\n", unitsWeakened)
}

func ConvertCoordinates(coordinateCode int, unitOwnerData [][]byte, gameMode int) (int, int, bool) {
	row := int(coordinateCode) / len(unitOwnerData[0])
	if gameMode == GameModeConquest { // subtract 2 from row if conquest
		row -= 2
	}

	col := int(coordinateCode) % len(unitOwnerData[0])

	// Check if coordinates are within valid bounds
	if row < 0 || row >= len(unitOwnerData) || col < 0 || col >= len(unitOwnerData[0]) {
		return row, col, false // Invalid coordinates
	}

	return row, col, true // Valid coordinates
}

func ConvertToCoordinateCode(row, col int, unitOwnerData [][]byte, gameMode int) int {
	coordinateCode := row*len(unitOwnerData[0]) + col
	if gameMode == GameModeConquest { // conquest mode
		coordinateCode += 2 * len(unitOwnerData[0])
	}
	return coordinateCode
}

func GetCityNameAtPosition(row, col int, saveOutput *WC4SaveOutput) string {
	for _, city := range saveOutput.Cities {
		cityRow, cityCol, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if valid && cityRow == row && cityCol == col {
			if name, hasName := GetCityName(city.CityId); hasName {
				return name
			}
			break
		}
	}
	return "Unknown"
}

func ConvertCoordinatesWithDebug(coordinateCode int, unitOwnerData [][]byte, gameMode int) (int, int, bool) {
	row, col, valid := ConvertCoordinates(coordinateCode, unitOwnerData, gameMode)
	if valid {
		fmt.Printf("Coordinate %v (row: %v, column: %v)\n", coordinateCode, row, col)
	}
	return row, col, valid
}

// ConvertPlayer converts all units owned by oldPlayer to newPlayer
func ConvertPlayer(inputFilename string, saveOutput *WC4SaveOutput, oldPlayer, newPlayer int) *ConversionStats {
	stats := NewConversionStats()

	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] != TileUnowned && saveOutput.UnitOwnerData[i][j] != MainPlayer {
				if saveOutput.UnitOwnerData[i][j] != byte(oldPlayer) {
					continue
				}

				oldValue := saveOutput.UnitOwnerData[i][j]
				saveOutput.UnitOwnerData[i][j] = byte(newPlayer)
				fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %v", i, j, oldPlayer, newPlayer))
				stats.AddChange(oldValue, byte(newPlayer), j, i)
			}
		}
	}
	WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
	return stats
}

// ConvertTile converts a specific tile at (x, y) to newPlayer
func ConvertTile(inputFilename string, saveOutput *WC4SaveOutput, x, y, newPlayer int) error {
	oldPlayer := saveOutput.UnitOwnerData[y][x]
	if oldPlayer == TileUnowned {
		return fmt.Errorf("can't convert tile at (%v, %v) - tile has no owner", y, x)
	}
	WriteUnitOwnerToFile(inputFilename, newPlayer, x, y)
	fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %v", y, x, oldPlayer, newPlayer))
	return nil
}

// ConvertAllAllies converts all allied units to player 0
func ConvertAllAllies(inputFilename string, saveOutput *WC4SaveOutput) *ConversionStats {
	stats := NewConversionStats()
	playerTeamId := saveOutput.PlayerData[MainPlayer].TeamId

	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] == TileUnowned {
				continue
			}

			if saveOutput.UnitOwnerData[i][j] == MainPlayer {
				continue
			}

			oldValue := saveOutput.UnitOwnerData[i][j]
			if saveOutput.PlayerData[oldValue].TeamId == playerTeamId {
				saveOutput.UnitOwnerData[i][j] = MainPlayer
				fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %d", i, j, oldValue, MainPlayer))
				stats.AddChange(oldValue, MainPlayer, j, i)
			}
		}
	}
	WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
	return stats
}

// ConvertTeam converts all players to the same team as player 0
func ConvertTeam(inputFilename string, saveOutput *WC4SaveOutput) {
	playerTeamId := saveOutput.PlayerData[MainPlayer].TeamId
	for i := 1; i < len(saveOutput.PlayerData); i++ {
		offset := GetFileOffsetMap()[BuildPlayerStartKey(i)]
		WriteUint32AtFileOffset(inputFilename, offset+24, int(playerTeamId))
		fmt.Println("Converting player", i, "from team", saveOutput.PlayerData[i].TeamId, "to team", playerTeamId)
	}
}

// ConvertAllPlayers converts all units to player 0
func ConvertAllPlayers(inputFilename string, saveOutput *WC4SaveOutput) *ConversionStats {
	stats := NewConversionStats()

	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] != TileUnowned && saveOutput.UnitOwnerData[i][j] != MainPlayer {
				oldValue := saveOutput.UnitOwnerData[i][j]
				saveOutput.UnitOwnerData[i][j] = MainPlayer
				fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %d", i, j, oldValue, MainPlayer))
				stats.AddChange(oldValue, MainPlayer, j, i)
			}
		}
	}
	WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
	return stats
}
