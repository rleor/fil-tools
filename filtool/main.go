package main

import (
	"context"
	"filtool/market"
	"filtool/snap"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "fix_market_publish" {
			dealsPerBatch := 40
			if len(os.Args) > 2 {
				var err error
				dealsPerBatch, err = strconv.Atoi(os.Args[2])
				if err != nil {
					log.Println(os.Args[2], "is an invalid int - dealsPerBatch, will use default value 40.")
				}
			}
			batches := 100
			if len(os.Args) > 3 {
				var err error
				batches, err = strconv.Atoi(os.Args[3])
				if err != nil {
					log.Println(os.Args[3], "is an invalid int - batches, will use default value 100.")
				}
			}
			market.FixPublishDeals(context.Background(), dealsPerBatch, batches)
			return
		} else if os.Args[1] == "mark_cc_available" {
			limit := 10
			if len(os.Args) > 2 {
				var err error
				limit, err = strconv.Atoi(os.Args[2])
				if err != nil {
					log.Println(os.Args[2], "is an invalid int - limit, will use default value 10.")
				}
			}
			snap.MarkSnaps(context.Background(), limit)
			return
		}
	}
	printHelp()
}

func printHelp() {
	// TODO:
}
