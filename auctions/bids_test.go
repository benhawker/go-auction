package auctions

import (
	"testing"
)

func Test_RegisterBid_Success(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	RegisterListing(sellRow)

	bidRow := []string{"17", "8", "BID", "toaster_1", "20.00"}
	returnVal := RegisterBid(bidRow)

	if returnVal != nil {
		t.Errorf("Expected nil, got: %s", returnVal)
	}

	expectedBid := bid{
		shared: shared{
			timestamp: 17,
			userID:    8,
			action:    "BID",
			item:      "toaster_1",
		},
		bidAmount: 20.00,
	}

	listing := listingsRegistry["toaster_1"]

	if *listing.bids[0] != expectedBid {
		t.Errorf("Bid not found on listing. Listing contains: %v", *listing.bids[0])
	}
}

func Test_RegisterBid_Failure_ListingNotFound(t *testing.T) {
	clearRegistries()

	bidRow := []string{"17", "8", "BID", "toaster_1", "20.00"}
	returnVal := RegisterBid(bidRow)

	if returnVal.Error() != "Listing Not Found" {
		t.Errorf("Expected 'Listing Not Found' error, got: %s", returnVal)
	}
}

func Test_RegisterBid_Failure_BidTooLow(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	RegisterListing(sellRow)

	bidRow := []string{"17", "8", "BID", "toaster_1", "20.00"}
	RegisterBid(bidRow)

	bidRow2 := []string{"18", "8", "BID", "toaster_1", "10.00"}
	returnVal := RegisterBid(bidRow2)

	if returnVal.Error() != "Bid of 10.00 needs to be larger than existing max bid of 20.00" {
		t.Errorf("Expected 'Bid too low' error, got: %s", returnVal)
	}
}

func Test_RegisterBid_Failure_AuctionClosed(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	RegisterListing(sellRow)

	bidRow := []string{"21", "8", "BID", "toaster_1", "20.00"}
	returnVal := RegisterBid(bidRow)

	if returnVal.Error() != "This auction closed at: 20. The time now is: 21" {
		t.Errorf("Expected 'Auction closed' error, got: %s", returnVal)
	}
}
