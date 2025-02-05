package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"image/png"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"golang.org/x/time/rate"
)

// Config represents the configuration structure
type Config struct {
	XIVApiUrl     string `json:"xivApiUrl"`
	Currencies    string `json:"currencies"`
	Icons         string `json:"icons"`
	ScssVariables string `json:"scssVariables"`
}

// Icons represents the icons.json structure
type Icons map[string]map[string]string

// Currencies represents the currencies.json structure
type Currencies struct {
	Currency  map[string]CurrencyItem   `json:"currency"`
	Tomestone map[string]TomestoneItem  `json:"tomestone"`
	Scrip     map[string]ScripItem      `json:"scrip"`
	Society   map[string]BeastTribeItem `json:"society"`
}

type CurrencyItem struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type TomestoneItem struct {
	Name     string  `json:"name"`
	ID       string  `json:"id"`
	Added    string  `json:"added"`
	Obsolete *string `json:"obsolete"`
}

type ScripItem struct {
	Gatherer string  `json:"gatherer"`
	Crafter  string  `json:"crafter"`
	Added    string  `json:"added"`
	Obsolete *string `json:"obsolete"`
}

type BeastTribeItem struct {
	Name       string `json:"name"`
	SocietyID  string `json:"societyId"`
	CurrencyID string `json:"currencyId"`
}

// Create a global rate limiter
var limiter = rate.NewLimiter(rate.Every(500*time.Millisecond), 1) // 1 request every 500ms

func generateSCSSVariables(currencies Currencies, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating SCSS file: %v", err)
	}
	defer f.Close()

	// Write header
	f.WriteString("// Generated from xiv_currencies.json\n\n")

	// Generate currency maps
	f.WriteString("$xiv-currencies: (\n")

	// Regular currencies
	f.WriteString("    \"currency\": (\n")
	for key := range currencies.Currency {
		className := strings.ReplaceAll(key, "_", "-")
		f.WriteString(fmt.Sprintf("        \"%s\",\n", className))
	}
	for key := range currencies.Society {
		className := strings.ReplaceAll(key, "_", "-")
		f.WriteString(fmt.Sprintf("        \"%s\",\n", className))
	}
	f.WriteString("    ),\n")

	// Tomestones
	f.WriteString("    \"tomestone\": (\n")
	for key := range currencies.Tomestone {
		className := strings.ToLower(key)
		f.WriteString(fmt.Sprintf("        \"%s\",\n", className))
	}
	f.WriteString("    ),\n")

	// Scrips
	f.WriteString("    \"scrip\": (\n")
	for key, scrip := range currencies.Scrip {
		if scrip.Gatherer != "" {
			f.WriteString(fmt.Sprintf("        \"%s-gatherer\",\n", key))
		}
		if scrip.Crafter != "" {
			f.WriteString(fmt.Sprintf("        \"%s-crafter\",\n", key))
		}
	}
	f.WriteString("    ),\n")

	// Close the map
	f.WriteString(");\n")

	return nil
}
func downloadIcon(baseURL, category, iconID, name string) error {
	// Wait for rate limiter
	ctx := context.Background()

	// Create the category directory if it doesn't exist
	dir := filepath.Join("static", "xiv", category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	destPath := filepath.Join(dir, name+".webp")
	if _, err := os.Stat(destPath); err == nil {
		// File already exists, skip
		return nil
	}

	// Generate URL
	folder := iconID[:3] + "000"
	url := fmt.Sprintf("%s/i/%s/%s.png", baseURL, folder, iconID)

	// Download the file
	err := limiter.Wait(ctx)
	if err != nil {
		return fmt.Errorf("rate limiter error: %v", err)
	}
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download icon: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download icon %s: status %d", iconID, resp.StatusCode)
	}

	// Convert the file
	img, err := png.Decode(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to decode PNG: %v", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer out.Close()

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 90)
	if err := webp.Encode(out, img, options); err != nil {
		return fmt.Errorf("failed to encode WebP: %v", err)
	}

	fmt.Printf("   -> Downloaded and converted to WebP: %s/%s.webp\n", category, name)
	return nil
}

func main() {
	// Load config
	configFile, err := os.ReadFile("scripts/config.json")
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		return
	}

	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		fmt.Printf("Error parsing config: %v\n", err)
		return
	}

	fmt.Println("==> Importing data...")
	iconsFile, err := os.ReadFile(config.Icons)
	if err != nil {
		fmt.Printf("Error reading icons.json: %v\n", err)
		return
	}
	var icons Icons
	if err := json.Unmarshal(iconsFile, &icons); err != nil {
		fmt.Printf("Error parsing icons.json: %v\n", err)
		return
	}

	currenciesFile, err := os.ReadFile(config.Currencies)
	if err != nil {
		fmt.Printf("Error reading currencies.json: %v\n", err)
		return
	}
	var currencies Currencies
	if err := json.Unmarshal(currenciesFile, &currencies); err != nil {
		fmt.Printf("Error parsing currencies.json: %v\n", err)
		return
	}

	fmt.Printf("Source URL: %s\n", config.XIVApiUrl)
	fmt.Println("==> Processing icons.json")
	for category, iconMap := range icons {
		for name, iconID := range iconMap {
			if err := downloadIcon(config.XIVApiUrl, category, iconID, name); err != nil {
				fmt.Printf("Error downloading icon: %v\n", err)
			}
		}

	}

	fmt.Println("==> Processing currencies.json")
	// We're not doing the same things as icons since there's a bit more logic...
	for name, curr := range currencies.Currency {
		if err := downloadIcon(config.XIVApiUrl, "currency", curr.ID, name); err != nil {
			fmt.Printf("Error downloading currency icon: %v\n", err)
		}
	}

	for name, tome := range currencies.Tomestone {
		if err := downloadIcon(config.XIVApiUrl, "tomestone", tome.ID, name); err != nil {
			fmt.Printf("Error downloading tomestone icon: %v\n", err)
		}
	}

	// Scrip
	for name, scrip := range currencies.Scrip {
		if scrip.Gatherer != "" {
			if err := downloadIcon(config.XIVApiUrl, "scrip", scrip.Gatherer, name+"-gatherer"); err != nil {
				fmt.Printf("Error downloading gatherer scrip icon: %v\n", err)
			}
		}
		if scrip.Crafter != "" {
			if err := downloadIcon(config.XIVApiUrl, "scrip", scrip.Crafter, name+"-crafter"); err != nil {
				fmt.Printf("Error downloading crafter scrip icon: %v\n", err)
			}
		}
	}

	// Society
	for name, society := range currencies.Society {
		if err := downloadIcon(config.XIVApiUrl, "society", society.SocietyID, name); err != nil {
			fmt.Printf("Error downloading society icon: %v\n", err)
		}
		if err := downloadIcon(config.XIVApiUrl, "currency", society.CurrencyID, name); err != nil {
			fmt.Printf("Error downloading society currency icon: %v\n", err)
		}
	}

	// Generate SCSS variables
	fmt.Println("==> Generating SCSS variables")
	if err := generateSCSSVariables(currencies, config.ScssVariables); err != nil {
		fmt.Printf("Error generating SCSS variables: %v\n", err)
		return
	}

}
