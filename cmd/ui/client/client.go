package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Prettyletto/post-dude/cmd/ui/collections"
)

func PostCollection(newCollection collections.CollectionItem) error {

	body, err := json.Marshal(newCollection)
	if err != nil {
		return fmt.Errorf("error marshilling the json: %w", err)
	}

	resp, err := http.Post("http://localhost:8080/collections", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error with the POST request in /collections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code not created: %d", resp.StatusCode)
	}

	return nil
}

func FetchCollections() ([]collections.CollectionItem, error) {
	resp, err := http.Get("http://localhost:8080/collections")
	if err != nil {
		return nil, fmt.Errorf("error with the GET request in /collections: %w", err)
	}
	defer resp.Body.Close()

	var fetchedCollections []collections.CollectionItem
	if err := json.NewDecoder(resp.Body).Decode(&fetchedCollections); err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return fetchedCollections, nil
}
