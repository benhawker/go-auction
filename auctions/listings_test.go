package auctions

import (
	"reflect"
	"testing"
)

func Test_RegisterListing_Success(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	returnVal := RegisterListing(sellRow)

	if returnVal != nil {
		t.Errorf("Expected nil, got: %s", returnVal)
	}

	expectedListing := listing{
		shared: shared{
			timestamp: 10,
			userID:    1,
			action:    "SELL",
			item:      "toaster_1",
		},
		reservePrice: 10.00,
		closeTime:    20,
	}

	expiringAuctions := expiryRegistry[20]
	listing := listingsRegistry["toaster_1"]

	if expiringAuctions[0] != expectedListing.item {
		t.Errorf("Listing not found in expiry registry. Registry contains: %v", expiryRegistry)
	}

	if reflect.DeepEqual(listing, expectedListing) {
		t.Errorf("Listing not found in listings registry. Registry contains: %v", listingsRegistry)
	}
}

func Test_RegisterListing_Failure_InvalidUserID(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "User123", "SELL", "toaster_1", "10.00", "20"}
	returnVal := RegisterListing(sellRow)

	if returnVal == nil {
		t.Errorf("Expected error, got: %s", returnVal)
	}

	if len(expiryRegistry) != 0 {
		t.Errorf("Something was added to the registry and it should not have been. Registry now contains: %v", expiryRegistry)
	}
}

func Test_HandleEndingListings_Success_NoEndingAuctionsForTimestamp(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	_ = RegisterListing(sellRow)

	output, _ := HandleEndingListings(1)

	if len(output) != 0 {
		t.Errorf("Expected empty slice, got: %v", output)
	}
}

func Test_HandleEndingListings_Success_EndingAuctionForTimestamp(t *testing.T) {
	clearRegistries()

	sellRow := []string{"10", "1", "SELL", "toaster_1", "10.00", "20"}
	_ = RegisterListing(sellRow)

	output, _ := HandleEndingListings(20)

	if len(output) != 1 {
		t.Errorf("Expected a single output, got: %v", output)
	}
}

func clearRegistries() {
	for k := range expiryRegistry {
		delete(expiryRegistry, k)
	}

	for k := range listingsRegistry {
		delete(listingsRegistry, k)
	}
}
