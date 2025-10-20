package fileio

import (
	"fmt"
	"strings"
)

// Data structures for table data

type PlayerTableData struct {
	PlayerID    int
	CountryName string
	CountryID   int
	TeamID      int
	UnitsOwned  int
}

type UnitTableData struct {
	UnitID   int
	UnitType string
	Level    int
	Health   string
	Position string
}

type GeneralTableData struct {
	UnitID      int
	UnitType    string
	GeneralName string
	Level       int
	Health      string
	Position    string
}

type CityTableData struct {
	CityID   int
	CityName string
	Position string
}

type PlayerSummaryData struct {
	PlayerID     int
	CountryName  string
	CountryID    int
	TeamID       int
	CityCount    int
	TotalTiles   int
	LandPercent  float64
	WorldPercent float64
}

type CityTileData struct {
	CityIndex      int
	CityName       string
	CityID         int
	Position       string
	CoordinateCode int
	TileCount      int
}

type TerritoryControlData struct {
	PlayerID     int
	CountryName  string
	CountryID    int
	TileCount    int
	LandPercent  float64
	WorldPercent float64
}

type TeamTerritoryData struct {
	TeamID       int
	PlayerCount  int
	TileCount    int
	LandPercent  float64
	WorldPercent float64
	Players      []string
}

type LandmineTableData struct {
	LandmineID int
	Position   string
	Owner      string
	Health     string
}

type TeamAnalysisData struct {
	TeamID       int
	PlayerCount  int
	CityCount    int
	UnitCount    int
	LandPercent  float64
	WorldPercent float64
	Players      []string
}

// ColumnConfig defines the configuration for a table column
type ColumnConfig struct {
	Header    string
	Width     int
	Alignment string // "left", "right", "center"
}

// TableFormatter provides utilities for creating formatted tables
type TableFormatter struct{}

// NewTableFormatter creates a new TableFormatter instance
func NewTableFormatter() *TableFormatter {
	return &TableFormatter{}
}

// PrintTable prints a formatted table with the given headers and data
func (tf *TableFormatter) PrintTable(headers []ColumnConfig, data [][]string) {
	if len(headers) == 0 || len(data) == 0 {
		return
	}

	// Calculate total width
	totalWidth := len(headers) - 1 // for separators
	for _, header := range headers {
		totalWidth += header.Width
	}

	// Print top border
	tf.printTopBorder(headers)

	// Print header row
	tf.printHeaderRow(headers)

	// Print separator row
	tf.printSeparatorRow(headers)

	// Print data rows
	for _, row := range data {
		tf.printDataRow(headers, row)
	}

	// Print bottom border
	tf.printBottomBorder(headers)
}

// printTopBorder prints the top border of the table
func (tf *TableFormatter) printTopBorder(headers []ColumnConfig) {
	fmt.Print("┌")
	for i, header := range headers {
		fmt.Print(strings.Repeat("─", header.Width))
		if i < len(headers)-1 {
			fmt.Print("┬")
		}
	}
	fmt.Println("┐")
}

// printHeaderRow prints the header row of the table
func (tf *TableFormatter) printHeaderRow(headers []ColumnConfig) {
	fmt.Print("│")
	for i, header := range headers {
		fmt.Printf(" %-*s ", header.Width-2, header.Header)
		if i < len(headers)-1 {
			fmt.Print("│")
		}
	}
	fmt.Println("│")
}

// printSeparatorRow prints the separator row between header and data
func (tf *TableFormatter) printSeparatorRow(headers []ColumnConfig) {
	fmt.Print("├")
	for i, header := range headers {
		fmt.Print(strings.Repeat("─", header.Width))
		if i < len(headers)-1 {
			fmt.Print("┼")
		}
	}
	fmt.Println("┤")
}

// printDataRow prints a single data row
func (tf *TableFormatter) printDataRow(headers []ColumnConfig, row []string) {
	fmt.Print("│")
	for i, header := range headers {
		value := ""
		if i < len(row) {
			value = row[i]
		}

		switch header.Alignment {
		case "right":
			fmt.Printf(" %*s ", header.Width-2, value)
		case "center":
			padding := (header.Width - 2 - len(value)) / 2
			fmt.Printf(" %*s%s%*s ", padding, "", value, header.Width-2-len(value)-padding, "")
		default: // left
			fmt.Printf(" %-*s ", header.Width-2, value)
		}

		if i < len(headers)-1 {
			fmt.Print("│")
		}
	}
	fmt.Println("│")
}

