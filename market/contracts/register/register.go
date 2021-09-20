package main

import (
	"log"
	"register/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Entry into the go application - chaincode for the register
func main() {
	firmChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating register-basic chaincode: %v", err)
	}

	if err := firmChaincode.Start(); err != nil {
		log.Panicf("Error starting register-basic chaincode: %v", err)
	}
}
