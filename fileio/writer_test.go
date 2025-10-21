package fileio

import (
	"testing"
)

func TestConvertCoordinates(t *testing.T) {
	// Create test unit owner data (3x3 grid)
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	tests := []struct {
		name           string
		coordinateCode int
		gameMode       int
		expectedRow    int
		expectedCol    int
		expectedValid  bool
	}{
		{
			name:           "Valid coordinates - Campaign mode, top-left",
			coordinateCode: 0,
			gameMode:       GameModeCampaign,
			expectedRow:    0,
			expectedCol:    0,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Campaign mode, center",
			coordinateCode: 4,
			gameMode:       GameModeCampaign,
			expectedRow:    1,
			expectedCol:    1,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Campaign mode, bottom-right",
			coordinateCode: 8,
			gameMode:       GameModeCampaign,
			expectedRow:    2,
			expectedCol:    2,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Conquest mode, top-left",
			coordinateCode: 6, // 0 + 2*3
			gameMode:       GameModeConquest,
			expectedRow:    0,
			expectedCol:    0,
			expectedValid:  true,
		},
		{
			name:           "Valid coordinates - Conquest mode, center",
			coordinateCode: 10, // 4 + 2*3
			gameMode:       GameModeConquest,
			expectedRow:    1,
			expectedCol:    1,
			expectedValid:  true,
		},
		{
			name:           "Invalid coordinates - negative",
			coordinateCode: -1,
			gameMode:       GameModeCampaign,
			expectedRow:    -1,
			expectedCol:    2,
			expectedValid:  false,
		},
		{
			name:           "Invalid coordinates - too high for campaign",
			coordinateCode: 9,
			gameMode:       GameModeCampaign,
			expectedRow:    3,
			expectedCol:    0,
			expectedValid:  false,
		},
		{
			name:           "Invalid coordinates - too high for conquest",
			coordinateCode: 15,
			gameMode:       GameModeConquest,
			expectedRow:    3,
			expectedCol:    0,
			expectedValid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			row, col, valid := ConvertCoordinates(tt.coordinateCode, unitOwnerData, tt.gameMode)
			if valid != tt.expectedValid {
				t.Errorf("ConvertCoordinates(%d, %d) valid = %v, want %v", tt.coordinateCode, tt.gameMode, valid, tt.expectedValid)
			}
			if valid {
				if row != tt.expectedRow {
					t.Errorf("ConvertCoordinates(%d, %d) row = %d, want %d", tt.coordinateCode, tt.gameMode, row, tt.expectedRow)
				}
				if col != tt.expectedCol {
					t.Errorf("ConvertCoordinates(%d, %d) col = %d, want %d", tt.coordinateCode, tt.gameMode, col, tt.expectedCol)
				}
			}
		})
	}
}

func TestConvertToCoordinateCode(t *testing.T) {
	// Create test unit owner data (3x3 grid)
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	tests := []struct {
		name         string
		row          int
		col          int
		gameMode     int
		expectedCode int
	}{
		{
			name:         "Campaign mode - top-left",
			row:          0,
			col:          0,
			gameMode:     GameModeCampaign,
			expectedCode: 0,
		},
		{
			name:         "Campaign mode - center",
			row:          1,
			col:          1,
			gameMode:     GameModeCampaign,
			expectedCode: 4,
		},
		{
			name:         "Campaign mode - bottom-right",
			row:          2,
			col:          2,
			gameMode:     GameModeCampaign,
			expectedCode: 8,
		},
		{
			name:         "Conquest mode - top-left",
			row:          0,
			col:          0,
			gameMode:     GameModeConquest,
			expectedCode: 6, // 0 + 2*3
		},
		{
			name:         "Conquest mode - center",
			row:          1,
			col:          1,
			gameMode:     GameModeConquest,
			expectedCode: 10, // 4 + 2*3
		},
		{
			name:         "Conquest mode - bottom-right",
			row:          2,
			col:          2,
			gameMode:     GameModeConquest,
			expectedCode: 14, // 8 + 2*3
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ConvertToCoordinateCode(tt.row, tt.col, unitOwnerData, tt.gameMode)
			if result != tt.expectedCode {
				t.Errorf("ConvertToCoordinateCode(%d, %d, %d) = %d, want %d", tt.row, tt.col, tt.gameMode, result, tt.expectedCode)
			}
		})
	}
}

// Test round-trip conversion (coordinate -> code -> coordinate)
func TestCoordinateRoundTrip(t *testing.T) {
	unitOwnerData := [][]byte{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
	}

	gameModes := []int{GameModeCampaign, GameModeConquest}

	for _, gameMode := range gameModes {
		for row := 0; row < len(unitOwnerData); row++ {
			for col := 0; col < len(unitOwnerData[0]); col++ {
				// Convert row,col to coordinate code
				code := ConvertToCoordinateCode(row, col, unitOwnerData, gameMode)

				// Convert coordinate code back to row,col
				resultRow, resultCol, valid := ConvertCoordinates(code, unitOwnerData, gameMode)

				if !valid {
					t.Errorf("Round-trip conversion failed: (%d,%d) -> %d -> invalid", row, col, code)
					continue
				}

				if resultRow != row || resultCol != col {
					t.Errorf("Round-trip conversion failed: (%d,%d) -> %d -> (%d,%d)",
						row, col, code, resultRow, resultCol)
				}
			}
		}
	}
}

// Test with different map sizes
func TestConvertCoordinatesDifferentSizes(t *testing.T) {
	testCases := []struct {
		name          string
		width         int
		height        int
		code          int
		gameMode      int
		expectedRow   int
		expectedCol   int
		expectedValid bool
	}{
		{
			name:          "1x1 map - Campaign",
			width:         1,
			height:        1,
			code:          0,
			gameMode:      GameModeCampaign,
			expectedRow:   0,
			expectedCol:   0,
			expectedValid: true,
		},
		{
			name:          "1x1 map - Conquest",
			width:         1,
			height:        1,
			code:          2, // 0 + 2*1
			gameMode:      GameModeConquest,
			expectedRow:   0,
			expectedCol:   0,
			expectedValid: true,
		},
		{
			name:          "2x2 map - Campaign",
			width:         2,
			height:        2,
			code:          3,
			gameMode:      GameModeCampaign,
			expectedRow:   1,
			expectedCol:   1,
			expectedValid: true,
		},
		{
			name:          "2x2 map - Conquest",
			width:         2,
			height:        2,
			code:          7, // 3 + 2*2
			gameMode:      GameModeConquest,
			expectedRow:   1,
			expectedCol:   1,
			expectedValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create unit owner data with specified dimensions
			unitOwnerData := make([][]byte, tc.height)
			for i := range unitOwnerData {
				unitOwnerData[i] = make([]byte, tc.width)
			}

			row, col, valid := ConvertCoordinates(tc.code, unitOwnerData, tc.gameMode)
			if valid != tc.expectedValid {
				t.Errorf("ConvertCoordinates(%d, %d) valid = %v, want %v", tc.code, tc.gameMode, valid, tc.expectedValid)
			}
			if valid {
				if row != tc.expectedRow {
					t.Errorf("ConvertCoordinates(%d, %d) row = %d, want %d", tc.code, tc.gameMode, row, tc.expectedRow)
				}
				if col != tc.expectedCol {
					t.Errorf("ConvertCoordinates(%d, %d) col = %d, want %d", tc.code, tc.gameMode, col, tc.expectedCol)
				}
			}
		})
	}
}
