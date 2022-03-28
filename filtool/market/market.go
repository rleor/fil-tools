package market

import (
	"context"
	"filtool/util"
	"github.com/filecoin-project/go-fil-markets/storagemarket"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/ipfs/go-cid"
	"log"
	"net/http"
	"sort"
	"time"
)

var marketnode v0api.StorageMiner

func FixPublishDeals(ctx context.Context, dealsPerBatch int, batches int) {
	// ---- init ----
	// read config
	config, err := util.ParseConfig("")
	if err != nil {
		log.Fatalln("read config failed: ", err)
		return
	}
	log.Println("read config success: ", config)

	// prepare miner api
	header := http.Header{"Authorization": []string{"Bearer " + config.MarketMiner.Token}}
	sm, closer, err := util.NewStorageMinerRPCV0(ctx, "ws://"+config.MarketMiner.Node+"/rpc/v0", header)
	if err != nil {
		log.Fatalf("connecting with marketminer failed: %s", err)
		return
	}
	defer closer()
	marketnode = sm

	// ---------- processing -------------
	deals, err := marketnode.MarketListIncompleteDeals(ctx)
	if err != nil {
		log.Fatalln("call MarketListIncompleteDeals failed: ", err)
		return
	}
	// earlier first.
	sort.Slice(deals, func(i, j int) bool {
		return deals[i].CreationTime.Time().Before(deals[j].CreationTime.Time())
	})

	successDeals := map[cid.Cid]storagemarket.MinerDeal{}
	for _, deal := range deals {
		if deal.State == storagemarket.StorageDealPublish {
			log.Println("retry-publish deal: proposal cid=", deal.ProposalCid, ",create=", deal.CreationTime.Time())
			err := marketnode.MarketRetryPublishDeal(ctx, deal.ProposalCid)
			if err != nil {
				log.Fatalln("error: ", err)
				log.Println()
				continue
			} else {
				log.Println("succeed.")
				successDeals[deal.ProposalCid] = deal
				log.Println()

				if len(successDeals)/dealsPerBatch >= batches {
					break
				}

				if len(successDeals)%dealsPerBatch == 0 {
					log.Println("wait msgs to be on-chain, wait 1 min...")
					time.Sleep(time.Minute)
				}
			}
		}
	}

	time.Sleep(time.Minute)
	log.Println("Summary:")
	log.Println("Processed Deals:")
	updatedDeals, err := marketnode.MarketListIncompleteDeals(ctx)
	if err != nil {
		log.Fatalln("call MarketListIncompleteDeals failed: ", err)
		return
	}
	for _, ud := range updatedDeals {
		if _, ok := successDeals[ud.ProposalCid]; ok {
			if ud.State == storagemarket.StorageDealPublishing {
				log.Println(ud.ProposalCid, "StorageDealPublishing")
			} else {
				log.Println(ud.ProposalCid, ud.State)
			}
		}
	}
}
