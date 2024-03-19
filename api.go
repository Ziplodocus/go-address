package main

func GetAddresses(postcode string) ([]Address, error) {
	db = getDB()

	addresses, err := SelectAddressesByPostCode(db, postcode)

	if err == nil {
		return addresses, nil
	}

	addresses, err = FetchAddresses(postcode)

	if err != nil {
		InsertAddresses(db, addresses)
	}

	return addresses, err
}
