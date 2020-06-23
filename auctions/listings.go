package auctions

import (
	"strconv"
)

// 1. Listings Registry
// The source of listings data (and bids on each listing)
// {
//   itemName: listing1{},
//   tv: listing2{},
//   toaster: listing3{},
// }

// 2. Expiry Registry
// Enables simple look up by itemId to collect all bids
// {
//   timestamp: [toaster, tv],
//   timestamp2: [shoes, book],
// }

type listingsRegistryDefinition map[string]*listing
type expiryRegistryDefinition map[int][]string

var listingsRegistry = make(listingsRegistryDefinition)
var expiryRegistry = make(expiryRegistryDefinition)

// RegisterListing => for each listing (SELL) row store the data in both the listingsRegistry & the expiryRegistry.
func RegisterListing(rowData []string) error {
	timestamp, err := strconv.Atoi(rowData[0])
	if err != nil {
		return err
	}

	userID, err := strconv.Atoi(rowData[1])
	if err != nil {
		return err
	}

	reservePrice, err := strconv.ParseFloat(rowData[4], 64)
	if err != nil {
		return err
	}

	closeTime, err := strconv.Atoi(rowData[5])
	if err != nil {
		return err
	}

	l := listing{
		shared: shared{
			timestamp: timestamp,
			userID:    userID,
			action:    rowData[2],
			item:      rowData[3],
		},
		reservePrice: reservePrice,
		closeTime:    closeTime,
		status:       inProgressStatus,
	}

	err = addListingToRegistries(l)
	if err != nil {
		return err
	}

	return nil
}

// HandleEndingListings => checks against the expiryRegistry for each timestamp and ends the listing if appropriate.
func HandleEndingListings(timestamp int) error {
	if endingListings, ok := expiryRegistry[timestamp]; ok {
		for _, listingName := range endingListings {

			l := listingsRegistry[listingName]

			err := l.endListing()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func addListingToRegistries(l listing) error {
	// This would overwrite an existing listing with the same name.
	// An acknowledged limitation of the design.
	listingsRegistry[l.item] = &l

	// Append to expiryRegistry
	current := expiryRegistry[l.closeTime]
	current = append(current, l.item)
	expiryRegistry[l.closeTime] = current

	return nil
}
