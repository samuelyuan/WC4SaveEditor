package fileio

import (
	"fmt"
)

// UnitTileInfo holds information about a tile with a unit
type UnitTileInfo struct {
	X        int
	Y        int
	Owner    byte
	UnitType string
	UnitName string
}

// InteractiveConvertPlayer provides an interactive interface for player conversion
func InteractiveConvertPlayer(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Player Conversion ===")

	// Show available players and get selection
	oldPlayer := getPlayerSelection("old player ID to convert FROM", saveOutput)
	if oldPlayer == -1 {
		return
	}

	newPlayer := getPlayerSelection("new player ID to convert TO", saveOutput)
	if newPlayer == -1 {
		return
	}

	// Show preview and get confirmation
	if showPlayerConversionPreview(saveOutput, oldPlayer, newPlayer) {
		performPlayerConversion(inputFilename, saveOutput, oldPlayer, newPlayer)
	}
}

// getPlayerSelection shows available players and gets user selection
func getPlayerSelection(prompt string, saveOutput *WC4SaveOutput) int {
	showAvailablePlayers(saveOutput)

	fmt.Printf("\nEnter %s: ", prompt)
	var playerStr string
	fmt.Scanln(&playerStr)
	player := parseInt(playerStr, prompt)
	if player == -1 {
		return -1
	}

	// Validate player exists
	if player < 0 || player >= len(saveOutput.PlayerData) {
		fmt.Printf("Error: Player ID %d does not exist\n", player)
		return -1
	}

	return player
}

// showAvailablePlayers displays all available players
func showAvailablePlayers(saveOutput *WC4SaveOutput) {
	fmt.Println("\nAvailable Players:")
	sortedPlayerIDs := GetSortedPlayerIDs(saveOutput.PlayerData)
	for _, playerID := range sortedPlayerIDs {
		player := saveOutput.PlayerData[playerID]
		countryName, _ := GetCountryInfoFromData(player)
		fmt.Printf("  %d: %s (CountryId %d, TeamId %d)\n", playerID, countryName, player.CountryId, player.TeamId)
	}
}

// showPlayerConversionPreview shows conversion preview and returns true if user confirms
func showPlayerConversionPreview(saveOutput *WC4SaveOutput, oldPlayer, newPlayer int) bool {
	oldPlayerData := saveOutput.PlayerData[oldPlayer]
	newPlayerData := saveOutput.PlayerData[newPlayer]
	oldCountryName, _ := GetCountryInfoFromData(oldPlayerData)
	newCountryName, _ := GetCountryInfoFromData(newPlayerData)

	fmt.Printf("\n=== Conversion Preview ===\n")
	fmt.Printf("FROM: Player %d (%s - CountryId %d, TeamId %d)\n", oldPlayer, oldCountryName, oldPlayerData.CountryId, oldPlayerData.TeamId)
	fmt.Printf("TO:   Player %d (%s - CountryId %d, TeamId %d)\n", newPlayer, newCountryName, newPlayerData.CountryId, newPlayerData.TeamId)

	// Count affected tiles
	affectedCount := countPlayerTiles(saveOutput, oldPlayer)
	fmt.Printf("This will convert %d units/tiles from Player %d to Player %d\n", affectedCount, oldPlayer, newPlayer)

	// Ask for confirmation
	fmt.Print("\nProceed with conversion? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	return confirm == "y" || confirm == "Y" || confirm == "yes" || confirm == "Yes"
}

// countPlayerTiles counts how many tiles a player owns
func countPlayerTiles(saveOutput *WC4SaveOutput, playerID int) int {
	count := 0
	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			if saveOutput.UnitOwnerData[i][j] == byte(playerID) {
				count++
			}
		}
	}
	return count
}

// performPlayerConversion executes the player conversion
func performPlayerConversion(inputFilename string, saveOutput *WC4SaveOutput, oldPlayer, newPlayer int) {
	fmt.Println("\nPerforming conversion...")
	stats := ConvertPlayer(inputFilename, saveOutput, oldPlayer, newPlayer)
	stats.PrintSummary("Convert Player", saveOutput)
	fmt.Println("Conversion completed!")
}

// InteractiveConvertTile provides an interactive interface for tile conversion
func InteractiveConvertTile(inputFilename string, saveOutput *WC4SaveOutput) {
	fmt.Println("\n=== Interactive Tile Conversion ===")

	// Find and display all tiles with units
	unitTiles := findAllUnitTiles(saveOutput)
	if len(unitTiles) == 0 {
		fmt.Println("No units found on the map.")
		return
	}

	// Show tiles and get user selection
	selectedTile := selectTileFromList(unitTiles, saveOutput)
	if selectedTile == nil {
		return
	}

	// Get new player and show preview
	newPlayer := getPlayerSelection("new player ID to convert TO", saveOutput)
	if newPlayer == -1 {
		return
	}

	// Show preview and get confirmation
	if showTileConversionPreview(selectedTile, newPlayer, saveOutput) {
		performTileConversion(inputFilename, saveOutput, selectedTile, newPlayer)
	}
}

// findAllUnitTiles scans the map and finds all tiles with units
func findAllUnitTiles(saveOutput *WC4SaveOutput) []UnitTileInfo {
	var unitTiles []UnitTileInfo

	for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
		for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
			owner := saveOutput.UnitOwnerData[i][j]
			if owner != TileUnowned { // Skip unowned tiles
				unitInfo := identifyUnitAtLocation(saveOutput, i, j, owner)
				unitTiles = append(unitTiles, unitInfo)
			}
		}
	}

	return unitTiles
}

