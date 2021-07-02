// Oscar Golding - UNSW thesis
package main

import (
	"certifier/chaincode"
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Entry into the go application - chaincode
func main() {
	firmChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	if err != nil {
		log.Panicf("Error creating energy-certifier-basic chaincode: %v", err)
	}

	if err := firmChaincode.Start(); err != nil {
		log.Panicf("Error starting energy-certifier-basic chaincode: %v", err)
	}
}
