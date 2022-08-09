package seal

import (
	"context"
	"filtool/util"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/specs-actors/v7/actors/builtin"
	"log"
	"net/http"
	"time"
)

var fullnode3 v1api.FullNode

func GetExpirationPower(ctx context.Context, start abi.ChainEpoch, end abi.ChainEpoch, batch int, minerId *address.Address) {
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
		log.Fatalf("connecting with fullnode2 failed: %s", err)
		return
	}
	defer closer()
	fullnode3 = fn

	startTs, err := fullnode3.ChainGetTipSetByHeight(ctx, start, types.EmptyTSK)
	if err != nil {
		log.Fatalln("ChainHead failed: ", err)
		return
	}

	var minersToCheck []address.Address
	if minerId == nil {
		minersToCheck, err := fullnode3.StateListMiners(ctx, startTs.Key())
		if err != nil {
			log.Fatalln("StateListMiners failed", err)
			return
		}
		log.Println(len(minersToCheck), " miners to check.")
	} else {
		minersToCheck = append(minersToCheck, *minerId)
		log.Printf("%v to check", minerId)
	}

	var skipMinerIds []address.Address
	var processed int

	power32G := abi.SectorSize(34359738368)
	power64G := abi.SectorSize(68719476736)

	var batchResultEpoch []abi.ChainEpoch
	batchResult := make(map[abi.ChainEpoch]abi.StoragePower)
	for temp := start; temp <= end; {
		batchResultEpoch = append(batchResultEpoch, temp)
		batchResult[temp] = abi.NewStoragePower(0)
		temp += abi.ChainEpoch(batch)
	}

	tsStart := build.Clock.Now()
	for _, minerId := range minersToCheck {
		processed++
		// get sector size from sector to save one rpc call.
		//minerInfo, err := fullnode3.StateMinerInfo(ctx, minerId, types.EmptyTSK)
		//if err != nil {
		//	log.Println("error: state miner info failed", err)
		//	skipMinerIds = append(skipMinerIds, minerId)
		//	continue
		//}
		//log.Println("processing ", minerId, " ", minerInfo.SectorSize)

		socis, err := fullnode3.StateMinerActiveSectors(ctx, minerId, startTs.Key())
		if err != nil {
			log.Println("error: StateMinerActiveSectors failed", err)
			skipMinerIds = append(skipMinerIds, minerId)
			continue
		}

		var ss abi.SectorSize

		var expireSectors int
		expirePower := abi.NewStoragePower(0)

		for _, soci := range socis {
			if ss == abi.SectorSize(0) {
				if soci.SealProof == abi.RegisteredSealProof_StackedDrg32GiBV1 ||
					soci.SealProof == abi.RegisteredSealProof_StackedDrg32GiBV1_1 {
					ss = power32G
				} else if soci.SealProof == abi.RegisteredSealProof_StackedDrg64GiBV1 ||
					soci.SealProof == abi.RegisteredSealProof_StackedDrg64GiBV1_1 {
					ss = power64G
				}
			}

			for k, _ := range batchResult {
				if soci.Expiration >= k && soci.Expiration < minEpoch(k+abi.ChainEpoch(batch), end+1) {
					power := QAPowerForSector(ss, soci)
					batchResult[k] = big.Add(batchResult[k], power)
					expireSectors++
					expirePower = big.Add(expirePower, power)
				}
			}
		}
		if processed%1000 == 0 {
			log.Println(processed, "miners has been processed: costs: ", time.Since(tsStart))
			tsStart = build.Clock.Now()
		}
		if expireSectors > 0 {
			log.Println(minerId, ": expire sectors: ", expireSectors, " expire power: ", expirePower)
		}
	}
	log.Println("Summary ", "[", start, ",", end, "]")
	totalPower := abi.NewStoragePower(0)
	for _, e := range batchResultEpoch {
		log.Println("[", e, ",", e+abi.ChainEpoch(batch), ")", ":", batchResult[e])
		totalPower = big.Add(totalPower, batchResult[e])
	}
	log.Println("Total to-expire power: ", totalPower)
}

func minEpoch(x, y abi.ChainEpoch) abi.ChainEpoch {
	if x < y {
		return x
	}
	return y
}

func QAPowerForSector(size abi.SectorSize, sector *miner.SectorOnChainInfo) abi.StoragePower {
	duration := sector.Expiration - sector.Activation
	return QAPowerForWeight(size, duration, sector.DealWeight, sector.VerifiedDealWeight)
}
func QAPowerForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.StoragePower {
	quality := QualityForWeight(size, duration, dealWeight, verifiedWeight)
	return big.Rsh(big.Mul(big.NewIntUnsigned(uint64(size)), quality), builtin.SectorQualityPrecision)
}

func QualityForWeight(size abi.SectorSize, duration abi.ChainEpoch, dealWeight, verifiedWeight abi.DealWeight) abi.SectorQuality {
	// sectorSpaceTime = size * duration
	sectorSpaceTime := big.Mul(big.NewIntUnsigned(uint64(size)), big.NewInt(int64(duration)))
	// totalDealSpaceTime = dealWeight + verifiedWeight
	totalDealSpaceTime := big.Add(dealWeight, verifiedWeight)

	// Base - all size * duration of non-deals
	// weightedBaseSpaceTime = (sectorSpaceTime - totalDealSpaceTime) * QualityBaseMultiplier
	weightedBaseSpaceTime := big.Mul(big.Sub(sectorSpaceTime, totalDealSpaceTime), builtin.QualityBaseMultiplier)
	// Deal - all deal size * deal duration * 10
	// weightedDealSpaceTime = dealWeight * DealWeightMultiplier
	weightedDealSpaceTime := big.Mul(dealWeight, builtin.DealWeightMultiplier)
	// Verified - all verified deal size * verified deal duration * 100
	// weightedVerifiedSpaceTime = verifiedWeight * VerifiedDealWeightMultiplier
	weightedVerifiedSpaceTime := big.Mul(verifiedWeight, builtin.VerifiedDealWeightMultiplier)
	// Sum - sum of all spacetime
	// weightedSumSpaceTime = weightedBaseSpaceTime + weightedDealSpaceTime + weightedVerifiedSpaceTime
	weightedSumSpaceTime := big.Sum(weightedBaseSpaceTime, weightedDealSpaceTime, weightedVerifiedSpaceTime)
	// scaledUpWeightedSumSpaceTime = weightedSumSpaceTime * 2^20
	scaledUpWeightedSumSpaceTime := big.Lsh(weightedSumSpaceTime, builtin.SectorQualityPrecision)

	// Average of weighted space time: (scaledUpWeightedSumSpaceTime / sectorSpaceTime * 10)
	return big.Div(big.Div(scaledUpWeightedSumSpaceTime, sectorSpaceTime), builtin.QualityBaseMultiplier)
}
