package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db     *sql.DB
	dbOnce sync.Once
)

func InitDatabase() {
	getDB()
}

func getDB() *sql.DB {
	dbOnce.Do(Connect)
	return db
}

func Connect() {
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASSWORD")
	name := os.Getenv("POSTGRES_DB")
	host := os.Getenv("POSTGRES_HOST")
	port := 5432

	connStr := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		host, port, user, name, pass)

	newDB, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = newDB.Ping()
	if err != nil {
		panic(err)
	}

	db = newDB

	log.Printf("Connected to db")

	// Create the address table
	createAddressesTable(db)
}

func createAddressesTable(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS addresses (
			id SERIAL PRIMARY KEY,
            postcode TEXT PRIMARY KEY,
            lat FLOAT,
            lng FLOAT,
            line_1 TEXT,
            line_2 TEXT,
            city TEXT,
            county TEXT,
            country TEXT
        )`)
	if err != nil {
		return err
	}
	return nil
}

func InsertAddresses(db *sql.DB, addrs []Address) error {
	// Prepare statement for inserting addresses
	stmt, err := db.Prepare(`INSERT INTO addresses (postcode, lat, lng, line_1, line_2, city, county, country)
                                VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Iterate over addresses and insert each one into the database
	for _, addr := range addresses {
		_, err := stmt.Exec(addr.Postcode, addr.Lat, addr.Lng, addr.Line_1, addr.Line_2, addr.City, addr.County, addr.Country)
		if err != nil {
			return err
		}
	}

	return nil
}

func SelectAddressesByPostCode(db *sql.DB, postcode string) ([]Address, error) {
	query := `SELECT postcode, lat, lng, line_1, line_2, city, county, country
                FROM addresses
                WHERE postcode = $1`

	rows, err := db.Query(query, postcode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(&address.Postcode, &address.Lat, &address.Lng, &address.Line_1, &address.Line_2, &address.City, &address.County, &address.Country); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func selectAllAddresses(db *sql.DB) ([]Address, error) {
	query := `SELECT postcode, lat, lng, line_1, line_2, city, county, country
                FROM addresses`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []Address
	for rows.Next() {
		var address Address
		if err := rows.Scan(&address.Postcode, &address.Lat, &address.Lng, &address.Line_1, &address.Line_2, &address.City, &address.County, &address.Country); err != nil {
			return nil, err
		}
		addresses = append(addresses, address)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}
