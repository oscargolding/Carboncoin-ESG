package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const SMALL = "small"
const MEDIUM = "medium"
const LARGE = "large"

// The type the module is dealing with - smart contract
type SmartContract struct {
	contractapi.Contract
}

// Represents the basic of the application - a firm with ID and size
type Firm struct {
	ID   string `json:"ID"`
	Size string `json:"size"`
}

// Create the energy certifier on chain - selection of firms
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	firms := []Firm{
		{ID: "oscarIndustry", Size: LARGE},
		{ID: "rioTinto", Size: MEDIUM},
		{ID: "smallFirm", Size: SMALL},
		{ID: "largeFirm", Size: LARGE},
		{ID: "oscar@gmail.com", Size: LARGE},
	}

	for _, firm := range firms {
		firmJSON, err := json.Marshal(firm)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(firm.ID, firmJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}
	return nil
}

// Get the size of the firm being handled
func (s *SmartContract) FirmSize(ctx contractapi.TransactionContextInterface, firm string) (string, error) {
	firmJSON, err := ctx.GetStub().GetState(firm)
	if err != nil {
		// Problem with the world state - return nill
		return "", fmt.Errorf("failed to read the world state. %v", err)
	}
	if firmJSON == nil {
		return SMALL, nil
	}
	var usingFirm Firm
	err = json.Unmarshal(firmJSON, &usingFirm)
	if err != nil {
		return "", err
	}
	// Return the size of the given firm
	return usingFirm.Size, nil
}

// Represents logic onchain calling the market
func (s *SmartContract) ReportProduction(
	ctx contractapi.TransactionContextInterface, firm string,
	carbonProduction int, day string, id string) error {
	var matrix [][]byte
	matrix = append(matrix, []byte("ProducerProduction"))
	matrix = append(matrix, []byte(firm))
	matrix = append(matrix, []byte(fmt.Sprintf("%d", carbonProduction)))
	matrix = append(matrix, []byte(day))
	matrix = append(matrix, []byte(id))
	matrix = append(matrix, []byte("Environmental"))
	matrix = append(matrix, []byte("Carbon Dioxide From Hydrogen Production"))
	res := ctx.GetStub().InvokeChaincode("basic", matrix, "mychannel")
	fmt.Printf("status code ->> %d", res.Status)
	if res.Status != 200 {
		return fmt.Errorf("err: failed calling chaincode: %v", res.Message)
	}
	return nil
}
