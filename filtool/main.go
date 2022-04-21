package main

import (
	"context"
	"filtool/market"
	"filtool/seal"
	"filtool/snap"
	"fmt"
	"github.com/filecoin-project/go-address"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error
	if len(os.Args) > 1 {
		if os.Args[1] == "recovery" {
			if len(os.Args) < 6 {
				printHelp()
				return
			} else {
				minerId, err := address.NewFromString(os.Args[2])
				if err != nil {
					printHelp()
					return
				}
				dl, err := strconv.Atoi(os.Args[3])
				if err != nil {
					printHelp()
					return
				}
				parIdxs := os.Args[4]
				if err != nil {
					printHelp()
					return
				}
				var parInts []uint64
				for _, s := range strings.Split(parIdxs, ",") {
					parIdx, err := strconv.Atoi(s)
					if err != nil {
						printHelp()
						return
					}

					parInts = append(parInts, uint64(parIdx))
				}

				controlAddr := os.Args[5]
				controlAddress, err := address.NewFromString(controlAddr)
				if err != nil {
					log.Fatalln("control address is invalid.")
					printHelp()
					return
				}
				seal.Recovery(context.Background(), minerId, uint64(dl), parInts, controlAddress)
				return
			}
		} else if os.Args[1] == "fix_market_publish" {
			dealsPerBatch := 40
			if len(os.Args) > 2 {
				dealsPerBatch, err = strconv.Atoi(os.Args[2])
				if err != nil {
					log.Println(os.Args[2], "is an invalid int - dealsPerBatch, will use default value 40.")
				}
			}
			batches := 100
			if len(os.Args) > 3 {
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
				limit, err = strconv.Atoi(os.Args[2])
				if err != nil {
					log.Println(os.Args[2], "is an invalid int - limit, will use default value 10.")
				}
			}
			snap.MarkSnaps(context.Background(), limit)
			return
		} else if os.Args[1] == "prenprove" {
			if len(os.Args) < 4 {
				printHelp()
				return
			}

			var startSN int
			var endSN int
			var minerAddress address.Address
			minerAddress, err := address.NewFromString(os.Args[2])
			if err != nil {
				printHelp()
				return
			}

			startSN, err = strconv.Atoi(os.Args[3])
			if err != nil {
				printHelp()
				return
			}

			endSN, err = strconv.Atoi(os.Args[4])
			if err != nil {
				printHelp()
				return
			}
			seal.FindPreButNotProveSectors(context.Background(), minerAddress, startSN, endSN)
			return
		}
	}
	printHelp()
}

func printHelp() {
	fmt.Println("./filtool fix_market_publish <deals_per_batch> <batches>")
	fmt.Println()
	fmt.Println("./filtool mark_cc_available <limit>")
	fmt.Println()
	fmt.Println("./filtool prenprove <minerId> <start_sector_number> <end_sector_number>")
	fmt.Println()
	fmt.Println("./filtool recovery <minerId> <dl> <partitions, separated by ,> <control address>")
	fmt.Println()
}
