package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Prettyletto/post-dude/cmd/ui/views"
)

const baseURL = "http://localhost:8080"

var httpClient = &http.Client{}

func PostCollection(newCollection views.Collection) error {
	url := fmt.Sprintf("%s/collections", baseURL)

	body, err := json.Marshal(newCollection)
	if err != nil {
		return fmt.Errorf("error marshilling the json: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error with the POST request in /collections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code not created: %d", resp.StatusCode)
	}

	return nil
}

func FetchCollections() ([]views.Collection, error) {
	url := fmt.Sprintf("%s/collections", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error with the GET request in /collections: %w", err)
	}
	defer resp.Body.Close()

	var fetchedCollections []views.Collection
	if err := json.NewDecoder(resp.Body).Decode(&fetchedCollections); err != nil {
		return nil, fmt.Errorf("error decoding json: %w", err)
	}

	return fetchedCollections, nil
}

func UpdateCollection(id int, updating views.Collection) error {
	body, err := json.Marshal(updating)
	if err != nil {
		return fmt.Errorf("error marshilling the json: %w", err)
	}

	url := fmt.Sprintf("%s/collections/%d", baseURL, id)

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error with the PUT request in /collections: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error with the PUT request in /collections: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("status code not created: %d", resp.StatusCode)
	}

	return nil
}

func DeleteCollection(id int) error {
	url := fmt.Sprintf("%s/collections/%d", baseURL, id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("error creating request to DELETE in/colections :%w", err)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error with the DELETE request in /collections: %w", err)
	}

	if resp.StatusCode != 202 {
		return fmt.Errorf("status code not created: %d", resp.StatusCode)
	}

	return nil
}
