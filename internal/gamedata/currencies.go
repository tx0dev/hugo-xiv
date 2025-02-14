// SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
//
// SPDX-License-Identifier: MIT

package gamedata

import (
	"fmt"

	"github.com/tx0dev/hugo-xiv/internal/types"
)

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
	return nil
}
