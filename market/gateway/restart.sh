#! /bin/sh
./network.sh down
./network.sh up createChannel -c mychannel -ca -s couchdb
./network.sh deployCC -ccn EnergyCertifier -ccp ../thesis/market/contracts/test_certifier -ccl go -cci 'InitLedger'
./network.sh deployCC -ccn Register -ccp ../thesis/market/contracts/register -ccl go 
./network.sh deployCC -ccn basic -ccp ../thesis/market/contracts/test_market -ccl go
