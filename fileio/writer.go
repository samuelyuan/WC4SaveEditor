package fileio

import (
	"fmt"
	"log"
	"os"
)


// Unit type constants (used in core writing operations)
const (
	UnitTypeCity = 39 // City unit type
)

// Tile ownership constants
const (
	TileUnowned = 255 // Value indicating unowned/neutral tile
)

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
			unitTypeName := "unit"
			cityInfo := ""
			if unitInfo.Unit.UnitType == UnitTypeCity {
				unitTypeName = "city"
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
			unitTypeName := "unit"
			cityInfo := ""
			if unitInfo.Unit.UnitType == UnitTypeCity {
				unitTypeName = "city"
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
	for unitType, count := range unitsByType {
		unitTypeName := "unit"
		if unitType == UnitTypeCity {
			unitTypeName = "city"
		}
		fmt.Printf("Type %d (%s): %d units\n", unitType, unitTypeName, count)
	}

	fmt.Printf("\nWeakened enemies. Changed %d units to have low health.\n", unitsWeakened)
}

func ConvertCoordinates(coordinateCode int, unitOwnerData [][]byte, gameMode int) (int, int, bool) {
	row := int(coordinateCode) / len(unitOwnerData[0])
	if gameMode == 2 { // subtract 2 from row if conquest
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
	if gameMode == 2 { // conquest mode
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

func GetCountryInfo(countryID uint32) (string, string) {
	countryName, hasName := GetCountryName(uint8(countryID))
	if !hasName {
		countryName = "Unknown"
	}
	return countryName, fmt.Sprintf("CountryId %d", countryID)
}

func ConvertCoordinatesWithDebug(coordinateCode int, unitOwnerData [][]byte, gameMode int) (int, int, bool) {
	row, col, valid := ConvertCoordinates(coordinateCode, unitOwnerData, gameMode)
	if valid {
		fmt.Printf("Coordinate %v (row: %v, column: %v)\n", coordinateCode, row, col)
	}
	return row, col, valid
}
