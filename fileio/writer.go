package fileio

import (
	"encoding/binary"
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

	WriteDataAtOffset(inputFilename, GetOffset(buildUnitOwnerStartKey()), byteData)
}

func SetPlayerMaxCurrency(filename string, playerID int, currencyValue int) {
	// Input validation
	if playerID < 0 {
		log.Fatal("Error: Player ID cannot be negative")
	}
	if currencyValue < 0 || currencyValue > 65535 {
		log.Fatal(fmt.Sprintf("Error: Currency value %d is invalid. Must be between 0 and 65535", currencyValue))
	}

	offset := GetOffset(BuildPlayerStartKey(playerID))
	fmt.Printf("Player %d base offset: 0x%04X\n", playerID, offset)

	// Currency [3]uint32 starts at offset +8 and is 12 bytes total (3 * 4 bytes)
	currencyStartOffset := offset + 8

	fmt.Println("Current currency values:")
	for i := 0; i < 3; i++ {
		currencyOffset := currencyStartOffset + (i * 4) // Each uint32 is 4 bytes
		oldValue := ReadUint32AtFileOffset(filename, currencyOffset)
		fmt.Printf("  Currency %d: %d at offset 0x%04X\n", i+1, oldValue, currencyOffset)
	}

	// Write all three currency values as a contiguous block
	currencyData := make([]byte, 12)                                         // 3 * uint32 = 12 bytes
	binary.LittleEndian.PutUint32(currencyData[0:4], uint32(currencyValue))  // Currency[0]
	binary.LittleEndian.PutUint32(currencyData[4:8], uint32(currencyValue))  // Currency[1]
	binary.LittleEndian.PutUint32(currencyData[8:12], uint32(currencyValue)) // Currency[2]

	// Use WriteDataAtOffset with calculated offset
	WriteDataAtOffset(filename, currencyStartOffset, currencyData)

	fmt.Println("Setting new currency values:")
	for i := 0; i < 3; i++ {
		currencyOffset := currencyStartOffset + (i * 4)
		fmt.Printf("  Currency %d: -> %d at offset 0x%04X\n", i+1, currencyValue, currencyOffset)
	}

	fmt.Printf("Set max currency to %d for player %d\n", currencyValue, playerID)
}

func SetPlayerMaxCityTech(filename string, saveOutput *WC4SaveOutput, playerID int, techLevel int) {
	// Input validation
	if playerID < 0 {
		log.Fatal("Error: Player ID cannot be negative")
	}
	if techLevel < 0 || techLevel > 255 {
		log.Fatal(fmt.Sprintf("Error: Tech level %d is invalid. Must be between 0 and 255", techLevel))
	}
	if playerID >= len(saveOutput.PlayerData) {
		log.Fatal(fmt.Sprintf("Error: Player ID %d does not exist. Valid range: 0-%d", playerID, len(saveOutput.PlayerData)-1))
	}

	fmt.Printf("Setting max city tech to level %d for player %d\n", techLevel, playerID)
	fmt.Println("Processing cities...")

	// Process cities owned by the player
	citiesModified := processPlayerCityTech(filename, saveOutput, playerID, techLevel)

	fmt.Printf("Set max city tech to level %d for %d cities owned by player %d\n", techLevel, citiesModified, playerID)
}

// processPlayerCityTech processes cities owned by the specified player
func processPlayerCityTech(filename string, saveOutput *WC4SaveOutput, playerID int, techLevel int) int {
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

		offset := GetOffset(BuildCityStartKey(i))
		cityName, hasName := GetCityName(city.CityId)
		if !hasName {
			cityName = "Unknown"
		}

		fmt.Printf("City %d: %s (ID:0x%02x) base offset: 0x%04X (owner: %d)\n", i, cityName, city.CityId, offset, owner)

		// TechLevels [6]byte starts at offset +24
		techStartOffset := offset + 24
		updateCityTechLevels(filename, techStartOffset, techLevel)

		citiesModified++
	}

	return citiesModified
}

