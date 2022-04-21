package seal

import (
	"context"
	"filtool/util"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/go-state-types/big"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/api/v1api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/actors"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
	"github.com/filecoin-project/lotus/chain/types"
	"log"
	"net/http"
)

var fullnode2 v1api.FullNode

func Recovery(ctx context.Context, minerId address.Address, dl uint64, parIdxs []uint64, controlAddress address.Address) {
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
	fullnode2 = fn

	params := &miner.DeclareFaultsRecoveredParams{
		Recoveries: []miner.RecoveryDeclaration{},
	}

	for _, parIdx := range parIdxs {
		partitions, err := fullnode2.StateMinerPartitions(ctx, minerId, dl, types.EmptyTSK)
		if err != nil {
			log.Fatalln("call StateMinerPartitions failed: ", err)
			return
		}
		params.Recoveries = append(params.Recoveries, miner.RecoveryDeclaration{
			Deadline:  dl,
			Partition: parIdx,
			Sectors:   partitions[parIdx].FaultySectors,
		})
		faultsCount, err := partitions[parIdx].FaultySectors.Count()
		if err == nil {
			log.Println("dl: ", dl, "partition:", parIdxs, "faults:", faultsCount)
		} else {
			log.Println("count faults failed", err)
		}
	}
	enc, aerr := actors.SerializeParams(params)
	if aerr != nil {
		log.Fatalln("could not serialize declare recoveries parameters: ", aerr)
	}

	msg := &types.Message{
		To:     minerId,
		Method: miner.Methods.DeclareFaultsRecovered,
		Params: enc,
		Value:  types.NewInt(0),
		From:   controlAddress,
	}
	max, _ := big.FromString("1500000000000000000")
	spec := &api.MessageSendSpec{MaxFee: abi.TokenAmount(max)}
	if err := prepareMessage(ctx, minerId, msg, spec); err != nil {
		log.Fatalln(err)
		return
	}

	sm, err := fullnode2.MpoolPushMessage(ctx, msg, &api.MessageSendSpec{MaxFee: abi.TokenAmount(max)})
	if err != nil {
		log.Fatalln("pushing message to mpool: ", err)
		return
	}

	log.Println("declare faults recovered Message CID", "cid", sm.Cid())

	rec, err := fullnode2.StateWaitMsg(context.TODO(), sm.Cid(), build.MessageConfidence, api.LookbackNoLimit, true)
	if err != nil {
		log.Fatalln("declare faults recovered wait error: %w", err)
		return
	}

	if rec.Receipt.ExitCode != 0 {
		log.Fatalln("declare faults recovered wait non-0 exit code: ", rec.Receipt.ExitCode)
		return
	}
	log.Println("Success")
}

func prepareMessage(ctx context.Context, minerId address.Address, msg *types.Message, spec *api.MessageSendSpec) error {
	// (optimal) initial estimation with some overestimation that guarantees
	// block inclusion within the next 20 tipsets.
	gm, err := fullnode2.GasEstimateMessageGas(ctx, msg, spec, types.EmptyTSK)
	if err != nil {
		log.Fatalln("estimating gas", "error", err)
		return nil
	}
	*msg = *gm

	// calculate a more frugal estimation; premium is estimated to guarantee
	// inclusion within 5 tipsets, and fee cap is estimated for inclusion
	// within 4 tipsets.
	minGasFeeMsg := *msg

	minGasFeeMsg.GasPremium, err = fullnode2.GasEstimateGasPremium(ctx, 5, msg.From, msg.GasLimit, types.EmptyTSK)
	if err != nil {
		log.Fatalln("failed to estimate minimum gas premium: ", err)
		minGasFeeMsg.GasPremium = msg.GasPremium
	}

	minGasFeeMsg.GasFeeCap, err = fullnode2.GasEstimateFeeCap(ctx, &minGasFeeMsg, 4, types.EmptyTSK)
	if err != nil {
		log.Fatalln("failed to estimate minimum gas fee cap: ", err)
		minGasFeeMsg.GasFeeCap = msg.GasFeeCap
	}
	return nil
}
