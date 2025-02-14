// SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
//
// SPDX-License-Identifier: MIT

package gamedata

import (
	"fmt"
	"os"
	"strings"

	"github.com/tx0dev/hugo-xiv/internal/types"
)

func GenerateSCSSVariables(icons types.Icons,
	currencies types.Currencies,
	societies types.Societies,
	outputPath string) error {

	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating SCSS file: %v", err)
	}
	defer f.Close()

	// Write header
	f.WriteString("// Generated from data json\n\n")

	// Generate token maps
	f.WriteString("$xiv-tokens: (\n")

	// Regular currencies
	f.WriteString("    \"currency\": (\n")
	for key := range currencies.Currency {
		className := strings.ReplaceAll(key, "_", "-")
		f.WriteString(fmt.Sprintf("        \"%s\",\n", className))
	}
	// Societies are within the currencies
	for key := range societies {
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

	// Generate the icons maps
	f.WriteString("$xiv-icons: (\n")
	for key := range icons {
		f.WriteString(fmt.Sprintf("    \"%s\": (\n", key))
		for icon := range icons[key] {
			className := strings.ReplaceAll(icon, "_", "-")
			f.WriteString(fmt.Sprintf("        \"%s\",\n", className))
		}
		f.WriteString("    ),\n")
	}
	f.WriteString(");\n")

	return nil
}
