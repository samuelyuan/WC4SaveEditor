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
* list-players
* list-player-tiles
* list-cities
* list-units
* list-generals

Write Commands:
* max-money: Sets max currency to 9999.
* max-city-tech: Sets all city tech levels to level 4.
* restore-allies: Heal all of your units and your allies units.
* weaken-enemy: Reduce all enemy units to have 1 health and all enemy cities to have 0 health.
* convert-player: Convert all tiles owned by one player and assign ownership to another player. May crash game.
* convert-tile: Convert one tile and assign ownership to another player. May crash game.
* convert-all-allies: Convert all allied tiles to be your own tiles. May crash game.
* convert-team: Convert all players to be on the same team.
* convert-all-players: Convert all tiles to be your tiles. May crash game.