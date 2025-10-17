package fileio

import "fmt"

// Unit type constants
const (
	UnitTypeLightInfantry      = 1  // Light Infantry
	UnitTypeAssaultInfantry    = 2  // Assault Infantry
	UnitTypeMotorizedInfantry  = 3  // Motorized Infantry
	UnitTypeMechanizedInfantry = 4  // Mechanized Infantry
	UnitTypeCommando           = 5  // Commando
	UnitTypeArmoredCar         = 6  // Armored Car
	UnitTypeLightTank          = 7  // Light Tank
	UnitTypeMediumTank         = 8  // Medium Tank
	UnitTypeHeavyTank          = 9  // Heavy Tank
	UnitTypeSuperTank          = 10 // Super Tank
	UnitTypeFieldArtillery     = 11 // Field Artillery
	UnitTypeHowitzer           = 12 // Howitzer
	UnitTypeRocketArtillery    = 13 // Rocket Artillery
	UnitTypeSuperArtillery     = 14 // Super Artillery
	UnitTypeSubmarine          = 15 // Submarine
	UnitTypeDestroyer          = 16 // Destroyer
	UnitTypeCruiser            = 17 // Cruiser
	UnitTypeCarrier            = 18 // Carrier
	UnitTypeBunker             = 35 // Bunker
	UnitTypeLandFort           = 36 // Land Fort
	UnitTypeCoastalArtillery   = 37 // Coastal Artillery
	UnitTypeRocketLauncher     = 38 // Rocket Launcher
	UnitTypeCity               = 39 // City
)

// GetUnitTypeName returns a human-readable name for the unit type
func GetUnitTypeName(unitType uint8) string {
	switch unitType {
	case UnitTypeLightInfantry:
		return "Light Infantry"
	case UnitTypeAssaultInfantry:
		return "Assault Infantry"
	case UnitTypeMotorizedInfantry:
		return "Motorized Infantry"
	case UnitTypeMechanizedInfantry:
		return "Mechanized Infantry"
	case UnitTypeCommando:
		return "Commando"
	case UnitTypeArmoredCar:
		return "Armored Car"
	case UnitTypeLightTank:
		return "Light Tank"
	case UnitTypeMediumTank:
		return "Medium Tank"
	case UnitTypeHeavyTank:
		return "Heavy Tank"
	case UnitTypeSuperTank:
		return "Super Tank"
	case UnitTypeFieldArtillery:
		return "Field Artillery"
	case UnitTypeHowitzer:
		return "Howitzer"
	case UnitTypeRocketArtillery:
		return "Rocket Artillery"
	case UnitTypeSuperArtillery:
		return "Super Artillery"
	case UnitTypeSubmarine:
		return "Submarine"
	case UnitTypeDestroyer:
		return "Destroyer"
	case UnitTypeCruiser:
		return "Cruiser"
	case UnitTypeCarrier:
		return "Carrier"
	case UnitTypeBunker:
		return "Bunker"
	case UnitTypeLandFort:
		return "Land Fort"
	case UnitTypeCoastalArtillery:
		return "Coastal Artillery"
	case UnitTypeRocketLauncher:
		return "Rocket Launcher"
	case UnitTypeCity:
		return "City"
	default:
		return fmt.Sprintf("Unknown Type %d", unitType)
	}
}
