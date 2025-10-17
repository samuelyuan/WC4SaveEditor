package main

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/samuelyuan/WC4SaveEditor/fileio"
)

func warnf(format string, a ...interface{}) {
	fmt.Printf("[WARN] "+format+"\n", a...)
}

func parseInt(s string, fieldName string) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("Error: Invalid %s '%s'. Must be a number.\n", fieldName, s)
		return -1
	}
	return val
}

func showHelp() {
	fmt.Println("WC4SaveEditor - World Conqueror 4 Save File Editor")
	fmt.Println("")
	fmt.Println("Usage: WC4SaveEditor.exe -input <savefile> -command <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  list-players, list-cities, list-units, list-units-by-map, list-tiles, list-generals")
	fmt.Println("  max-money, max-city-tech, restore-allies, weaken-enemy")
	fmt.Println("  convert-player, convert-tile, convert-all-allies, convert-team, convert-all-players")
	fmt.Println("")
	fmt.Println("Use -help to show this message")
}

func main() {
	inputFilenamePtr := flag.String("input", "", "Path to the World Conqueror 4 save file")
	commandPtr := flag.String("command", "", "Command to execute (use -help for list of commands)")
	oldValuePtr := flag.String("oldvalue", "", "Old value (used with convert-player command)")
	newValuePtr := flag.String("value", "", "New value (used with various commands)")
	xPtr := flag.Int("x", -1, "X coordinate (used with convert-tile command)")
	yPtr := flag.Int("y", -1, "Y coordinate (used with convert-tile command)")
	helpPtr := flag.Bool("help", false, "Show help information")
	flag.Parse()

	// Show help if requested or no command provided
	if *helpPtr || *commandPtr == "" {
		showHelp()
		return
	}

	inputFilename := *inputFilenamePtr
	command := *commandPtr

	// Validate required parameters
	if inputFilename == "" {
		fmt.Println("Error: -input flag is required")
		fmt.Println("Use -help for usage information")
		return
	}

	saveOutput, err := fileio.ReadSaveFile(inputFilename)
	if err != nil {
		fmt.Printf("Error reading save file '%s': %v\n", inputFilename, err)
		fmt.Println("Make sure the file exists and is a valid World Conqueror 4 save file")
		return
	}

	if command == "list-players" {
		fileio.ListPlayers(saveOutput)
	} else if command == "list-cities" {
		fileio.ListCities(saveOutput)
	} else if command == "list-units" {
		fileio.ListUnits(saveOutput)
	} else if command == "list-units-by-map" {
		fileio.ListUnitsByMap(saveOutput)
	} else if command == "list-tiles" {
		fileio.ListTilesByOwner(saveOutput)
	} else if command == "list-generals" {
		fileio.ListGenerals(saveOutput)
	} else if command == "max-money" {
		fileio.SetPlayerMaxCurrency(inputFilename, 0, 9999)
	} else if command == "max-city-tech" {
		fileio.SetPlayerMaxCityTech(inputFilename, saveOutput, 0, 4)
	} else if command == "restore-allies" {
		fileio.RestoreAlliesHealth(inputFilename, saveOutput, 0)
	} else if command == "weaken-enemy" {
		fileio.WeakenEnemies(inputFilename, saveOutput, 0)
	} else if command == "convert-player" {
		if *oldValuePtr == "" || *newValuePtr == "" {
			fmt.Println("Error: Both -oldvalue and -value flags are required for convert-player command")
			fmt.Println("Example: -command convert-player -oldvalue 2 -value 0")
			return
		}
		oldPlayer := parseInt(*oldValuePtr, "old player ID")
		if oldPlayer == -1 {
			return
		}
		newPlayer := parseInt(*newValuePtr, "new player ID")
		if newPlayer == -1 {
			return
		}
		stats := fileio.ConvertPlayer(inputFilename, saveOutput, oldPlayer, newPlayer)
		stats.PrintSummary("Convert Player", saveOutput)
	} else if command == "convert-tile" {
		if *xPtr == -1 || *yPtr == -1 || *newValuePtr == "" {
			fmt.Println("Error: -x, -y, and -value flags are required for convert-tile command")
			fmt.Println("Example: -command convert-tile -x 10 -y 5 -value 0")
			return
		}
		targetX := *xPtr
		targetY := *yPtr
		newPlayer := parseInt(*newValuePtr, "player ID")
		if newPlayer == -1 {
			return
		}
		err := fileio.ConvertTile(inputFilename, saveOutput, targetX, targetY, newPlayer)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
	} else if command == "convert-all-allies" {
		stats := fileio.ConvertAllAllies(inputFilename, saveOutput)
		stats.PrintSummary("Convert All Allies", saveOutput)
	} else if command == "convert-team" {
		fileio.ConvertTeam(inputFilename, saveOutput)
	} else if command == "convert-all-players" {
		stats := fileio.ConvertAllPlayers(inputFilename, saveOutput)
		stats.PrintSummary("Convert All Players", saveOutput)
	} else {
		fmt.Printf("Error: Unrecognized command '%s'\n", command)
		fmt.Println("Use -help to see available commands")
		return
	}
}
