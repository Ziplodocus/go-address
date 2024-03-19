package main

import (
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Address struct {
	Postcode string `json:"postcode"`
	Line_1   string `json:"line_1"`
	Line_2   string `json:"line_2"`
	City     string `json:"city"`
	County   string `json:"county"`
	Country  string `json:"country"`
}

type Postcode struct {
	Value        string `json:"value"`
	last_checked time.Time
}

var addresses = []Address{
	{
		Postcode: "NR1 1NN",
		Line_1:   "40 St. Faiths Lane",
		Line_2:   "",
		City:     "Norwich",
		County:   "Norfolk",
		Country:  "GB",
	},
}

func main() {
	gin.SetMode(os.Getenv("GIN_MODE"))

	InitDatabase()

	router := gin.Default()

	router.GET("/addresses/:postcode", getAddressesByPostcode)

	router.Run(":8080")
}

func getAddressesByPostcode(c *gin.Context) {
	postcode := c.Param("postcode")

	postcode, err := sanitisePostcode(postcode)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "invalid UK postcode format"})
		return
	}

	addresses, err := GetAddresses(postcode)

	if err != nil {
		c.IndentedJSON(http.StatusOK, addresses)
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "addresses not found"})
}

func sanitisePostcode(postcode string) (string, error) {
	// Convert postcode to uppercase
	postcode = strings.ToUpper(postcode)

	// Regular expression to match UK postcode format
	regex := regexp.MustCompile(`^([A-Z]{1,2}\d{1,2}[A-Z]?)\s*(\d[A-Z]{2})$`)
	matches := regex.FindStringSubmatch(postcode)

	if len(matches) == 0 {
		return "", fmt.Errorf("invalid UK postcode format")
	}

	// Concatenate the two parts of the postcode with a space separator
	sanitizedPostcode := fmt.Sprintf("%s %s", matches[1], matches[2])
	return sanitizedPostcode, nil
}
