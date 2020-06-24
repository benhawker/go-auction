package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/benhawker/thoughtmachine-auction-v2/auctions"
)

var inputFile string
var defaultInputFile = "inputs/input.txt"

func init() {
	flag.StringVar(&inputFile, "inputFile", defaultInputFile, "Specifies dir of input file")
}

func main() {
	flag.Parse()

	err := parseFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
}

func parseFile(inputFile string) error {
	file, err := os.Open(inputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rowNumber := 0
	for scanner.Scan() {
		rowNumber++

		rowData := strings.Split(scanner.Text(), "|")
		err = handleRow(rowData, rowNumber)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func handleRow(rowData []string, rowNumber int) error {
	timestamp, err := strconv.Atoi(rowData[0])
	if err != nil {
		return fmt.Errorf("No timestamp found on row %d: %s", rowNumber, err)
	}

	if len(rowData) == 1 {
		// no-op (heartbeat)
	} else {
		switch rowData[2] {
		case "SELL":
			err = auctions.RegisterListing(rowData)
			if err != nil {
				return fmt.Errorf("An error occured whilst reading `SELL` input at row %d: %s", rowNumber, err)
			}
		case "BID":
			err = auctions.RegisterBid(rowData)
			if err != nil {
				return fmt.Errorf("An error occured whilst reading `BID` input at row %d: %s", rowNumber, err)
			}
		default:
			return fmt.Errorf("Unidentified line at row %d", rowNumber)
		}
	}

	output, err := auctions.HandleEndingListings(timestamp)
	if err != nil {
		return err
	}

	if len(output) > 0 {
		for _, o := range output {
			printFormattedOutput(o)
		}
	}

	return nil
}

func printFormattedOutput(o auctions.Output) {
	log.Printf("%d|%s|%s|%s|%.2f|%d|%.2f|%.2f \n",
		o.CloseTime,
		o.Item,
		o.BuyerID,
		o.Status,
		o.PricePaid,
		o.TotalBidCount,
		o.HighestBid,
		o.LowestBid,
	)
}
