if ```from``` that sends publish storage deal message is lack of funds, deals will stuck at ```Publish``` state.

below command is used to retry publish:
```shell
lotus-miner storage-deals retry-publish <proposal cid>
```

in case there are many such ```Publish``` deals, if you try to retry publish all these deals, you may encounter 'not enough funds' 
since balance is checked against all pending msgs in mempool by gaslimit * gasfeecap which may large than real gas a lot.

fix_market_publish command can fix this. 
it retries publish for ```deals per batch``` at one time, and wait 1 mins for these messages to be on-chain, then continue.
if ```max batches``` is not provided, all ```Publish``` deals will be retry-published.
```shell
./filtool fix_market_publish <deals per batch> <max batches>
```

`deals per batch`: how many retry-publish to be processed in each interval

`max batches`: how many intervals