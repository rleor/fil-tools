#!/bin/bash
minerid="f033672"
index=$1
echo "processing index $index ..."
datapath="$PWD/1k_$index.bin"
echo "generating file $datapath ..."
dd if=/dev/urandom of=$datapath bs=1K count=1 > /dev/null
echo "import data file ..."
datacid=`lotus client import $datapath`
datacid=`echo $datacid | awk '{print $4}'`
echo "datacid=$datacid"
echo "generating car ..."
carpath="$PWD/1k_$index.car"
lotus client generate-car $datapath $carpath
echo "calculating commp"
lotus client commP $carpath
commp=`lotus client commP $carpath | head -n 1 | awk '{print $2}'`
echo "proposing deal"
dealId=`lotus client deal --manual-piece-cid=$commp --manual-piece-size=2032 $datacid $minerid 0.0000000006 518400`
echo "proposal id: $dealId"
echo "lotus-miner storage-deals import-data $dealId /root/deals2/test/1k_$index.car"
cp $PWD/1k_$index.car /root/deals2/test/