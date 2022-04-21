package seal

import (
	"context"
	"filtool/util"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/build"
	"log"
	"net/http"
	"time"
)

var fullnode v1api.FullNode

func FindPreButNotProveSectors(ctx context.Context, minerId address.Address, startSectorNumber int, endSectorNumber int) {
	// read config
	config, err := util.GetConfig()
	if err != nil {
		log.Fatalln("read config failed: ", err)
		return
	}
	log.Println("read config success: ", config)

	// prepare daemon api
	header := http.Header{"Authorization": []string{"Bearer " + config.Fullnode.Token}}
	fn, closer, err := util.NewFullNodeRPCV1(ctx, "ws://"+config.Fullnode.Node+"/rpc/v1", header)
	if err != nil {
		log.Fatalf("connecting with fullnode failed: %s", err)
		return
	}
	defer closer()
	fullnode = fn

	ts, err := fullnode.ChainHead(ctx)
	if err != nil {
		log.Fatalln("ChainHead failed: ", err)
		return
	}

	const limit = 2880
	log.Printf("start searching sector range[%d,%d]...", startSectorNumber, endSectorNumber)
	tsStart := build.Clock.Now()
	var preSectors []abi.SectorNumber
	for i := startSectorNumber; i <= endSectorNumber; i++ {
		sn := abi.SectorNumber(i)
		//log.Println("processing sector", sn)
		preinfo, err := fullnode.StateSectorPreCommitInfo(ctx, minerId, abi.SectorNumber(i), ts.Key())
		if err == nil {
			if preinfo.PreCommitEpoch < ts.Height()-limit {
				preSectors = append(preSectors, sn)
				log.Println("Found precommit sector", i, ", epoch", preinfo.PreCommitEpoch)
			}
		} else {
			//log.Println(i, ": call StateSectorPreCommitInfo failed: ", err)
		}
	}

	elapsed := time.Since(tsStart)
	log.Println("search complete")
	log.Println("cost", elapsed, "for ", endSectorNumber-startSectorNumber, "sectors.")
	log.Println(len(preSectors), "Sectors has precommit info earlier than", limit, "epochs are:")
	log.Println(preSectors)
}
