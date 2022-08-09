package main

import (
	"context"
	"filtool/car"
	"filtool/env"
	"filtool/market"
	"filtool/seal"
	"filtool/snap"
	"fmt"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var err error
	if len(os.Args) > 1 {
		if os.Args[1] == "expiration" {
			if len(os.Args) < 5 {
				printHelp()
				return
			}
			startEpoch, err := strconv.Atoi(os.Args[2])
			if err != nil {
				printHelp()
				return
			}
			endEpoch, err := strconv.Atoi(os.Args[3])
			if err != nil {
				printHelp()
				return
			}
			batchEpoch, err := strconv.Atoi(os.Args[4])
			if err != nil {
				printHelp()
				return
			}
			var minerId *address.Address
			if len(os.Args) >= 6 {
				mi, err := address.NewFromString(os.Args[5])
				if err != nil {
					printHelp()
					return
				}
				minerId = &mi
			}
			seal.GetExpirationPower(context.Background(), abi.ChainEpoch(startEpoch), abi.ChainEpoch(endEpoch), batchEpoch, minerId)
			return
		} else if os.Args[1] == "recovery" {
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
		} else if os.Args[1] == "redis_clean" {
			if len(os.Args) < 5 {
				printHelp()
				return
			}

			minerId, err := strconv.Atoi(os.Args[2])
			if err != nil {
				printHelp()
				return
			}
			min, err := strconv.Atoi(os.Args[3])
			if err != nil {
				printHelp()
				return
			}
			max, err := strconv.Atoi(os.Args[4])
			if err != nil {
				printHelp()
				return
			}
			env.RedisClean(context.Background(), minerId, min, max)
		} else if os.Args[1] == "redis_sectorinfo" {
			if len(os.Args) < 3 {
				printHelp()
				return
			}
			minerId, err := strconv.Atoi(os.Args[2])
			if err != nil {
				printHelp()
				return
			}
			sn, err := strconv.Atoi(os.Args[3])
			if err != nil {
				printHelp()
				return
			}
			env.RedisGetSectorInfo(context.Background(), minerId, sn)
		} else if os.Args[1] == "read_car" {
			file, err := os.Open(os.Args[2])
			if err != nil {
				fmt.Println(err)
				printHelp()
				return
			}
			v, err := car.ReadVersion(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			fmt.Println("version ", v)
		}
	} else {
		printHelp()
	}
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
	fmt.Println("./filtool redis_clean <minerId> <min sector number><max sector number>")
	fmt.Println()
	fmt.Println("./filtool redis_sectorinfo <minerId> <sector number>")
	fmt.Println()
	fmt.Println("./filtool expiration <start epoch> <end epoch> <epoch batch count> <minerId optional>")
	fmt.Println()
}
