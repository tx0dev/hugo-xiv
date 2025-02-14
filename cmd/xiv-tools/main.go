package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tx0dev/hugo-xiv/internal/config"
	"github.com/tx0dev/hugo-xiv/internal/gamedata"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Configuration error: %v\n", err)
		return
	}

	fmt.Println("==> Starting data import...")
	fmt.Printf("Icon Sources: %s\n", cfg.XIVApiUrl)
	fmt.Printf("Game Data sources: %s\n", cfg.Sources)

	icons, err := config.LoadIcons(cfg)
	if err != nil {
		fmt.Printf("Error loading icons: %v\n", err)
		os.Exit(1)
		return
	}

	currencies, err := config.LoadCurrencies(cfg)
	if err != nil {
		fmt.Printf("Error loading currencies: %v\n", err)
		os.Exit(1)
		return
	}
	societies, err := config.LoadSocieties(cfg)
	if err != nil {
		fmt.Printf("Error loading societies: %v\n", err)
		os.Exit(1)
		return
	}

	if err := gamedata.ProcessIcons(cfg, icons); err != nil {
		fmt.Printf("Icon processing error: %v\n", err)
		os.Exit(1)
		return
	}

	if err := gamedata.ProcessCurrencies(cfg, currencies); err != nil {
		fmt.Printf("Currency processing error: %v\n", err)
		os.Exit(1)
		return
	}
	if err := gamedata.ProcessSocieties(cfg, societies); err != nil {
		fmt.Printf("Society processing error: %v\n", err)
		os.Exit(1)
		return
	}

	// Generate SCSS variables
	fmt.Println("==> Generating SCSS variables")
	if err := gamedata.GenerateSCSSVariables(icons, currencies, societies, cfg.ScssVariables); err != nil {
		fmt.Printf("Error generating SCSS variables: %v\n", err)
		os.Exit(1)
		return
	}

	os.Exit(0) // Just aborting before processing maps
	ctx := context.Background()
	mapData, err := gamedata.DownloadMapData(ctx, cfg)
	if err != nil {
		fmt.Printf("Game data processing error: %v\n", err)
		os.Exit(1)
		return
	}
	fmt.Println("==> Data import completed successfully!")

	mapsAndAetherytes, err := gamedata.SynthesisMapData(cfg, mapData)
	if err != nil { // processMaps is in maps.go
		fmt.Printf("Map processing error: %v\n", err)
		os.Exit(1)
		return
	}
	gamedata.SaveGameDataJSON(filepath.Join("data", "xiv_maps.json"), mapsAndAetherytes)
}
