package auctions

import (
	"fmt"
	"log"
	"sort"
	"strconv"
)

const (
	soldStatus       = "SOLD"
	unsoldStatus     = "UNSOLD"
	inProgressStatus = "IN PROGRESS"
)

// Attributes shared across listing & bid.
type shared struct {
	timestamp int
	userID    int
	strUserID string
	action    string
	item      string
}

type listing struct {
	shared
	reservePrice  float64
	closeTime     int
	lowestBid     float64
	highestBid    float64
	totalBidCount int
	pricePaid     float64
	status        string
	bids          []*bid
}

type bid struct {
	shared
	bidAmount float64
}

// A bid is valid if it:
// - arrives after the auction start time and before or on the closing time.
// - is larger than any previous valid bids submitted by the user.
func (l listing) validateBid(newBid bid) error {
	maxBidAmount := l.maxBid()

	if newBid.bidAmount <= maxBidAmount {
		return fmt.Errorf("Bid of %.2f needs to be larger than existing max bid of %.2f", newBid.bidAmount, maxBidAmount)
	}

	if newBid.timestamp > listingsRegistry[newBid.item].closeTime {
		return fmt.Errorf("This auction closed at: %d. The time now is: %d", listingsRegistry[newBid.item].closeTime, newBid.timestamp)
	}

	return nil
}

func (l listing) maxBid() float64 {
	var max float64

	for _, b := range l.bids {
		if max < b.bidAmount {
			max = b.bidAmount
		}
	}

	return max
}

func (l *listing) getBids() ([]*bid, error) {
	return l.bids, nil
}

func (l *listing) addBid(bid bid) error {
	err := l.validateBid(bid)

	if err != nil {
		return err
	}

	l.bids = append(l.bids, &bid)

	return nil
}

func (l *listing) endListing() error {
	l.strUserID = ""
	l.status = unsoldStatus

	// When there are bids run the following block and assign values to the output struct.
	bids, err := l.getBids()
	if err != nil {
		return err
	}

	if len(bids) > 0 {
		l.sortBidsByAmountAndTimestamp()

		min := l.getMinBid()
		max := l.getMaxBid()
		status := l.calculateStatus(max)

		l.status = status
		l.pricePaid = l.calculatePricePaid(status)
		l.totalBidCount = len(l.bids) // Valid Bids only
		l.highestBid = max
		l.lowestBid = min
	}

	l.printOutput()

	return nil
}

func (l *listing) printOutput() {
	log.Printf("%d|%s|%s|%s|%.2f|%d|%.2f|%.2f \n",
		l.closeTime,
		l.item,
		l.getBuyerUserID(),
		l.status,
		l.pricePaid,
		l.totalBidCount,
		l.highestBid,
		l.lowestBid,
	)
}

func (l *listing) getMinBid() float64 {
	return l.bids[0].bidAmount
}

func (l *listing) getMaxBid() float64 {
	return l.bids[len(l.bids)-1].bidAmount
}

func (l *listing) getBuyerUserID() string {
	if l.status == soldStatus {
		return strconv.Itoa(l.bids[len(l.bids)-1].userID)
	}

	return " " // Return empty string for unsold item (naturally it has no buyer id)
}

func (l *listing) calculateStatus(highestBid float64) string {
	if highestBid >= l.reservePrice {
		return soldStatus
	}

	return unsoldStatus
}

func (l *listing) calculatePricePaid(status string) float64 {
	if status == soldStatus {
		// If there is only a single VALID bid they will pay the reserve price of the auction.
		if len(l.bids) == 1 {
			// NB: This assumes there is ALWAYS a reservePrice
			// An item with no reserve presents issues given we do not consider a 'starting_price' (it is implicitly 0)
			// A sensible addition would be the concept of a `minimumValidBidAmount` or a
			// 'defaultStartingPrice' (e.g. 1) to ensure that a price is always paid in this scenario.
			return l.reservePrice
		}

		// The requirements state:
		//  - At the end of the auction the winner will pay the price of the second highest bidder.
		//  - If there is only a single valid bid they will pay the reserve price of the auction.
		//  - If two bids are received for the same amount then the earliest bid wins the item.
		//
		// In reality the situation is more complex:
		// - At the end of the auction the winner will pay the price of the second highest bidder only if BOTH of these bids exceeded the reserve price.
		//
		// When >= 2 bids
		// - 1 exceeding reserve, 1 below reserve => highest bidder will pay reserve price
		// - Both exceeding or equal to reserve => highester bidder will pay second highest bid
		// - Both below reserve => return 0.00 - UNSOLD item. (NB: this block is not executed in this case)

		secondHighestBidAmount := l.bids[len(l.bids)-2].bidAmount

		// We have already checked that highestBid >= listing.reservePrice in `calculateStatus`
		if secondHighestBidAmount >= l.reservePrice {
			return secondHighestBidAmount
		}

		return l.reservePrice
	}

	// The item was not sold and thus price is 0.00.
	return 0.00
}

// Sorts the bids slice by bids in ascending order (i.e. highest bids last)
// Secondary sort by timestamp in descending order (i.e. earliest bid last)
//
// This allows us to access the last element for the earliest AND highest bid.
// The second last element for the second highest bid whether this is:
// 1. A bid of the same amount but made at a later date
// 2. A bid of a lower amount but made at any time.
func (l *listing) sortBidsByAmountAndTimestamp() error {
	sort.SliceStable(l.bids, func(i, j int) bool {
		if l.bids[i].bidAmount < l.bids[j].bidAmount {
			return true
		}

		if l.bids[i].bidAmount > l.bids[j].bidAmount {
			return false
		}

		return l.bids[i].timestamp > l.bids[j].timestamp
	})

	return nil
}
