package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/samuelyuan/WorldConqueror4SaveEditor/fileio"
)

func ConvertCoordinates(coordinateCode int, unitOwnerData [][]byte, gameMode int) (int, int) {
	row := int(coordinateCode) / len(unitOwnerData[0])
	if gameMode == 2 { // subtract 2 from row if conquest
		row -= 2
	}

	col := int(coordinateCode) % len(unitOwnerData[0])
	fmt.Printf("Coordinate %v (row: %v, column: %v)\n", coordinateCode, row, col)
	return row, col
}

func main() {
	inputFilenamePtr := flag.String("input", "", "input filename")
	commandPtr := flag.String("command", "", "")
	oldValuePtr := flag.String("oldvalue", "", "Old value")
	newValuePtr := flag.String("value", "", "New value")
	xPtr := flag.Int("x", -1, "x")
	yPtr := flag.Int("y", -1, "y")
	flag.Parse()

	inputFilename := *inputFilenamePtr
	command := *commandPtr

	saveOutput, err := fileio.ReadSaveFile(inputFilename)
	if err != nil {
		log.Fatal(err)
	}

	if command == "list-players" {
		countMap := make(map[byte]int)
		for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
			for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
				if saveOutput.UnitOwnerData[i][j] == 255 {
					continue
				}
				owner := saveOutput.UnitOwnerData[i][j]
				existingCount, ok := countMap[owner]
				if !ok {
					countMap[owner] = 1
				} else {
					countMap[owner] = existingCount + 1
				}
			}
		}

		for i := 0; i < len(saveOutput.PlayerData); i++ {
			player := saveOutput.PlayerData[i]
			fmt.Println(fmt.Sprintf("Player %v: CountryId %v, TeamId %v, units owned: %v", i, player.CountryId, player.TeamId, countMap[byte(i)]))
		}
	} else if command == "list-player-tiles" {
		player, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		for i := 0; i < len(saveOutput.UnitOwnerData); i++ {
			for j := 0; j < len(saveOutput.UnitOwnerData[i]); j++ {
				if saveOutput.UnitOwnerData[i][j] != 255 && saveOutput.UnitOwnerData[i][j] != 0 {
					if saveOutput.UnitOwnerData[i][j] != byte(player) {
						continue
					}

					fmt.Println(fmt.Sprintf("Player %v owns unit at tile (%v, %v)", player, i, j))
				}
			}
		}
	} else if command == "list-cities" {
		fmt.Println("Map rows:", len(saveOutput.UnitOwnerData), ", columns:", len(saveOutput.UnitOwnerData[0]))
		for i := 0; i < len(saveOutput.Cities); i++ {
			city := saveOutput.Cities[i]
			row, col := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			fmt.Printf("City %v (owner: %v): %+v\n", i, saveOutput.UnitOwnerData[row][col], city)
		}
	} else if command == "list-units" {
		for i := 0; i < len(saveOutput.Units); i++ {
			unit := saveOutput.Units[i]
			row, col := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			fmt.Printf("Unit %v (owner: %v): %+v\n", i, saveOutput.UnitOwnerData[row][col], unit)
		}
	} else if command == "list-generals" {
		for i := 0; i < len(saveOutput.Units); i++ {
			unit := saveOutput.Units[i]
			row, col := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			if unit.GeneralId > 0 {
				fmt.Printf("General (unit %v, owner: %v): %+v\n", i, saveOutput.UnitOwnerData[row][col], unit)
			}
		}
	} else if command == "max-money" {
		offset := fileio.GetFileOffsetMap()[fileio.BuildPlayerStartKey(0)]
		fileio.WriteUint16AtFileOffset(inputFilename, offset + 8, 9999)
		fileio.WriteUint16AtFileOffset(inputFilename, offset + 12, 9999)
		fileio.WriteUint16AtFileOffset(inputFilename, offset + 16, 9999)
		fmt.Println("Set max currency to 9999 for player 0")
	} else if command == "max-city-tech" {
		for i := 0; i < len(saveOutput.Cities); i++ {
			city := saveOutput.Cities[i]
			row, col := ConvertCoordinates(int(city.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			owner := saveOutput.UnitOwnerData[row][col]
			fmt.Printf("City %v (owner: %v): %+v\n", i, owner, city)
			if owner != 0 {
				continue
			}

			offset := fileio.GetFileOffsetMap()[fileio.BuildCityStartKey(i)]
			for techCount := 0; techCount < 6; techCount++ {
				fileio.WriteUint8AtFileOffset(inputFilename, offset + 24 + techCount, 4)
			}
		}

		fmt.Println("Set max city tech to level 4 for player 0")
	} else if command == "restore-allies" {
		playerTeamId := saveOutput.PlayerData[0].TeamId

		count := 0
		for i := 0; i < len(saveOutput.Units); i++ {
			unit := saveOutput.Units[i]
			row, col := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			fmt.Println("Coordinate ")
			owner := saveOutput.UnitOwnerData[row][col]
			if saveOutput.PlayerData[owner].TeamId == playerTeamId {
				fmt.Println("Restore unit", i, "health to", unit.MaxHealth)
				fileio.WriteUint16AtFileOffset(inputFilename, fileio.GetFileOffsetMap()[fileio.BuildUnitHealthKey(i)], int(unit.MaxHealth))
				count += 1
			}
		}
		fmt.Println("Restored allies. Changed", count, "units to have max health.")
	} else if command == "weaken-enemy" {
		playerTeamId := saveOutput.PlayerData[0].TeamId

		fmt.Println("Current player teamId", playerTeamId)

		count := 0
		for i := 0; i < len(saveOutput.Units); i++ {
			unit := saveOutput.Units[i]
			row, col := ConvertCoordinates(int(unit.CoordinateCode), saveOutput.UnitOwnerData, int(saveOutput.SaveHeader.GameMode))
			owner := saveOutput.UnitOwnerData[row][col]

			if int(owner) >= len(saveOutput.PlayerData) {
				fmt.Println("Invalid owner", owner, ", skip")
				continue
			}

			if saveOutput.PlayerData[owner].TeamId != playerTeamId {
				if unit.UnitType == 39 {
					fmt.Println("Reduce enemy city", i, "health to 0")
					fileio.WriteUint16AtFileOffset(inputFilename, fileio.GetFileOffsetMap()[fileio.BuildUnitHealthKey(i)], 0)
				} else {
					fmt.Println("Reduce enemy unit", i, "health to 1")
					fileio.WriteUint16AtFileOffset(inputFilename, fileio.GetFileOffsetMap()[fileio.BuildUnitHealthKey(i)], 1)
				}

				count += 1
			}
		}
		fmt.Println("Weakened enemies. Changed", count, "units to have 1 health.")
	} else if command == "convert-player" {
		oldPlayer, err := strconv.Atoi(*oldValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		newPlayer, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
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
		targetX := *xPtr
		targetY := *yPtr
		newPlayer, err := strconv.Atoi(*newValuePtr)
		if err != nil {
			log.Fatal(err)
		}
		oldPlayer := saveOutput.UnitOwnerData[targetY][targetX]
		if oldPlayer == 255 {
			log.Fatal(fmt.Sprintf("Can't convert tile at (%v, %v) without owner. Row: %v", targetY, targetX, saveOutput.UnitOwnerData[targetY]))
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
			fileio.WriteUint32AtFileOffset(inputFilename, offset + 24, int(playerTeamId))
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
		log.Fatal("Unrecognized command:", command)
	}
}
