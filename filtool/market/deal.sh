#!/bin/bash

while true
do
        echo "send deal..."
        dealid=$(lotus client deal --manual-piece-cid=baga6ea4seaqosxcuhahklrgcvtszllj4nuyn6ck7gu3rfftlb4e6pfxmwi3vumy --manual-piece-size=34091302912 bafykbzacedw653luv4o3detvsjdvjxpk6te3kaz7opvtvuxsep7zl3wzjxzyy f033672 0.000000016 518400)
        echo "complete send deal: $dealid"

        state=$(lotus client get-deal $dealid | jq --raw-output '."DealInfo: ".State')
        until [ "$state" = "13" ]
        do
                state=$(lotus client get-deal $dealid | jq --raw-output '."DealInfo: ".State')
          sleep 30
        done

        echo "import car..."
        lotus-miner storage-deals import-data $dealid /root/deals2/test/17G.car
        echo "complete import car"

        sleep 20m
done