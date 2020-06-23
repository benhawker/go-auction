package auctions

import (
	"errors"
	"strconv"
)

// RegisterBid => For each bid row store the data within the appropriate listing
func RegisterBid(rowData []string) error {
	timestamp, err := strconv.Atoi(rowData[0])
	if err != nil {
		return err
	}

	userID, err := strconv.Atoi(rowData[1])
	if err != nil {
		return err
	}

	bidAmount, err := strconv.ParseFloat(rowData[4], 64)
	if err != nil {
		return err
	}

	b := bid{
		shared: shared{
			timestamp: timestamp,
			userID:    userID,
			action:    rowData[2],
			item:      rowData[3],
		},
		bidAmount: bidAmount,
	}

	if listing, ok := listingsRegistry[b.shared.item]; ok {
		listing.addBid(b)
	} else {
		return errors.New("Listing Not Found")
	}

	return nil
}
