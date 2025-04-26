// SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
//
// SPDX-License-Identifier: MIT

package gamedata

import (
	"context"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/tx0dev/hugo-xiv/internal/config"
	"github.com/tx0dev/hugo-xiv/internal/types"
)

func DownloadIcon(baseURL, category, iconID, name string) error {
	// Create the category directory if it doesn't exist
	dir := filepath.Join("static", "xiv", category)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	// if the IconID is "0", it's a placeholder, skip it
	if iconID == "0" {
		fmt.Printf("   -> Skipping placeholder icon %s\n", name)
		return nil
	}

	destPath := filepath.Join(dir, name+".webp")
	if _, err := os.Stat(destPath); err == nil {
		// File already exists, skip
		return nil
	}

	// Wait for rate limiter (limiter is in fetch.go)
	ctx := context.Background()
	limiter := config.GetLimiter()
	err := limiter.Wait(ctx)
	if err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}

	// Generate URL
	// v2 Format: https://v2.xivapi.com/api/asset?path=ui/icon/060000/060858_hr1.tex&format=png
	// https://v2.xivapi.com/api/asset?path=ui/icon/060000/060858_hr1.tex&format=png
	folder := iconID[:3] + "000"
	url := fmt.Sprintf("%s/asset?path=ui/icon/%s/%s_hr1.tex&format=png", baseURL, folder, iconID)

	// Download the file (fetchData is in fetch.go, but for icons we need image decoding)
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

	fmt.Printf("   -> Fetched: %s/%s.webp\n", category, name)
	return nil
}

func ProcessIcons(config types.Config, icons types.Icons) error {
	fmt.Printf("==> Processing %s\n", config.Icons)
	for category, iconCategoryMap := range icons {
		for iconKey, iconItem := range iconCategoryMap {
			if err := DownloadIcon(config.XIVApiUrl, category, iconItem.ID, iconKey); err != nil {
				fmt.Printf("Error downloading icon: %v\n", err)
				return err // Return error to stop processing if icon download fails
			}
		}
	}
	return nil
}
