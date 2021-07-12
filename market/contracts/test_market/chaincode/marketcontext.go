// Custom context to make querying the blockchain easier
package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const OFFER_TYPE = "offer"

// The custom object for dealing with the market
type CustomMarketContextInterface interface {
	contractapi.TransactionContextInterface
	CheckProducer(string) (bool, error)
	GetProducer(string) *Producer
	GetUserId() (string, error)
	GetUserType() (string, error)
	CreateOffer(string, int, int, string) error
	GetSellable(string) (int, error)
	IteratorResults(shim.StateQueryIteratorInterface, interface{}) error
	GetResult(string, interface{}) error
}

type CustomMarketContext struct {
	contractapi.TransactionContext
}

// Check if the producer exists or not
func (s *CustomMarketContext) CheckProducer(producerId string) (bool, error) {
	producerJSON, err := s.GetStub().GetState(producerId)
	if err != nil {
		return false, fmt.Errorf("failed to read world state. %v", err)
	}
	if producerJSON == nil {
		return false, nil
	}
	return true, nil
}

// Get the user identifier
func (s *CustomMarketContext) GetUserId() (string, error) {
	id, err := s.GetClientIdentity().GetID()
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
	if len(match) < 2 {
		return "", fmt.Errorf("invalid ID")
	}
	fmt.Println("the id -->", clientIDString)
	fmt.Println("found", match[1], len(match))
	return match[1], nil
}

// Get the user type
func (s *CustomMarketContext) GetUserType() (string, error) {
	userId, err := s.GetUserId()
	if err != nil {
		return "", err
	}
	if userId == "admin" {
		return userId, nil
	}
	att, found, err := s.GetClientIdentity().GetAttributeValue("usertype")
	if !found {
		return "", fmt.Errorf("no `usertype` attribute value found")
	}
	if err != nil {
		return "", err
	}
	return att, nil
}

// Get the specified producer
func (s *CustomMarketContext) GetProducer(producerId string) *Producer {
	producerJSON, err := s.GetStub().GetState(producerId)
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

// Get the sellable tokens available to a producer
func (s *CustomMarketContext) GetSellable(producerId string) (int, error) {
	producer := s.GetProducer(producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get sellable with name: %v", producerId)
	}
	return producer.Sellable, nil
}

// Create the offer on the blockchain for a user
func (s *CustomMarketContext) CreateOffer(producerId string, amount int,
	tokens int, offerID string) error {
	duplicateOfferJSON, err := s.GetStub().GetState(offerID)
	if err != nil || duplicateOfferJSON != nil {
		return fmt.Errorf("offer with id already exists on the market")
	}
	// Create the offer so that it is always active
	offer := Offer{DocType: OFFER_TYPE, Producer: producerId, Amount: amount,
		Tokens: tokens, Active: true, OfferID: offerID}
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("unable to format offer conversion %v", err)
	}
	return s.GetStub().PutState(offerID, offerJSON)
}

func (s *CustomMarketContext) IteratorResults(
	iterator shim.StateQueryIteratorInterface, callBackFuncI interface{}) error {
	funcType := reflect.TypeOf(callBackFuncI)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(callBackFuncI)
	for iterator.HasNext() {
		doc := reflect.New(docType)
		docInterface := doc.Interface()
		queryresult, err := iterator.Next()
		if err != nil {
			return err
		}
		err = json.Unmarshal(queryresult.Value, &docInterface)
		if err != nil {
			return err
		}
		callBackFunc.Call([]reflect.Value{doc})
	}
	return nil
}

func (s *CustomMarketContext) GetResult(name string,
	callBackFuncI interface{}) error {
	funcType := reflect.TypeOf(callBackFuncI)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(callBackFuncI)
	doc := reflect.New(docType)
	docInterface := doc.Interface()
	result, err := s.GetStub().GetState(name)
	if result == nil || err != nil {
		return fmt.Errorf("failed getting state")
	}
	err = json.Unmarshal(result, &docInterface)
	if err != nil {
		return err
	}
	callBackFunc.Call([]reflect.Value{doc})
	return nil
}
