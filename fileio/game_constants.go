package fileio

// Game mode constants
const (
	GameModeCampaign = 1 // Campaign mode
	GameModeConquest = 2 // Conquest mode
	GameModeFrontier = 6 // Frontier mode
)

// GetGameModeName returns a human-readable name for the game mode
func GetGameModeName(gameMode uint32) string {
	switch gameMode {
	case GameModeCampaign:
		return "Campaign"
	case GameModeConquest:
		return "Conquest"
	case GameModeFrontier:
		return "Frontier"
	default:
		return "Unknown"
	}
}
