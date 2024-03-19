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

	template := `{
		"postcode":"` + postcode + `",
		"line_1":"{line_1}",
		"line_2":"{line_2} {line_3} {line_4}",
		"city":"{town_or_city}",
		"county": "{county}",
		"country": "{country}"
	}`
	url := fmt.Sprintf("https://api.getAddress.io/autocomplete/%s?api-key=%s&template=%s", escapedPostcode, os.Getenv("GETADDRESS_API_KEY"), template)

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