// processEnemyCityTech processes cities owned by enemies
func processEnemyCityTech(filename string, saveOutput *WC4SaveOutput, playerID int, playerTeamId uint32, techLevel int) int {
	citiesModified := 0

	for i := 0; i < len(saveOutput.Cities); i++ {
		city := saveOutput.Cities[i]
		row, col, valid := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
		if !valid {
			continue
		}

		owner := saveOutput.UnitOwnerData[row][col]
		if int(owner) >= len(saveOutput.PlayerData) {
			continue
		}

		// Skip if owner is main player or ally
		if int(owner) == playerID || saveOutput.PlayerData[owner].TeamId == playerTeamId {
			continue
		}

		offset := GetOffset(BuildCityStartKey(i))
		cityName, hasName := GetCityName(city.CityId)
		if !hasName {
			cityName = "Unknown"
		}

		player := saveOutput.PlayerData[owner]
		countryName, _ := GetCountryInfo(player.CountryId)
		fmt.Printf("City %d: %s (owner: Player %d - %s) - Minimizing tech\n", i, cityName, owner, countryName)

		// TechLevels [6]byte starts at offset +24
		techStartOffset := offset + 24
		updateCityTechLevels(filename, techStartOffset, techLevel)

		citiesModified++
	}

	return citiesModified
}

// updateCityTechLevels updates the tech levels for a city at the given offset
func updateCityTechLevels(filename string, techStartOffset int, techLevel int) {
	fmt.Println("  Current tech levels:")
	for j := 0; j < 6; j++ {
		techOffset := techStartOffset + j
		oldValue := ReadUint8AtFileOffset(filename, techOffset)
		fmt.Printf("    Tech %d: %d at offset 0x%04X\n", j+1, oldValue, techOffset)
	}

	// Write all 6 tech levels as a contiguous block
	techData := make([]byte, 6)
	for j := 0; j < 6; j++ {
		techData[j] = byte(techLevel)
	}

	WriteDataAtOffset(filename, techStartOffset, techData)

	fmt.Println("  Setting new tech levels:")
	for j := 0; j < 6; j++ {
		techOffset := techStartOffset + j
		fmt.Printf("    Tech %d: -> %d at offset 0x%04X\n", j+1, techLevel, techOffset)
	}
}

func RestoreAlliesHealth(filename string, saveOutput *WC4SaveOutput, playerID int) {
	// Input validation
	if playerID < 0 {
		log.Fatal("Error: Player ID cannot be negative")
	}
	if playerID >= len(saveOutput.PlayerData) {
		log.Fatal(fmt.Sprintf("Error: Player ID %d does not exist. Valid range: 0-%d", playerID, len(saveOutput.PlayerData)-1))
	}

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

		offset := fileOffsetMap[BuildUnitStartKey(unitInfo.Index)] + 12
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
	tableFormatter := NewTableFormatter()
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

		// Prepare data for table
		var healAlliesData []HealAlliesUnitData
		for _, unitInfo := range unitInfos {
			unitTypeName := GetUnitTypeName(unitInfo.Unit.UnitType)
			unitTypeStr := fmt.Sprintf("%d (%s)", unitInfo.Unit.UnitType, unitTypeName)

			// Add city name in parentheses if it's a city unit
			if unitInfo.Unit.UnitType == UnitTypeCity {
				cityName := GetCityNameAtPosition(unitInfo.Row, unitInfo.Col, saveOutput)
				if cityName != "" {
					unitTypeStr = fmt.Sprintf("%d (%s) - [%s]", unitInfo.Unit.UnitType, unitTypeName, cityName)
				}
			}

			healAlliesData = append(healAlliesData, HealAlliesUnitData{
				UnitID:    unitInfo.Index,
				UnitType:  unitTypeStr,
				Level:     int(unitInfo.Unit.Level),
				OldHealth: int(unitInfo.Unit.CurrentHealth),
				NewHealth: int(unitInfo.Unit.MaxHealth),
				Position:  fmt.Sprintf("(%d,%d)", unitInfo.Row, unitInfo.Col),
			})
		}

		tableFormatter.PrintHealAlliesTable(healAlliesData)
		fmt.Println()
	}

	fmt.Printf("\nRestored allies health. Changed %d units to have max health.\n", unitsRestored)
}

