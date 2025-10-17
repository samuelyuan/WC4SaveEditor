package fileio

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	fileOffsetMap = make(map[string]int)
)

type SaveHeader struct {
	Magic              [4]byte
	UnknownInt1        uint32
	MapId              uint32
	GameMode           uint32 // 1 - campaign, 2 - conquest, 6 - frontier
	UnknownInt2        uint32
	UnknownInt3        uint32
	Camera             [3]float32
	UnknownInt4        uint32
	TurnNumber         uint32
	UnknownArr2        [12]byte
	SaveTimestamp      [5]uint32
	UnknownArr3        [16]byte
	UnknownInt7        uint32 // only seems to be set to non-zero value in frontier mode last mission
	UnknownInt8        uint32 // only seems to be set to non-zero value in frontier mode last mission
	MapWidth           uint32
	MapHeight          uint32
	CountryCount       uint32
	CityCount          uint32
	UnitCount          uint32
	UnknownCount1      uint32
	UnknownCount2      uint32
	UnknownArr4        [8]byte
	TurnCount1         uint32
	TurnCount2         uint32
	UnknownCount3      uint32
	UnknownCount4      uint32
	UnknownCount5      uint32
	UnknownCount6      uint32
	ImportantCityCount uint32
	UnknownArr5        [4]byte
	UnknownInt9        uint32
	UnknownInt10       uint32
	UnknownArr6        [12]byte
	LandmineCount      uint32
	UnknownArr7        [16]byte
	UnknownCount9      uint32
}

type CountryData struct {
	TurnOrder    uint32
	CountryId    uint32
	Currency     [3]uint32
	BotFlag      uint32
	TeamId       uint32
	UnknownArr2  [4]byte
	UnknownColor [2][4]byte
	PrimaryColor [4]byte
	UnknownArr4  [16]byte
	UnknownArr5  [460]byte
}

type CityData struct {
	CoordinateCode    uint16
	CityId            uint16
	BuildingType      uint8
	Apperance         uint8
	UnknownByte1      uint8
	Wonders           uint8
	UnknownArr2       [6]byte
	UnknownArr3       [8]byte
	AntiAirWeaponType uint8
	AntiAirRange      uint8
	TechLevels        [6]byte
	UnknownArr4       [2]byte
}

type UnitData struct {
	CoordinateCode      uint16
	UnitType            uint8
	Level               uint8
	Personnel           uint8
	Direction           uint8 // 0 - left, 1 - right
	Movement            uint16
	Experience          uint16
	UnknownHealth       uint16
	CurrentHealth       uint16
	MaxHealth           uint16
	GeneralId           uint16
	GeneralMilitaryRank uint8
	GeneralTitle        uint8
	GeneralBadges       [3]byte
	GeneralSkillLevels  [5]byte
	UnknownArr5         [12]byte
	MoraleValue         int8
	MoraleTurnsLeft     uint16
	UnknownArr6         [21]byte
}

type LandmineData struct {
	CoordinateCode uint16
	Owner          uint16
	UnknownArr1    [2]byte
	Health         uint16
	UnknownArr2    [4]byte
}

type WC4SaveOutput struct {
	SaveHeader    SaveHeader
	PlayerData    []CountryData
	CityTiles     [][]uint16
	UnitOwnerData [][]byte
	Cities        []CityData
	Units         []UnitData
}

func DeserializeMapHeaderFromBytes(streamReader *io.SectionReader) SaveHeader {
	mapHeaderInput := SaveHeader{}
	if err := binary.Read(streamReader, binary.LittleEndian, &mapHeaderInput); err != nil {
		log.Fatal("Failed to load MapHeaderInput: ", err)
	}
	fmt.Printf("Map Header Input: %+v\n", mapHeaderInput)
	return mapHeaderInput
}

func DeserializeCountryDataFromBytes(streamReader *io.SectionReader, count int) []CountryData {
	allPlayerData := make([]CountryData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMap(fileOffsetMap, streamReader, BuildPlayerStartKey(i))

		countryData := CountryData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &countryData); err != nil {
			log.Fatal("Failed to load country data: ", err)
		}
		allPlayerData[i] = countryData
		fmt.Printf("Player %v data: %+v\n", i, countryData)
	}
	return allPlayerData
}

func DeserializeCityTileOwnershipFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) [][]uint16 {
	allCityTiles := make([][]uint16, 0)
	for i := 0; i < mapHeight; i++ {
		cityRow := make([]uint16, 0)
		for j := 0; j < mapWidth; j++ {
			cityCoordinates := uint16(0)
			if err := binary.Read(streamReader, binary.LittleEndian, &cityCoordinates); err != nil {
				log.Fatal("Failed to load city tile ownership: ", err)
			}
			cityRow = append(cityRow, cityCoordinates)
		}
		allCityTiles = append(allCityTiles, cityRow)
		fmt.Println("City Tile Owner Row", i, ":", cityRow)
	}
	return allCityTiles
}

func DeserializeUnknownCampaignBlockFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) {
	for i := 0; i < mapHeight; i++ {
		for j := 0; j < mapWidth; j++ {
			unknownBlock := make([]byte, 16)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load city tile ownership: ", err)
			}
			fmt.Println("Unknown block:", unknownBlock)
		}
	}
}

func DeserializeUnitOwnerDataFromBytes(streamReader *io.SectionReader, mapWidth int, mapHeight int) [][]byte {
	unitOwnerData := make([][]byte, 0)

	for i := 0; i < mapHeight; i++ {
		unitOwnerRow := make([]byte, mapWidth*1)

		for j := 0; j < mapWidth; j++ {
			updateFileOffsetMap(fileOffsetMap, streamReader, buildSingleUnitOwnerStartKey(j, i))

			if err := binary.Read(streamReader, binary.LittleEndian, &unitOwnerRow[j]); err != nil {
				log.Fatal("Failed to load unit owner: ", err)
			}
		}

		unitOwnerData = append(unitOwnerData, unitOwnerRow)
		fmt.Println("Unit Owner Row", i, ":", unitOwnerRow)
	}

	return unitOwnerData
}

func DeserializeCityDataFromBytes(streamReader *io.SectionReader, count int) []CityData {
	allCities := make([]CityData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMap(fileOffsetMap, streamReader, BuildCityStartKey(i))

		cityData := CityData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &cityData); err != nil {
			log.Fatal("Failed to load cityData: ", err)
		}

		allCities[i] = cityData
		fmt.Printf("City data: %+v\n", cityData)

		if i > 0 && cityData.CoordinateCode == 0 {
			log.Fatal("Invalid city data")
		}
	}
	return allCities
}

func DeserializeUnitDataFromBytes(streamReader *io.SectionReader, count int) []UnitData {
	allUnits := make([]UnitData, count)
	for i := 0; i < count; i++ {
		updateFileOffsetMapForField(fileOffsetMap, streamReader, BuildUnitHealthKey(i), 12)

		unitData := UnitData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &unitData); err != nil {
			log.Fatal("Failed to load unit data: ", err)
		}
		allUnits[i] = unitData
		fmt.Printf("Unit data: %+v\n", unitData)
	}
	return allUnits
}

func DeserializeLandmineDataFromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		landmineData := LandmineData{}
		if err := binary.Read(streamReader, binary.LittleEndian, &landmineData); err != nil {
			log.Fatal("Failed to load landmine data: ", err)
		}
		fmt.Printf("Landmine: %+v\n", landmineData)
	}
}

func DeserializeUnknownData2FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 16)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Unknown block 2:", unknownBlock)
	}
}

func DeserializeUnknownData3FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 44)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Unknown block 3:", unknownBlock)
	}
}

func DeserializeUnknownData4FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 80)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Unknown block 4:", unknownBlock)
	}
}

func DeserializeUnknownData5FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 8)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Unknown block 5:", unknownBlock)
	}
}

func DeserializeImportantCityDataFromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 4)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Important city:", unknownBlock)
	}
}

func DeserializeUnknownData7FromBytes(streamReader *io.SectionReader, count int) {
	for i := 0; i < count; i++ {
		unknownBlock := make([]byte, 16)
		if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
			log.Fatal("Failed to load unknownBlock: ", err)
		}
		fmt.Println("Unknown block 7:", unknownBlock)
	}
}

