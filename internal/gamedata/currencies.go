package gamedata

import (
	"fmt"
	"os"
	"strings"

	"github.com/tx0dev/hugo-xiv/internal/types"
)

func GenerateSCSSVariables(currencies types.Currencies, outputPath string) error {
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating SCSS file: %v", err)
	}
	defer f.Close()

	// Write header
	f.WriteString("// Generated from data json\n\n")

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

func ProcessCurrencies(config types.Config, currencies types.Currencies) error {
	fmt.Printf("==> Processing %s\n", config.Currencies)
	// We're not doing the same things as icons since there's a bit more logic...
	for name, curr := range currencies.Currency {
		if err := DownloadIcon(config.XIVApiUrl, "currency", curr.ID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading currency icon: %v\n", err)
			return err
		}
	}

	for name, tome := range currencies.Tomestone {
		if err := DownloadIcon(config.XIVApiUrl, "tomestone", tome.ID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading tomestone icon: %v\n", err)
			return err
		}
	}

	// Scrip
	for name, scrip := range currencies.Scrip {
		if scrip.Gatherer != "" {
			if err := DownloadIcon(config.XIVApiUrl, "scrip", scrip.Gatherer, name+"-gatherer"); err != nil { // downloadIcon is in icons.go
				fmt.Printf("Error downloading gatherer scrip icon: %v\n", err)
				return err
			}
		}
		if scrip.Crafter != "" {
			if err := DownloadIcon(config.XIVApiUrl, "scrip", scrip.Crafter, name+"-crafter"); err != nil { // downloadIcon is in icons.go
				fmt.Printf("Error downloading crafter scrip icon: %v\n", err)
				return err
			}
		}
	}

	// Society
	for name, society := range currencies.Society {
		if err := DownloadIcon(config.XIVApiUrl, "society", society.SocietyID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading society icon: %v\n", err)
			return err
		}
		if err := DownloadIcon(config.XIVApiUrl, "currency", society.CurrencyID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading society currency icon: %v\n", err)
			return err
		}
	}

	// Generate SCSS variables
	fmt.Println("==> Generating SCSS variables")
	if err := GenerateSCSSVariables(currencies, config.ScssVariables); err != nil {
		fmt.Printf("Error generating SCSS variables: %v\n", err)
		return err
	}
	return nil
}
