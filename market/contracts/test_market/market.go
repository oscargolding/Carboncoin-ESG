package main

import (
	"log"
	"market/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	marketChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating market chaincode: %v", err)
	}

	if err := marketChaincode.Start(); err != nil {
		log.Panicf("Error starting market chaincode: %v", err)
	}
}
