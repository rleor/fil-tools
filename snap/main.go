package main

import (
	"context"
	"fmt"
	"github.com/filecoin-project/go-jsonrpc"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
	"net/http"
)

var fullnode v1api.FullNode
var storageminernode v0api.StorageMiner

func NewFullNodeRPCV1(ctx context.Context, addr string, requestHeader http.Header) (api.FullNode, jsonrpc.ClientCloser, error) {
	var res v1api.FullNodeStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		api.GetInternalStructs(&res), requestHeader)

	return &res, closer, err
}

func NewStorageMinerRPCV0(ctx context.Context, addr string, requestHeader http.Header, opts ...jsonrpc.Option) (v0api.StorageMiner, jsonrpc.ClientCloser, error) {
	var res v0api.StorageMinerStruct
	closer, err := jsonrpc.NewMergeClient(ctx, addr, "Filecoin",
		api.GetInternalStructs(&res), requestHeader, opts...)

	return &res, closer, err
}

func main() {
	// read config
	config, err := ParseConfig("")
	if err != nil {
		log.Fatalln("read config failed: ", err)
		return
	}
	log.Println("read config success: ", config)

	ctx := context.Background()

	// prepare daemon api
	header := http.Header{"Authorization": []string{"Bearer " + config.Fullnode.Token}}
	fn, closer, err := NewFullNodeRPCV1(ctx, "ws://"+config.Fullnode.Node+"/rpc/v1", header)
	if err != nil {
		log.Fatalf("connecting with fullnode failed: %s", err)
		return
	}
	defer closer()
	fullnode = fn

	// prepare miner api
	header = http.Header{"Authorization": []string{"Bearer " + config.StorageMiner.Token}}
	sm, closer, err := NewStorageMinerRPCV0(ctx, "ws://"+config.StorageMiner.Node+"/rpc/v0", header)
	if err != nil {
		log.Fatalf("connecting with storageminer failed: %s", err)
		return
	}
	defer closer()
	storageminernode = sm

	minerId, err := storageminernode.ActorAddress(ctx)
	if err != nil {
		log.Fatalln("call ActorAddress failed: ", err)
		return
	}
	log.Println("processing miner ", minerId)

	socis, err := fullnode.StateMinerActiveSectors(ctx, minerId, types.EmptyTSK)
	if err != nil {
		log.Fatalln("call StateMinerActiveSectors failed: ", err)
		return
	}

	totalActiveCCSectorCount := 0
	var msg string
	failedSectors := make(map[abi.SectorNumber]string)
	var succeedSectors []abi.SectorNumber
	for _, soci := range socis {
		if soci.SectorKeyCID == nil {
			totalActiveCCSectorCount++
			log.Println("process active cc sector", soci.SectorNumber, "...")

			si, err := storageminernode.SectorsStatus(ctx, soci.SectorNumber, false)
			if err != nil {
				msg = fmt.Sprintf("SectorsStatus failed: %v", err)
				log.Fatalln(msg)
				failedSectors[soci.SectorNumber] = msg
			} else {
				if si.State == "Proving" {
					err := storageminernode.SectorsUpdate(ctx, soci.SectorNumber, api.SectorState("Proving"))
					if err != nil {
						msg = fmt.Sprintf("SectorsUpdate to Proving failed: %v", err)
						log.Fatalln(msg)
						failedSectors[soci.SectorNumber] = msg
					} else {
						//snap-up
						err := storageminernode.SectorMarkForUpgrade(ctx, soci.SectorNumber, true)
						if err != nil {
							msg = fmt.Sprintf("SectorMarkForUpgrade failed: %v", err)
							log.Fatalln(msg)
							failedSectors[soci.SectorNumber] = msg
						} else {
							// TODO: check status again Available
							log.Println("snap-up succeed.")
							succeedSectors = append(succeedSectors, soci.SectorNumber)
						}
					}
				} else {
					msg = fmt.Sprintf("skip, its state is %s", si.State)
					log.Println(msg)
					failedSectors[soci.SectorNumber] = msg
				}
			}

			log.Println()
		}
		if len(succeedSectors) >= config.Limit {
			break
		}
	}

	log.Println("Summary:")
	log.Println("Total Active CC Sector:", totalActiveCCSectorCount)
	log.Println("Snap-up CC Sectors(", len(succeedSectors), "):", succeedSectors)
	log.Println("Unprocessed sectors(", len(failedSectors), "):")
	for k, v := range failedSectors {
		log.Println(k, ":", v)
	}
}
