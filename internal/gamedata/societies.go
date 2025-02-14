// SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
//
// SPDX-License-Identifier: MIT

package gamedata

import (
	"fmt"

	"github.com/tx0dev/hugo-xiv/internal/types"
)

func ProcessSocieties(config types.Config, societies types.Societies) error {
	for name, society := range societies {
		if err := DownloadIcon(config.XIVApiUrl, "society", society.SocietyID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading society icon: %v\n", err)
			return err
		}
		if err := DownloadIcon(config.XIVApiUrl, "currency", society.CurrencyID, name); err != nil { // downloadIcon is in icons.go
			fmt.Printf("Error downloading society currency icon: %v\n", err)
			return err
		}
	}
	return nil
}
