package fileio

import (
	"encoding/binary"
	"log"
	"os"
)

// ReadUint8AtFileOffset reads a uint8 value from a file at the specified offset
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

// ReadUint16AtFileOffset reads a uint16 value from a file at the specified offset
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

func ReadUint32AtFileOffset(inputFilename string, offset int) int {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDONLY, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	byteData := make([]byte, 4)
	if _, err := inputFile.ReadAt(byteData, int64(offset)); err != nil {
		log.Fatal("Failed to read uint32 from file:", err)
	}

	return int(binary.LittleEndian.Uint32(byteData))
}

// WriteUint8AtFileOffset writes a uint8 value at the specified file offset
func WriteUint8AtFileOffset(inputFilename string, offset int, value int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if value >= 256 {
		log.Fatal("Value is too large for uint8")
	}
	if _, err := inputFile.WriteAt([]byte{uint8(value)}, int64(offset)); err != nil {
		log.Fatal("Failed to write uint8 to file:", err)
	}
}

// WriteUint16AtFileOffset writes a uint16 value at the specified file offset
func WriteUint16AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 65536 {
		log.Fatal("Value is too large for uint16")
	}
	byteArrUnitType := make([]byte, 2)
	binary.LittleEndian.PutUint16(byteArrUnitType, uint16(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

// WriteUint32AtFileOffset writes a uint32 value at the specified file offset
func WriteUint32AtFileOffset(inputFilename string, offset int, updatedValue int) {
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state: ", err)
	}

	if updatedValue >= 4294967295 {
		log.Fatal("Value is too large for uint32")
	}
	byteArrUnitType := make([]byte, 4)
	binary.LittleEndian.PutUint32(byteArrUnitType, uint32(updatedValue))
	if _, err := inputFile.WriteAt(byteArrUnitType, int64(offset)); err != nil {
		log.Fatal(err)
	}
}

// WriteDataAtOffset writes data to a file at the specified offset (simple overwrite)
func WriteDataAtOffset(inputFilename string, offset int, newData []byte) {
	// Open file to modify
	inputFile, err := os.OpenFile(inputFilename, os.O_RDWR, 0644)
	defer inputFile.Close()
	if err != nil {
		log.Fatal("Failed to load save state:", err)
	}

	// Simple overwrite - no shifting needed for fixed-size data
	if _, err := inputFile.WriteAt(newData, int64(offset)); err != nil {
		log.Fatal("Failed to write data:", err)
	}
}
