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
	fmt.Println("  list-players, list-cities, list-units, list-units-by-map, list-tiles, list-generals, list-landmines, list-teams")
	fmt.Println("  max-money, max-city-tech, heal-allies, weaken-enemies")
	fmt.Println("  conquer-player (interactive), transfer-tile (interactive)")
	fmt.Println("  conquer-allies, unite-team, conquer-all")
	fmt.Println("  visualize-map")
	fmt.Println("")
	fmt.Println("Use -help to show this message")
}

func main() {
	inputFilenamePtr := flag.String("input", "", "Path to the World Conqueror 4 save file")
	commandPtr := flag.String("command", "", "Command to execute (use -help for list of commands)")
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
	} else if command == "list-landmines" {
		fileio.ListLandmines(saveOutput)
	} else if command == "list-teams" {
		fileio.ListTeams(saveOutput)
	} else if command == "max-money" {
		fileio.SetPlayerMaxCurrency(inputFilename, 0, 9999)
	} else if command == "max-city-tech" {
		fileio.SetPlayerMaxCityTech(inputFilename, saveOutput, 0, 4)
	} else if command == "heal-allies" {
		fileio.RestoreAlliesHealth(inputFilename, saveOutput, 0)
	} else if command == "weaken-enemies" {
		fileio.WeakenEnemies(inputFilename, saveOutput, 0)
	} else if command == "conquer-player" {
		fileio.InteractiveConvertPlayer(inputFilename, saveOutput)
	} else if command == "transfer-tile" {
		fileio.InteractiveConvertTile(inputFilename, saveOutput)
	} else if command == "conquer-allies" {
		stats := fileio.ConvertAllAllies(inputFilename, saveOutput)
		stats.PrintSummary("Conquer Allies", saveOutput)
	} else if command == "unite-team" {
		fileio.ConvertTeam(inputFilename, saveOutput)
	} else if command == "conquer-all" {
		stats := fileio.ConvertAllPlayers(inputFilename, saveOutput)
		stats.PrintSummary("Conquer All Players", saveOutput)
	} else if command == "visualize-map" {
		outputPath := "map_visualization.svg"
		err := fileio.VisualizeMap(saveOutput, outputPath)
		if err != nil {
			fmt.Printf("Error creating map visualization: %v\n", err)
			return
		}
		fmt.Printf("Map visualization saved to: %s\n", outputPath)
	} else {
		fmt.Printf("Error: Unrecognized command '%s'\n", command)
		fmt.Println("Use -help to see available commands")
		return
	}
}
