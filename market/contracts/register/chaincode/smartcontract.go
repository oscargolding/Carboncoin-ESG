package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const CARBON_MARKET = "basic"
const CARBON_CHANNEL = "mychannel"

type SmartContract struct {
	contractapi.Contract
}

type ESGCertificate struct {
	Name       string `json:"name"`
	Category   string `json:"category"`
	Min        int    `json:"min"`
	Max        int    `json:"max"`
	Multiplier int    `json:"multiplier"`
}

// Create a simple ESG certificate to put on the chain
func (s *SmartContract) CreateCertificate(
	ctx contractapi.TransactionContextInterface, name string, category string,
	min int, max int, multiplier int) error {
	certificate := ESGCertificate{Name: name, Category: category, Min: min,
		Max: max, Multiplier: multiplier}
	certificateJSON, err := json.Marshal(certificate)
	if err != nil {
		return fmt.Errorf("error: could not marshal certificate %v", err)
	}
	err = ctx.GetStub().PutState(certificate.Name, certificateJSON)
	if err != nil {
		return fmt.Errorf("err: could not put the firm on the chain %v", err)
	}
	return nil
}

// Register that a user was allowed to create a certificate and gain reputation
func (s *SmartContract) RegisterUserCertificate(
	ctx contractapi.TransactionContextInterface, name string, value int,
	firm string, day string, id string) error {
	certificateJSON, err := ctx.GetStub().GetState(name)
	if err != nil {
		return fmt.Errorf("err: unable to get relevant ESG certificate %v", err)
	}
	var certificate ESGCertificate
	err = json.Unmarshal(certificateJSON, &certificate)
	if err != nil {
		return err
	}
	if value < certificate.Min || value > certificate.Max {
		return fmt.Errorf("err: value does not have the right format")
	}
	reputationValue := value * certificate.Multiplier
	var matrix [][]byte
	matrix = append(matrix, []byte("ProducerProduction"))
	matrix = append(matrix, []byte(firm))
	matrix = append(matrix, []byte(fmt.Sprintf("%d", reputationValue)))
	matrix = append(matrix, []byte(day))
	matrix = append(matrix, []byte(id))
	matrix = append(matrix, []byte(certificate.Category))
	matrix = append(matrix, []byte(certificate.Name))
	res := ctx.GetStub().InvokeChaincode(CARBON_MARKET, matrix, CARBON_CHANNEL)
	fmt.Printf("status code ->> %d", res.Status)
	if res.Status != 200 {
		return fmt.Errorf("err: failed calling chaincode: %v", res.Message)
	}
	return nil
}
