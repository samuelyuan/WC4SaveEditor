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
		count := 0
		for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
			for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
				if saveOutput.UnitOwnerData[i][j] != 255 && saveOutput.UnitOwnerData[i][j] != 0 {
					if saveOutput.UnitOwnerData[i][j] != byte(oldPlayer) {
						continue
					}

					fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %v", i, j, oldPlayer, newPlayer))
					saveOutput.UnitOwnerData[i][j] = byte(newPlayer)
					count += 1
				}
			}
		}
		fileio.WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
		fmt.Println("Changed", count, "tiles")
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
		oldPlayer := saveOutput.UnitOwnerData[targetY][targetX]
		if oldPlayer == 255 {
			fmt.Printf("Error: Can't convert tile at (%v, %v) - tile has no owner\n", targetY, targetX)
			return
		}
		fileio.WriteUnitOwnerToFile(inputFilename, newPlayer, targetX, targetY)
		fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to %v", targetY, targetX, oldPlayer, newPlayer))
	} else if command == "convert-all-allies" {
		playerTeamId := saveOutput.PlayerData[0].TeamId

		count := 0
		for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
			for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
				if saveOutput.UnitOwnerData[i][j] == 255 {
					continue
				}

				if saveOutput.UnitOwnerData[i][j] == 0 {
					continue
				}

				oldValue := saveOutput.UnitOwnerData[i][j]
				if saveOutput.PlayerData[oldValue].TeamId == playerTeamId {
					fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to 0", i, j, oldValue))
					saveOutput.UnitOwnerData[i][j] = 0
					count += 1
				}
			}
		}
		fileio.WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
		fmt.Println("Converted all allies. Changed", count, "allied units")
	} else if command == "convert-team" {
		playerTeamId := saveOutput.PlayerData[0].TeamId
		for i := 1; i < len(saveOutput.PlayerData); i++ {
			offset := fileio.GetFileOffsetMap()[fileio.BuildPlayerStartKey(i)]
			fileio.WriteUint32AtFileOffset(inputFilename, offset+24, int(playerTeamId))
			fmt.Println("Converting player", i, "from team", saveOutput.PlayerData[i].TeamId, "to team", playerTeamId)
		}
	} else if command == "convert-all-players" {
		count := 0
		for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
			for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
				if saveOutput.UnitOwnerData[i][j] != 255 && saveOutput.UnitOwnerData[i][j] != 0 {
					oldValue := saveOutput.UnitOwnerData[i][j]
					fmt.Println(fmt.Sprintf("Changed owner at (%v, %v) from %v to 0", i, j, oldValue))
					saveOutput.UnitOwnerData[i][j] = 0
					count += 1
				}
			}
		}
		fileio.WriteAllUnitOwnersToFile(inputFilename, saveOutput.UnitOwnerData)
		fmt.Println("Converted all players. Changed", count, "units")
	} else {
		fmt.Printf("Error: Unrecognized command '%s'\n", command)
		fmt.Println("Use -help to see available commands")
		return
	}
}
