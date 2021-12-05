# Carboncoin-ESG
Carboncoin - a carbon trading platform for environmental assets using 
ESG reputation. 

## Motivation
A blockchain application to motivate how the use of trustful 
carbon trading can increase market quality through 
a reputation system employing ESG data. A unique fungible asset called 
*Carboncoin* tokenises the right of energy producers to emit carbon.
No on-chain permits required. 

## Architecture
Unique architecture where ESG reputation data is first asynchronously 
submitted and then passed through a pipeline of smart contracts before 
it is validated and weighted by a carbon market.

## Dependencies
* Hyperledger Fabric 2.2
* Node v17
* Go 1.17
* Docker 20
* Linux v5.15
* jq 1.6

## Quick Install Guide from Source
### Install the source code
`Git clone` this repository into the example [fabric samples](https://github.com/hyperledger/fabric-samples) repository. Ideal testing scenario is 
on Linux.
Run the following commands:


### Initialise the chaincode and Fabric

Copy the chain initialisation over with
`cp market/gateway/restart.sh ../test-network/`

Launch Fabric with defined chaincode versions.

`cd ../test-network/` 

`./restart.sh`

### Run the API

Launch the HTTPS API.

`cd ../thesis/market/application/server`

`npm install`

`npm start`

### Run the frontend

`cd ../../client`

`yarn start`

### (Optional) populate the carbon market with dummy data

`cd ../application/server/certifier`

`./create.sh`

## FAQ

*Where are the architecture diagrams held?*

They are found in `diagrams`. 

*Where are the thesis reports and presentations held?*

They are found in `thesisA`, `thesisB` and `thesisC` directories 
respectively. 




