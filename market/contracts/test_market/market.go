package main

import (
	"log"
	"market/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Entrypoint
func main() {
	contract := new(chaincode.SmartContract)
	contract.TransactionContextHandler = new(chaincode.CustomMarketContext)
	marketChaincode, err := contractapi.NewChaincode(contract)
	if err != nil {
		log.Panicf("Error creating market chaincode: %v", err)
	}
	if err := marketChaincode.Start(); err != nil {
		log.Panicf("Error starting market chaincode: %v", err)
	}
}
