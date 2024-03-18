package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type address struct {
	Postcode string `json:"postcode"`
	Line_1   string `json:"line_1"`
	Line_2   string `json:"line_2"`
	City     string `json:"city"`
	County   string `json:"county"`
	Country  string `json:"country"`
}

type Postcode struct {
	Value string `json:"value"`
}

var addresses = []address{
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
	router := gin.Default()

	router.GET("/addresses", getAddresses)
	router.GET("/addresses/:postcode", getAddressByPostcode)
	router.POST("/addresses", postAddresses)

	router.Run("localhost:8080")
}

func getAddresses(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, addresses)
}

// postAddresses adds an address from JSON received in the request body.
func postAddresses(c *gin.Context) {
	var newAddress address

	// Call BindJSON to bind the received JSON to
	// newAddress.
	if err := c.BindJSON(&newAddress); err != nil {
		return
	}

	// Add the new address to the slice.
	addresses = append(addresses, newAddress)
	c.IndentedJSON(http.StatusCreated, newAddress)
}

func getAddressByPostcode(c *gin.Context) {
	postcode := c.Param("postcode")

	for _, a := range addresses {
		if a.Postcode != postcode {
			continue
		}
		c.IndentedJSON(http.StatusOK, a)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "address not found"})
}
