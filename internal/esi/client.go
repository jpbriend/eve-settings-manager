package esi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	baseURL        = "https://esi.evetech.net/latest"
	requestTimeout = 10 * time.Second
)

// CharacterInfo represents public character information from ESI.
type CharacterInfo struct {
	Name          string `json:"name"`
	CorporationID int64  `json:"corporation_id"`
	Birthday      string `json:"birthday"`
	Description   string `json:"description"`
	Gender        string `json:"gender"`
	RaceID        int    `json:"race_id"`
	BloodlineID   int    `json:"bloodline_id"`
}

// Client is an ESI API client with caching.
type Client struct {
	httpClient *http.Client
	cache      map[int64]*CharacterInfo
	cacheMu    sync.RWMutex
}

// NewClient creates a new ESI API client.
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: requestTimeout,
		},
		cache: make(map[int64]*CharacterInfo),
	}
}

// GetCharacter fetches public character information by ID.
func (c *Client) GetCharacter(characterID int64) (*CharacterInfo, error) {
	// Check cache first
	c.cacheMu.RLock()
	if info, ok := c.cache[characterID]; ok {
		c.cacheMu.RUnlock()
		return info, nil
	}
	c.cacheMu.RUnlock()

	// Fetch from API
	url := fmt.Sprintf("%s/characters/%d/", baseURL, characterID)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch character %d: %w", characterID, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("character %d not found", characterID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ESI API returned status %d for character %d", resp.StatusCode, characterID)
	}

	var info CharacterInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode character info: %w", err)
	}

	// Cache the result
	c.cacheMu.Lock()
	c.cache[characterID] = &info
	c.cacheMu.Unlock()

	return &info, nil
}

// GetCharacterName is a convenience method to get just the character name.
func (c *Client) GetCharacterName(characterID int64) (string, error) {
	info, err := c.GetCharacter(characterID)
	if err != nil {
		return "", err
	}
	return info.Name, nil
}

// GetCharacterNameOrFallback returns the character name or a fallback string if lookup fails.
func (c *Client) GetCharacterNameOrFallback(characterID int64) string {
	name, err := c.GetCharacterName(characterID)
	if err != nil {
		return fmt.Sprintf("Unknown (%d)", characterID)
	}
	return name
}

// BatchGetCharacterNames fetches names for multiple character IDs concurrently.
func (c *Client) BatchGetCharacterNames(characterIDs []int64) map[int64]string {
	results := make(map[int64]string)
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Limit concurrency to avoid overwhelming the API
	sem := make(chan struct{}, 5)

	for _, id := range characterIDs {
		wg.Add(1)
		go func(charID int64) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			name := c.GetCharacterNameOrFallback(charID)
			mu.Lock()
			results[charID] = name
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	return results
}
