package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/tx0dev/hugo-xiv/internal/types"
	"golang.org/x/time/rate"
)

var config types.Config

// loadConfig reads and parses the configuration file.
func LoadConfig() (types.Config, error) {
	configFile, err := os.ReadFile("config.json") // Assuming config.json is in the scripts directory
	if err != nil {
		return types.Config{}, fmt.Errorf("error reading config: %w", err)
	}

	if err := json.Unmarshal(configFile, &config); err != nil { // Unmarshal directly into globalConfig
		return types.Config{}, fmt.Errorf("error parsing config: %w", err)
	}
	return config, nil
}

// loadIcons reads and parses the icons JSON file.
func LoadIcons(config types.Config) (types.Icons, error) {
	iconsFile, err := os.ReadFile(config.Icons)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", config.Icons, err)
	}
	var icons types.Icons
	if err := json.Unmarshal(iconsFile, &icons); err != nil {
		return nil, fmt.Errorf("error parsing icons.json: %w", err)
	}
	return icons, nil
}

// loadCurrencies reads and parses the currencies JSON file.
func LoadCurrencies(config types.Config) (types.Currencies, error) {
	currenciesFile, err := os.ReadFile(config.Currencies)
	if err != nil {
		return types.Currencies{}, fmt.Errorf("error reading %s: %w", config.Currencies, err)
	}
	var currencies types.Currencies
	if err := json.Unmarshal(currenciesFile, &currencies); err != nil {
		return types.Currencies{}, fmt.Errorf("error parsing currencies.json: %w", err)
	}
	return currencies, nil
}

func LoadSocieties(config types.Config) (types.Societies, error) {
	societiesFile, err := os.ReadFile(config.Societies)
	if err != nil {
		return types.Societies{}, fmt.Errorf("error reading %s: %w", config.Societies, err)
	}
	var societies types.Societies
	if err := json.Unmarshal(societiesFile, &societies); err != nil {
		return types.Societies{}, fmt.Errorf("error parsing societies.json: %w", err)
	}
	return societies, nil
}

// Return a limiter object
func GetLimiter() *rate.Limiter {
	return rate.NewLimiter(rate.Every(time.Duration(config.RateLimit)*time.Millisecond), 1) // 1 request every 500ms
}
