package gamedata

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tx0dev/hugo-xiv/internal/fetch"
	"github.com/tx0dev/hugo-xiv/internal/types"
)

// Data is sourced from Teamcraft, which gets it from the XIVAPI and SaintCoinach.
// It's just already formatted for my purposes. Thanks Teamcraft!

type MapData struct {
	Maps       map[string]types.Map
	Places     map[string]types.Place
	Aetherytes []types.Aetheryte
}

// Fetch or load the game data
func DownloadMapData(ctx context.Context, config types.Config) (MapData, error) {
	fmt.Println("==> Fetching/loading game data")

	var data MapData

	// Define data types and their associated configuration
	dataConfigs := []struct {
		dataType  string
		target    interface{}
		sourceURL string
		cacheName string
	}{
		{
			dataType:  "maps",
			target:    &data.Maps,
			sourceURL: config.Sources.Maps,
			cacheName: "xiv_raw_maps.json",
		},
		{
			dataType:  "places",
			target:    &data.Places,
			sourceURL: config.Sources.Places,
			cacheName: "xiv_raw_places.json",
		},
		{
			dataType:  "aetherytes",
			target:    &data.Aetherytes,
			sourceURL: config.Sources.Aetherytes,
			cacheName: "xiv_raw_aetherytes.json",
		},
	}

	// Process each data type
	for _, dc := range dataConfigs {
		cachePath := filepath.Join("data", dc.cacheName)

		// Check if cache exists
		if _, err := os.Stat(cachePath); os.IsNotExist(err) {
			// Download data if cache doesn't exist
			fmt.Printf("  -> Fetching %s from: %s\n", dc.dataType, dc.sourceURL)
			if err := fetch.FetchData(ctx, dc.sourceURL, dc.target); err != nil {
				fmt.Printf("Error fetching %s: %v\n", dc.dataType, err)
				return MapData{}, err
			}
			// Count fetched items depending on data type
			var count int
			switch v := dc.target.(type) {
			case *map[string]types.Map:
				count = len(*v)
			case *map[string]types.Place:
				count = len(*v)
			case *[]types.Aetheryte:
				count = len(*v)
			}
			fmt.Printf("  -> Fetched %d %s\n", count, dc.dataType)

			// Save to cache
			if err := SaveGameDataJSON(cachePath, dc.target); err != nil {
				return MapData{}, err
			}
		} else {
			// Load from cache if it exists
			fmt.Printf("  -> %s cache already exists, skipping download\n", dc.dataType)
			if err := LoadGameDataJSON(cachePath, dc.target); err != nil {
				fmt.Printf("Error loading %s from file: %v\n", dc.dataType, err)
				return MapData{}, err
			}

			// Count loaded items depending on data type
			var count int
			switch v := dc.target.(type) {
			case *map[string]types.Map:
				count = len(*v)
			case *map[string]types.Place:
				count = len(*v)
			case *[]types.Aetheryte:
				count = len(*v)
			}
			fmt.Printf("  -> Loaded %d %s from file\n", count, dc.dataType)
		}
	}

	fmt.Println("==> Fetch/load complete")
	return data, nil
}

func SaveGameDataJSON(outputPath string, data interface{}) error {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %v", outputPath, err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print JSON
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to file %s: %v", outputPath, err)
	}
	fmt.Printf("  -> Saved data to: %s\n", outputPath)
	return nil
}

func LoadGameDataJSON(inputPath string, data interface{}) error {
	file, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file %s: %v", inputPath, err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(data); err != nil {
		return fmt.Errorf("failed to decode JSON from file %s: %v", inputPath, err)
	}
	fmt.Printf("  -> Loaded data from: %s\n", inputPath)
	return nil
}

// Get the map data and add the place name and aetherytes
func SynthesisMapData(config types.Config, data MapData) (map[string]types.Map, error) {
	fmt.Println("==> Synthesizing map data")

	// Create a reverse lookup for places
	placeMap := make(map[string]types.Place)
	for key, place := range data.Places {
		placeMap[key] = place
	}

	// Create a new map to hold the filtered maps
	filteredMaps := make(map[string]types.Map)

	// Iterate through maps and enrich them with related aetherytes
	for _, gameMap := range data.Maps {
		var relatedAetherytes []types.Aetheryte
		//Update Map place and sub name
		if place, ok := placeMap[fmt.Sprintf("%d", gameMap.PlaceNameID)]; ok {
			gameMap.PlaceName = place
			fmt.Printf("    -> Found Map PlaceName: ID=%d, PlaceNameID=%d, Name=%s\n", gameMap.ID, gameMap.PlaceNameID, place.EN)
		} else {
			fmt.Printf("    -> Map PlaceName NOT Found: ID=%d, PlaceNameID=%d\n", gameMap.ID, gameMap.PlaceNameID)
		}
		if place, ok := placeMap[fmt.Sprintf("%d", gameMap.PlaceNameSubID)]; ok {
			gameMap.SubName = place
			fmt.Printf("    -> Found Map SubName: ID=%d, PlaceNameSubID=%d, Name=%s\n", gameMap.ID, gameMap.PlaceNameSubID, place.EN)
		} else {
			fmt.Printf("    -> Map SubName NOT Found: ID=%d, PlaceNameSubID=%d\n", gameMap.ID, gameMap.PlaceNameSubID)
		}
		// Search for aetherytes
		for _, aetheryte := range data.Aetherytes {
			if aetheryte.MapID == gameMap.ID {
				//Add name to Aetheryte
				if place, ok := placeMap[fmt.Sprintf("%d", aetheryte.NameID)]; ok {
					fmt.Printf("    -> Found Aetheryte Name: ID=%d, NameID=%d, Name=%s\n", aetheryte.ID, aetheryte.NameID, place.EN)
					aetheryte.Name = place
				} else {
					fmt.Printf("    -> Aetheryte Name NOT Found: ID=%d, NameID=%d\n", aetheryte.ID, aetheryte.NameID)
				}
				relatedAetherytes = append(relatedAetherytes, aetheryte)
			}
		}

		// Assign the aetherytes to the map
		gameMap.Aetherytes = relatedAetherytes

		// Generating map ID for the image URL to be used for XIVAPI v2
		// Example: "https://xivapi.com/m/s1f3/s1f3.01.jpg" -> "s1f3/01"
		imageParts := strings.Split(gameMap.Image, "/")
		if len(imageParts) != 0 {
			mapIDPart := strings.TrimSuffix(imageParts[len(imageParts)-1], ".jpg")
			gameMap.Image = strings.Join(strings.Split(mapIDPart, "."), "/")
		}

		// Filter the map, only keeping maps that have aetherytes and are neither housing nor dungeon
		if len(relatedAetherytes) > 0 && !gameMap.Dungeon {
			filteredMaps[gameMap.PlaceName.EN] = gameMap // Add the map to the filtered map
		}
	}

	fmt.Println("==> Synthesis complete")
	return filteredMaps, nil // Return the filtered map
}