// identifyUnitAtLocation identifies what unit/city is at a specific location
func identifyUnitAtLocation(saveOutput *WC4SaveOutput, x, y int, owner byte) UnitTileInfo {
	coordinateCode := ConvertToCoordinateCode(x, y, saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))

	// Check if there's a unit at this coordinate
	for _, unit := range saveOutput.Units {
		if int(unit.CoordinateCode) == coordinateCode {
			unitTypeName := GetUnitTypeName(unit.UnitType)
			unitType := unitTypeName
			if unitType == "" {
				unitType = fmt.Sprintf("Type %d", unit.UnitType)
			}
			return UnitTileInfo{
				X:        x,
				Y:        y,
				Owner:    owner,
				UnitType: unitType,
				UnitName: unitType,
			}
		}
	}

	// Check if there's a city at this coordinate
	for _, city := range saveOutput.Cities {
		if int(city.CoordinateCode) == coordinateCode {
			cityName, exists := GetCityName(city.CityId)
			unitName := cityName
			if !exists {
				unitName = fmt.Sprintf("City (ID %d)", city.CityId)
			}
			return UnitTileInfo{
				X:        x,
				Y:        y,
				Owner:    owner,
				UnitType: "City",
				UnitName: unitName,
			}
		}
	}

	// Default to structure if nothing else found
	return UnitTileInfo{
		X:        x,
		Y:        y,
		Owner:    owner,
		UnitType: "Structure",
		UnitName: "Structure",
	}
}

// selectTileFromList displays tiles and gets user selection
func selectTileFromList(unitTiles []UnitTileInfo, saveOutput *WC4SaveOutput) *UnitTileInfo {
	// Show available tiles
	fmt.Printf("\nFound %d tiles with units:\n", len(unitTiles))

	// Create table data
	tableFormatter := NewTableFormatter()
	headers := []ColumnConfig{
		{"Index", 7, "right"},
		{"Coordinates", 13, "center"},
		{"Owner", 20, "left"},
		{"Unit Type/Name", 25, "left"},
	}

	var data [][]string
	for i, tile := range unitTiles {
		ownerName := getOwnerDisplayName(tile.Owner, saveOutput)
		displayName := getUnitDisplayName(tile)
		row := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("(%d,%d)", tile.X, tile.Y),
			ownerName,
			displayName,
		}
		data = append(data, row)
	}

	tableFormatter.PrintTable(headers, data)

	// Get tile selection
	fmt.Printf("\nSelect tile to convert (1-%d): ", len(unitTiles))
	var tileStr string
	fmt.Scanln(&tileStr)
	tileIndex := parseInt(tileStr, "tile index")
	if tileIndex == -1 {
		return nil
	}

	if tileIndex < 1 || tileIndex > len(unitTiles) {
		fmt.Printf("Error: Invalid tile index %d\n", tileIndex)
		return nil
	}

	return &unitTiles[tileIndex-1]
}

// getOwnerDisplayName returns a formatted string for the owner
func getOwnerDisplayName(owner byte, saveOutput *WC4SaveOutput) string {
	if int(owner) < len(saveOutput.PlayerData) {
		player := saveOutput.PlayerData[owner]
		countryName, _ := GetCountryInfoFromData(player)
		return fmt.Sprintf("P%d (%s)", owner, countryName)
	}
	return "Unknown"
}

// getUnitDisplayName returns the display name for a unit
func getUnitDisplayName(tile UnitTileInfo) string {
	if tile.UnitType == "City" {
		return tile.UnitName
	}
	return tile.UnitType
}

// showTileConversionPreview shows conversion preview and returns true if user confirms
func showTileConversionPreview(selectedTile *UnitTileInfo, newPlayer int, saveOutput *WC4SaveOutput) bool {
	newPlayerData := saveOutput.PlayerData[newPlayer]
	newCountryName, _ := GetCountryInfoFromData(newPlayerData)

	fmt.Printf("\n=== Conversion Preview ===\n")
	displayName := getUnitDisplayName(*selectedTile)
	fmt.Printf("Tile: (%d, %d) - %s\n", selectedTile.X, selectedTile.Y, displayName)
	fmt.Printf("Current owner: Player %d\n", selectedTile.Owner)
	fmt.Printf("New owner: Player %d (%s - CountryId %d, TeamId %d)\n", newPlayer, newCountryName, newPlayerData.CountryId, newPlayerData.TeamId)

	// Ask for confirmation
	fmt.Print("\nProceed with conversion? (y/N): ")
	var confirm string
	fmt.Scanln(&confirm)

	return confirm == "y" || confirm == "Y" || confirm == "yes" || confirm == "Yes"
}

// performTileConversion executes the tile conversion
func performTileConversion(inputFilename string, saveOutput *WC4SaveOutput, selectedTile *UnitTileInfo, newPlayer int) {
	fmt.Println("\nPerforming conversion...")
	err := ConvertTile(inputFilename, saveOutput, selectedTile.Y, selectedTile.X, newPlayer)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Conversion completed!")
	}
}

// parseInt is a helper function to parse integer input with error handling
func parseInt(input, fieldName string) int {
	var result int
	_, err := fmt.Sscanf(input, "%d", &result)
	if err != nil {
		fmt.Printf("Error: Invalid %s '%s'. Please enter a valid number.\n", fieldName, input)
		return -1
	}
	return result
}
