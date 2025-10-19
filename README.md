# WC4SaveEditor

## Save File Location

Before using this tool, you need to locate your World Conqueror 4 save file:

**Windows (Microsoft Store version):**
```
C:\Users\<YourUserName>\AppData\Local\Packages\EasyTech.WorldConqueror4_<PackageIdentifier>\LocalState
```

Replace `<YourUserName>` with your actual Windows username. The `<PackageIdentifier>` is a unique string that may vary.

## Documentation

- [File Format Documentation](FILE_FORMAT.md) - Detailed description of the World Conqueror 4 save file format

## How to Use

There are various commands to modify the save file. 

Make sure you quit your current game and go to the main menu before overwriting the save file. If you overwrite the file while the game is still in progress, the game will overwrite the file when you leave and none of your new changes will apply.

Read Commands:
* list-players: Display all players with country information, team membership, and unit counts
* list-cities: Show all cities with their positions and ownership details
* list-units: Display all military units grouped by owner with detailed unit information
* list-units-by-map: Analyze unit distribution across the map
* list-tiles: Show city tile ownership analysis with territory control statistics
* list-generals: Display all units with assigned generals
* list-landmines: Show all landmines grouped by owner with position and health data
* list-teams: Display team analysis with player counts, territories, and alliance structure

Write Commands:
* max-money: Sets max currency to 9999.
* max-city-tech: Sets all city tech levels to level 4.
* heal-allies: Heal all of your units and your allies units.
* weaken-enemies: Minimize all enemy money to 0, reduce all enemy city tech to 0, and reduce all enemy units to have 1 health and all enemy cities to have 0 health.
* conquer-player: Conquer all territories from a specific player (interactive). May crash game.
* transfer-tile: Convert one tile and assign ownership to another player (interactive). May crash game.
* conquer-allies: Conquer all allied territories and make them yours. May crash game.
* unite-team: Convert all players to be on the same team.
* conquer-all: Convert all tiles to be your tiles. May crash game.

## Usage Examples

```bash
# Template: WC4SaveEditor.exe -input <savefile> -command <command>
```