func ReadSaveFile(inputFilename string) (*WC4SaveOutput, error) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
		return nil, err
	}
	fi, err := inputFile.Stat()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fileLength := fi.Size()
	streamReader := io.NewSectionReader(inputFile, int64(0), fileLength)

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("READING SAVE FILE SECTIONS")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n[1] SAVE HEADER")
	fmt.Println(strings.Repeat("-", 20))
	saveHeader := DeserializeMapHeaderFromBytes(streamReader)

	fmt.Println("\n[2] PLAYER DATA")
	fmt.Println(strings.Repeat("-", 20))
	allPlayerData := DeserializeCountryDataFromBytes(streamReader, int(saveHeader.CountryCount))

	isConquest := (int(saveHeader.GameMode) == 2)
	if !isConquest {
		if saveHeader.UnknownInt7 == 0 {
			fmt.Println("\n[3] CAMPAIGN BLOCK")
			fmt.Println(strings.Repeat("-", 20))
			DeserializeUnknownCampaignBlockFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))
		}
	}

	fmt.Println("\n[4] CITY TILES")
	fmt.Println(strings.Repeat("-", 20))
	allCityTiles := DeserializeCityTileOwnershipFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))

	if !isConquest {
		// required for some maps because the data is shifted
		if int(saveHeader.MapWidth)*int(saveHeader.MapHeight) != int(saveHeader.UnknownInt10) {
			fmt.Println("\n[5] UNKNOWN BLOCK 1")
			fmt.Println(strings.Repeat("-", 20))
			unknownBlock := make([]byte, 8)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load unknownBlock: ", err)
			}
		}
	}

	fmt.Println("\n[6] UNIT OWNER DATA")
	fmt.Println(strings.Repeat("-", 20))
	updateFileOffsetMap(fileOffsetMap, streamReader, buildUnitOwnerStartKey())
	unitOwnerData := DeserializeUnitOwnerDataFromBytes(streamReader, int(saveHeader.MapWidth), int(saveHeader.MapHeight))
	updateFileOffsetMap(fileOffsetMap, streamReader, buildUnitOwnerEndKey())

	if !isConquest {
		if int(saveHeader.MapWidth)*int(saveHeader.MapHeight) != int(saveHeader.UnknownInt10) {
			fmt.Println("\n[7] UNKNOWN BLOCK 2")
			fmt.Println(strings.Repeat("-", 20))
			unknownBlock := make([]byte, 4)
			if err := binary.Read(streamReader, binary.LittleEndian, &unknownBlock); err != nil {
				log.Fatal("Failed to load unknownBlock: ", err)
			}
		}
	}

	fmt.Println("\n[8] CITIES")
	fmt.Println(strings.Repeat("-", 20))
	allCities := DeserializeCityDataFromBytes(streamReader, int(saveHeader.CityCount))

	fmt.Println("\n[9] UNITS")
	fmt.Println(strings.Repeat("-", 20))
	allUnits := DeserializeUnitDataFromBytes(streamReader, int(saveHeader.UnitCount))

	fmt.Println("\n[10] LANDMINES")
	fmt.Println(strings.Repeat("-", 20))
	DeserializeLandmineDataFromBytes(streamReader, int(saveHeader.LandmineCount))

	fmt.Println("\n[11] UNKNOWN DATA SECTIONS")
	fmt.Println(strings.Repeat("-", 20))
	DeserializeUnknownData2FromBytes(streamReader, int(saveHeader.UnknownCount1))
	DeserializeUnknownData3FromBytes(streamReader, int(saveHeader.UnknownCount2))
	DeserializeUnknownData4FromBytes(streamReader, int(saveHeader.UnknownCount3))
	DeserializeUnknownData5FromBytes(streamReader, int(saveHeader.UnknownCount5))
	DeserializeUnknownData5FromBytes(streamReader, int(saveHeader.UnknownCount6))
	DeserializeImportantCityDataFromBytes(streamReader, int(saveHeader.ImportantCityCount))
	DeserializeUnknownData7FromBytes(streamReader, int(saveHeader.UnknownCount9))

	saveOutput := &WC4SaveOutput{
		SaveHeader:    saveHeader,
		PlayerData:    allPlayerData,
		CityTiles:     allCityTiles,
		UnitOwnerData: unitOwnerData,
		Cities:        allCities,
		Units:         allUnits,
	}

	// Print file summary with byte ranges
	printFileSummary(saveHeader, fileLength)

	return saveOutput, nil
}

