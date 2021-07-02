package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const CHANNEL = "mychannel"
const SMALL = "small"
const MEDIUM = "medium"
const LARGE = "large"

type SmartContract struct {
	contractapi.Contract
}

type Producer struct {
	ID       string `json:"ID"`
	Tokens   int    `json:"tokens"`
	Sellable int    `json:"sellable"`
}

type Offer struct {
	Producer string `json:"producer"`
	Amount   int    `json:"amount"`
	Tokens   int    `json:"tokens"`
}

// Public Functions //
// The Public exposable on-chain functions - for dealing with producers //

func (s *SmartContract) AddProducer(ctx contractapi.TransactionContextInterface, producerId string) error {
	exists, err := s.CheckProducer(ctx, producerId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("producer with name %v exists", producerId)
	} else {
		// Here we know the producer does not exists
		permissions, err := s.getUserType(ctx)
		if err != nil {
			return fmt.Errorf("error with user attribute %v", err)
		}
		if permissions != "admin" {
			return fmt.Errorf("cannot access unless admin")
		}
		// Call the chaincode
		var matrix [][]byte
		matrix = append(matrix, []byte("FirmSize"))
		matrix = append(matrix, []byte(producerId))
		res := ctx.GetStub().InvokeChaincode("EnergyCertifier", matrix, CHANNEL)
		size := string(res.GetPayload())
		tokenAllocation := s.determineToken(size)
		producer := Producer{ID: producerId, Tokens: tokenAllocation, Sellable: tokenAllocation}
		producerJSON, err := json.Marshal(producer)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutState(producer.ID, producerJSON)
		if err != nil {
			return fmt.Errorf("error putting to world state. %v", err)
		}
		// Return nil, producer was successfully put inside the world state
		return nil
	}
}

// Get the balance of a given producer
func (s *SmartContract) GetBalance(ctx contractapi.TransactionContextInterface, producerId string) (int, error) {
	producer := s.GetProducer(ctx, producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get producer with name: %v", producerId)
	}
	return producer.Tokens, nil
}

// Add an offer for the sale of tokens on chain
func (s *SmartContract) AddOffer(ctx contractapi.TransactionContextInterface,
	producerId string, amountGiven int, tokensGiven int) error {
	exists, err := s.CheckProducer(ctx, producerId)
	if !exists || err != nil {
		return fmt.Errorf("failed to determine the existence of producer")
	} else {
		userType, err := s.getUserType(ctx)
		if err != nil {
			return err
		}
		userId, err := s.getUserId(ctx)
		if err != nil {
			return err
		}
		if userType != "producer" || userId != producerId {
			return fmt.Errorf("carboncoin offers only allowed by valid producers")
		}
		sellable, err := s.GetSellable(ctx, producerId)
		if err != nil {
			return err
		}
		if sellable < tokensGiven {
			return fmt.Errorf("%v does not have enough sellable tokens", producerId)
		}
		err = s.createOffer(ctx, producerId, amountGiven, tokensGiven)
		if err != nil {
			return err
		}
		return s.changeSellableForProducer(ctx, producerId, tokensGiven)
	}
}

// Private class functions //
// Helper functions for the smart contract //

func (s *SmartContract) changeSellableForProducer(
	ctx contractapi.TransactionContextInterface,
	producerId string,
	tokens int) error {
	producer, err := ctx.GetStub().GetState(producerId)
	if producer == nil || err != nil {
		return fmt.Errorf("unabled to change sellable for %v", producerId)
	}
	var producerObj Producer
	err = json.Unmarshal(producer, &producerObj)
	if err != nil {
		return err
	}
	producerObj.Sellable -= tokens
	jsonProducer, err := json.Marshal(producerObj)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(producerId, jsonProducer)
}

func (s *SmartContract) createOffer(ctx contractapi.TransactionContextInterface,
	producerId string, amount int, tokens int) error {
	offer := Offer{Producer: producerId, Amount: amount, Tokens: tokens}
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("%s~%d~%d", producerId, amount, tokens)
	ctx.GetStub().PutState(key, offerJSON)
	return nil
}

// Get the sellable tokens
func (s *SmartContract) GetSellable(ctx contractapi.TransactionContextInterface,
	producerId string) (int, error) {
	producer := s.GetProducer(ctx, producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get sellable with name: %v", producerId)
	}
	return producer.Sellable, nil
}

// Get the specified producer
func (s *SmartContract) GetProducer(ctx contractapi.TransactionContextInterface, producerId string) *Producer {
	producerJSON, err := ctx.GetStub().GetState(producerId)
	if err != nil || producerJSON == nil {
		return nil
	}
	var usingProducer Producer
	err = json.Unmarshal(producerJSON, &usingProducer)
	if err != nil {
		return nil
	}
	return &usingProducer
}

// Check if the producer exists or not
func (s *SmartContract) CheckProducer(ctx contractapi.TransactionContextInterface, producerId string) (bool, error) {
	producerJSON, err := ctx.GetStub().GetState(producerId)
	if err != nil {
		return false, fmt.Errorf("failed to read world state. %v", err)
	}
	if producerJSON == nil {
		return false, nil
	}
	return true, nil
}

func (s *SmartContract) getUserId(ctx contractapi.TransactionContextInterface) (string, error) {
	id, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile("x509::CN=(.*?),.")
	clientID, err := base64.StdEncoding.DecodeString(id)
	clientIDString := string(clientID)
	if err != nil {
		return "", err
	}
	// Find the substring in the regex
	match := re.FindStringSubmatch(clientIDString)
	fmt.Println("the id -->", clientIDString)
	fmt.Println("found", match[1], len(match))
	if len(match) < 2 {
		return "", fmt.Errorf("invalid ID")
	}
	return match[1], nil
}

func (s *SmartContract) getUserType(ctx contractapi.TransactionContextInterface) (string, error) {
	userId, err := s.getUserId(ctx)
	if err != nil {
		return "", err
	}
	if userId == "admin" {
		return userId, nil
	}
	att, found, err := ctx.GetClientIdentity().GetAttributeValue("usertype")
	if !found {
		return "", fmt.Errorf("no `usertype` attribute value found")
	}
	if err != nil {
		return "", err
	}
	return att, nil
}

func (s *SmartContract) determineToken(size string) int {
	switch size {
	case SMALL:
		return 100
	case MEDIUM:
		return 200
	case LARGE:
		return 300
	default:
		return 100
	}
}