func WeakenEnemies(filename string, saveOutput *WC4SaveOutput, playerID int) {
	// Input validation
	if playerID < 0 {
		log.Fatal("Error: Player ID cannot be negative")
	}
	if playerID >= len(saveOutput.PlayerData) {
		log.Fatal(fmt.Sprintf("Error: Player ID %d does not exist. Valid range: 0-%d", playerID, len(saveOutput.PlayerData)-1))
	}

	playerTeamId := saveOutput.PlayerData[playerID].TeamId
	fmt.Printf("Weakening enemies of player %d (TeamId: %d)\n", playerID, playerTeamId)

	// First, minimize enemy money
	fmt.Println("Minimizing enemy money...")
	enemiesWeakened := 0
	for i := 0; i < len(saveOutput.PlayerData); i++ {
		if i == playerID || saveOutput.PlayerData[i].TeamId == playerTeamId {
			continue // Skip main player and allies
		}

		player := saveOutput.PlayerData[i]
		countryName, _ := GetCountryInfo(player.CountryId)
		fmt.Printf("Player %d (%s): Minimizing currency\n", i, countryName)

		// Use the same logic as SetPlayerMaxCurrency but with value 0
		SetPlayerMaxCurrency(filename, i, 0)
		enemiesWeakened++
	}

	fmt.Printf("Minimized money for %d enemy players\n", enemiesWeakened)

	// Then, minimize enemy city tech levels
	fmt.Println("\nMinimizing enemy city tech levels...")
	citiesWeakened := processEnemyCityTech(filename, saveOutput, playerID, playerTeamId, 0)

	fmt.Printf("Minimized tech for %d enemy cities\n", citiesWeakened)

	// Finally, weaken enemy units
	fmt.Println("\nProcessing units...")
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

		offset := fileOffsetMap[BuildUnitStartKey(unitInfo.Index)] + 12
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

	fmt.Printf("\nWeakened enemies. Minimized money for %d players, tech for %d cities, and reduced %d units to low health.\n", enemiesWeakened, citiesWeakened, unitsWeakened)
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
	// Input validation
	if oldPlayer < 0 || newPlayer < 0 {
		log.Fatal("Error: Player IDs cannot be negative")
	}
	if oldPlayer >= len(saveOutput.PlayerData) {
		log.Fatal(fmt.Sprintf("Error: Old player ID %d does not exist. Valid range: 0-%d", oldPlayer, len(saveOutput.PlayerData)-1))
	}
	if newPlayer >= len(saveOutput.PlayerData) {
		log.Fatal(fmt.Sprintf("Error: New player ID %d does not exist. Valid range: 0-%d", newPlayer, len(saveOutput.PlayerData)-1))
	}
	if oldPlayer == newPlayer {
		log.Fatal("Error: Old player and new player cannot be the same")
	}

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
	// Input validation
	if x < 0 || y < 0 {
		return fmt.Errorf("coordinates cannot be negative: (%d, %d)", x, y)
	}
	if y >= len(saveOutput.UnitOwnerData) || x >= len(saveOutput.UnitOwnerData[0]) {
		return fmt.Errorf("coordinates out of bounds: (%d, %d). Map size: %dx%d", x, y, len(saveOutput.UnitOwnerData[0]), len(saveOutput.UnitOwnerData))
	}
	if newPlayer < 0 {
		return fmt.Errorf("new player ID cannot be negative: %d", newPlayer)
	}
	if newPlayer >= len(saveOutput.PlayerData) {
		return fmt.Errorf("new player ID %d does not exist. Valid range: 0-%d", newPlayer, len(saveOutput.PlayerData)-1)
	}

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
	convertedCount := 0

	fmt.Println("Converting all players to team", playerTeamId)
	fmt.Println("----------------------------------------")

	for i := 1; i < len(saveOutput.PlayerData); i++ {
		oldTeamId := saveOutput.PlayerData[i].TeamId
		if oldTeamId != playerTeamId {
			offset := GetFileOffsetMap()[BuildPlayerStartKey(i)]
			WriteUint32AtFileOffset(inputFilename, offset+24, int(playerTeamId))

			player := saveOutput.PlayerData[i]
			countryName, _ := GetCountryInfo(player.CountryId)
			fmt.Printf("Player %d (%s): Team %d → Team %d\n", i, countryName, oldTeamId, playerTeamId)
			convertedCount++
		}
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("Team conversion completed: %d players converted to team %d\n", convertedCount, playerTeamId)
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