func printFileSummary(header SaveHeader, fileLength int64) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("FILE SUMMARY")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Total file size: %d bytes (%.2f KB)\n", fileLength, float64(fileLength)/1024)
	fmt.Printf("Game mode: %d (%s)\n", header.GameMode, getGameModeName(header.GameMode))
	fmt.Printf("Map dimensions: %dx%d\n", header.MapWidth, header.MapHeight)
	fmt.Printf("Turn number: %d\n", header.TurnNumber)
	fmt.Println()

	fmt.Println("SECTION BREAKDOWN:")
	fmt.Println("-----------------")

	offset := int64(0)

	// Save Header
	headerSize := int64(208) // Size of SaveHeader struct
	fmt.Printf("Save Header:        bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+headerSize-1, headerSize)
	offset += headerSize

	// Player Data
	playerDataSize := int64(header.CountryCount * 520) // 520 bytes per country
	fmt.Printf("Player Data:        bytes 0x%04X-0x%04X (0x%04X bytes, %d countries)\n", offset, offset+playerDataSize-1, playerDataSize, header.CountryCount)
	offset += playerDataSize

	// Campaign block (if not conquest)
	if header.GameMode != 2 && header.UnknownInt7 == 0 {
		campaignBlockSize := int64(header.MapWidth * header.MapHeight * 16)
		fmt.Printf("Campaign Block:     bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+campaignBlockSize-1, campaignBlockSize)
		offset += campaignBlockSize
	}

	// City Tiles
	cityTilesSize := int64(header.MapWidth * header.MapHeight * 2) // 2 bytes per tile
	fmt.Printf("City Tiles:         bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+cityTilesSize-1, cityTilesSize)
	offset += cityTilesSize

	// Unknown block (if needed)
	if header.GameMode != 2 && int64(header.MapWidth*header.MapHeight) != int64(header.UnknownInt10) {
		unknownBlockSize := int64(8)
		fmt.Printf("Unknown Block:      bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+unknownBlockSize-1, unknownBlockSize)
		offset += unknownBlockSize
	}

	// Unit Owner Data
	unitOwnerSize := int64(header.MapWidth * header.MapHeight * 1) // 1 byte per tile
	fmt.Printf("Unit Owner Data:    bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+unitOwnerSize-1, unitOwnerSize)
	offset += unitOwnerSize

	// Another unknown block (if needed)
	if header.GameMode != 2 && int64(header.MapWidth*header.MapHeight) != int64(header.UnknownInt10) {
		unknownBlockSize := int64(4)
		fmt.Printf("Unknown Block 2:    bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, offset+unknownBlockSize-1, unknownBlockSize)
		offset += unknownBlockSize
	}

	// Cities
	citiesSize := int64(header.CityCount * 32) // 32 bytes per city
	fmt.Printf("Cities:             bytes 0x%04X-0x%04X (0x%04X bytes, %d cities)\n", offset, offset+citiesSize-1, citiesSize, header.CityCount)
	offset += citiesSize

	// Units
	unitsSize := int64(header.UnitCount * 48) // 48 bytes per unit
	fmt.Printf("Units:              bytes 0x%04X-0x%04X (0x%04X bytes, %d units)\n", offset, offset+unitsSize-1, unitsSize, header.UnitCount)
	offset += unitsSize

	// Landmines
	landminesSize := int64(header.LandmineCount * 12) // 12 bytes per landmine
	fmt.Printf("Landmines:          bytes 0x%04X-0x%04X (0x%04X bytes, %d landmines)\n", offset, offset+landminesSize-1, landminesSize, header.LandmineCount)
	offset += landminesSize

	// Unknown data sections
	unknownData2Size := int64(header.UnknownCount1 * 16)
	fmt.Printf("Unknown Data 2:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData2Size-1, unknownData2Size, header.UnknownCount1)
	offset += unknownData2Size

	unknownData3Size := int64(header.UnknownCount2 * 44)
	fmt.Printf("Unknown Data 3:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData3Size-1, unknownData3Size, header.UnknownCount2)
	offset += unknownData3Size

	unknownData4Size := int64(header.UnknownCount3 * 80)
	fmt.Printf("Unknown Data 4:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData4Size-1, unknownData4Size, header.UnknownCount3)
	offset += unknownData4Size

	unknownData5Size := int64(header.UnknownCount5 * 8)
	fmt.Printf("Unknown Data 5:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData5Size-1, unknownData5Size, header.UnknownCount5)
	offset += unknownData5Size

	unknownData6Size := int64(header.UnknownCount6 * 8)
	fmt.Printf("Unknown Data 6:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData6Size-1, unknownData6Size, header.UnknownCount6)
	offset += unknownData6Size

	importantCitiesSize := int64(header.ImportantCityCount * 4)
	fmt.Printf("Important Cities:   bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+importantCitiesSize-1, importantCitiesSize, header.ImportantCityCount)
	offset += importantCitiesSize

	unknownData7Size := int64(header.UnknownCount9 * 16)
	fmt.Printf("Unknown Data 7:     bytes 0x%04X-0x%04X (0x%04X bytes, %d entries)\n", offset, offset+unknownData7Size-1, unknownData7Size, header.UnknownCount9)
	offset += unknownData7Size

	// Remaining bytes
	remainingBytes := fileLength - offset
	if remainingBytes > 0 {
		fmt.Printf("Remaining Data:     bytes 0x%04X-0x%04X (0x%04X bytes)\n", offset, fileLength-1, remainingBytes)
	}

	fmt.Println(strings.Repeat("=", 60))
}

func getGameModeName(gameMode uint32) string {
	switch gameMode {
	case 1:
		return "Campaign"
	case 2:
		return "Conquest"
	case 6:
		return "Frontier"
	default:
		return "Unknown"
	}
}

func ReadUint8AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	byteData := make([]byte, 1)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint8 from file:", err)
	}

	return int(byteData[0])
}

func ReadUint16AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	byteData := make([]byte, 2)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint16 from file:", err)
	}

	return int(binary.LittleEndian.Uint16(byteData))
}
