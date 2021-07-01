#! /bin/sh
./network.sh down
./network.sh up createChannel -c mychannel -ca
./network.sh deployCC -ccn EnergyCertifier -ccp ../thesis/market/contracts/energy_certifier/ -ccl javascript -cci 'InitLedger'
./network.sh deployCC -ccn basic -ccp ../thesis/market/contracts/carbon_market/ -ccl javascript