// printBottomBorder prints the bottom border of the table
func (tf *TableFormatter) printBottomBorder(headers []ColumnConfig) {
	fmt.Print("└")
	for i, header := range headers {
		fmt.Print(strings.Repeat("─", header.Width))
		if i < len(headers)-1 {
			fmt.Print("┴")
		}
	}
	fmt.Println("┘")
}

// Convenience methods for common table types

// PrintPlayerTable prints a table for player data
func (tf *TableFormatter) PrintPlayerTable(players []PlayerTableData) {
	headers := []ColumnConfig{
		{"Player", 8, "right"},
		{"Country", 20, "left"},
		{"Country ID", 13, "right"},
		{"Team", 8, "right"},
		{"Units Owned", 13, "right"},
	}

	var data [][]string
	for _, player := range players {
		row := []string{
			fmt.Sprintf("%d", player.PlayerID),
			player.CountryName,
			fmt.Sprintf("%d", player.CountryID),
			fmt.Sprintf("%d", player.TeamID),
			fmt.Sprintf("%d", player.UnitsOwned),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintUnitTable prints a table for unit data
func (tf *TableFormatter) PrintUnitTable(units []UnitTableData) {
	headers := []ColumnConfig{
		{"Unit", 7, "right"},
		{"Type", 21, "left"},
		{"Level", 7, "right"},
		{"Health", 11, "right"},
		{"Position", 13, "right"},
	}

	var data [][]string
	for _, unit := range units {
		row := []string{
			fmt.Sprintf("%d", unit.UnitID),
			unit.UnitType,
			fmt.Sprintf("%d", unit.Level),
			unit.Health,
			unit.Position,
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintGeneralTable prints a table for general data
func (tf *TableFormatter) PrintGeneralTable(generals []GeneralTableData) {
	headers := []ColumnConfig{
		{"Unit", 7, "right"},
		{"Type", 21, "left"},
		{"General", 13, "left"},
		{"Level", 7, "right"},
		{"Health", 11, "right"},
		{"Position", 13, "right"},
	}

	var data [][]string
	for _, general := range generals {
		row := []string{
			fmt.Sprintf("%d", general.UnitID),
			general.UnitType,
			general.GeneralName,
			fmt.Sprintf("%d", general.Level),
			general.Health,
			general.Position,
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintCityTable prints a table for city data
func (tf *TableFormatter) PrintCityTable(cities []CityTableData) {
	headers := []ColumnConfig{
		{"City", 8, "right"},
		{"Name", 25, "left"},
		{"Position", 12, "right"},
	}

	var data [][]string
	for _, city := range cities {
		row := []string{
			fmt.Sprintf("%d", city.CityID),
			city.CityName,
			city.Position,
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintPlayerSummaryTable prints a table for player summary data
func (tf *TableFormatter) PrintPlayerSummaryTable(players []PlayerSummaryData) {
	headers := []ColumnConfig{
		{"Player", 7, "right"},
		{"Country", 15, "left"},
		{"Cities", 7, "right"},
		{"Total Tiles", 12, "right"},
		{"Land %", 8, "right"},
		{"World %", 9, "right"},
	}

	var data [][]string
	for _, player := range players {
		row := []string{
			fmt.Sprintf("%d", player.PlayerID),
			player.CountryName,
			fmt.Sprintf("%d", player.CityCount),
			fmt.Sprintf("%d", player.TotalTiles),
			fmt.Sprintf("%.1f%%", player.LandPercent),
			fmt.Sprintf("%.1f%%", player.WorldPercent),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintCityTileTable prints a table for city tile data
func (tf *TableFormatter) PrintCityTileTable(cities []CityTileData) {
	headers := []ColumnConfig{
		{"City", 7, "right"},
		{"Name", 20, "left"},
		{"Position", 12, "center"},
		{"Coord Code", 13, "right"},
		{"Tile Count", 12, "right"},
	}

	var data [][]string
	for _, city := range cities {
		row := []string{
			fmt.Sprintf("%d", city.CityIndex),
			city.CityName,
			city.Position,
			fmt.Sprintf("%d", city.CoordinateCode),
			fmt.Sprintf("%d", city.TileCount),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintTerritoryControlTable prints a table for territory control data
func (tf *TableFormatter) PrintTerritoryControlTable(territories []TerritoryControlData) {
	headers := []ColumnConfig{
		{"Player", 9, "right"},
		{"Country", 15, "left"},
		{"Country ID", 13, "right"},
		{"Tiles", 8, "right"},
		{"Land %", 8, "right"},
		{"World %", 9, "right"},
	}

	var data [][]string
	for _, territory := range territories {
		row := []string{
			fmt.Sprintf("%d", territory.PlayerID),
			territory.CountryName,
			fmt.Sprintf("%d", territory.CountryID),
			fmt.Sprintf("%d", territory.TileCount),
			fmt.Sprintf("%.1f%%", territory.LandPercent),
			fmt.Sprintf("%.1f%%", territory.WorldPercent),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// HealAlliesUnitData represents a unit in the heal-allies breakdown
type HealAlliesUnitData struct {
	UnitID    int
	UnitType  string
	Level     int
	OldHealth int
	NewHealth int
	Position  string
}

// PrintHealAlliesTable prints a table for heal-allies unit breakdown
func (tf *TableFormatter) PrintHealAlliesTable(units []HealAlliesUnitData) {
	headers := []ColumnConfig{
		{"Unit", 6, "right"},
		{"Type", 30, "left"},
		{"Level", 8, "right"},
		{"Health", 12, "right"},
		{"Position", 12, "center"},
	}

	var data [][]string
	for _, unit := range units {
		healthStr := fmt.Sprintf("%d->%d", unit.OldHealth, unit.NewHealth)
		row := []string{
			fmt.Sprintf("%d", unit.UnitID),
			unit.UnitType,
			fmt.Sprintf("%d", unit.Level),
			healthStr,
			unit.Position,
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintTeamTerritoryControlTable prints a table for team territory control data
func (tf *TableFormatter) PrintTeamTerritoryControlTable(teams []TeamTerritoryData) {
	// First, print a compact summary table
	headers := []ColumnConfig{
		{"Team", 8, "right"},
		{"Players", 10, "right"},
		{"Tiles", 10, "right"},
		{"Land %", 10, "right"},
		{"World %", 11, "right"},
	}

	var data [][]string
	for _, team := range teams {
		row := []string{
			fmt.Sprintf("%d", team.TeamID),
			fmt.Sprintf("%d", team.PlayerCount),
			fmt.Sprintf("%d", team.TileCount),
			fmt.Sprintf("%.1f%%", team.LandPercent),
			fmt.Sprintf("%.1f%%", team.WorldPercent),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)

	// Then print team member details separately
	fmt.Println("\nTeam Member Details:")
	fmt.Println("-------------------")
	for _, team := range teams {
		fmt.Printf("Team %d (%d players, %d tiles, %.1f%% land, %.1f%% world):\n",
			team.TeamID, team.PlayerCount, team.TileCount, team.LandPercent, team.WorldPercent)

		// Print players in groups of 5 for better readability
		for i := 0; i < len(team.Players); i += 5 {
			end := i + 5
			if end > len(team.Players) {
				end = len(team.Players)
			}
			playerGroup := team.Players[i:end]
			fmt.Printf("  %s\n", strings.Join(playerGroup, ", "))
		}
		fmt.Println()
	}
}

// PrintLandmineTable prints a table for landmine data
func (tf *TableFormatter) PrintLandmineTable(landmines []LandmineTableData) {
	headers := []ColumnConfig{
		{"ID", 6, "right"},
		{"Position", 14, "center"},
		{"Health", 10, "right"},
	}

	var data [][]string
	for _, landmine := range landmines {
		row := []string{
			fmt.Sprintf("%d", landmine.LandmineID),
			landmine.Position,
			landmine.Health,
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}

// PrintTeamAnalysisTable prints a table for team analysis data
func (tf *TableFormatter) PrintTeamAnalysisTable(teams []TeamAnalysisData) {
	headers := []ColumnConfig{
		{"Team", 8, "right"},
		{"Players", 10, "right"},
		{"Cities", 8, "right"},
		{"Units", 8, "right"},
		{"Land %", 10, "right"},
		{"World %", 11, "right"},
	}

	var data [][]string
	for _, team := range teams {
		row := []string{
			fmt.Sprintf("%d", team.TeamID),
			fmt.Sprintf("%d", team.PlayerCount),
			fmt.Sprintf("%d", team.CityCount),
			fmt.Sprintf("%d", team.UnitCount),
			fmt.Sprintf("%.1f%%", team.LandPercent),
			fmt.Sprintf("%.1f%%", team.WorldPercent),
		}
		data = append(data, row)
	}

	tf.PrintTable(headers, data)
}
