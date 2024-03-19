package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

func FetchAddresses(postcode string) ([]Address, error) {
	escapedPostcode := url.QueryEscape(postcode)

	url := fmt.Sprintf("https://api.getAddress.io/v2/uk/%s?api-key=%s", escapedPostcode, os.Getenv("GETADDRESS_API_KEY"))

	// Make GET request
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decode JSON response
	var addresses []Address
	if err := json.NewDecoder(response.Body).Decode(&addresses); err != nil {
		return nil, err
	}

	return addresses, nil
}
