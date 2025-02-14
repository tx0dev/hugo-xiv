// SPDX-FileCopyrightText: 2025 Mathieu C. <mathieu@tx0.dev>
//
// SPDX-License-Identifier: MIT

package fetch

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tx0dev/hugo-xiv/internal/config"
)

// Fetches JSON data from a URL and unmarshals it into the provided interface.
func FetchData(ctx context.Context, url string, data interface{}) error {
	limiter := config.GetLimiter()
	err := limiter.Wait(ctx)
	if err != nil {
		return fmt.Errorf("rate limiter error: %w", err)
	}

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("HTTP GET error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status code error: %d for URL: %s", resp.StatusCode, url)
	}

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		return fmt.Errorf("JSON decode error: %w", err)
	}

	return nil
}
