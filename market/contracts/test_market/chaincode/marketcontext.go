// Custom context to make querying the blockchain easier
package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const OFFER_TYPE = "offer"
const PRODUCTION_TYPE = "production"

// The custom object for dealing with the market
type CustomMarketContextInterface interface {
	contractapi.TransactionContextInterface
	CheckProducer(string) (bool, error)
	GetProducer(string) *Producer
	GetUserId() (string, error)
	GetUserType() (string, error)
	CreateOffer(string, int, int, string, int) error
	CreateChip(string, int) error
	CreateProduction(string, int, string, string, bool, bool) error
	IteratorResults(shim.StateQueryIteratorInterface, interface{}) error
	GetResult(string, interface{}) error
	UpdateHighThrough(string, string, int) error
	GetHighThrough(string) (int, error)
	OfferStringGenerator(string, bool) string
}

type MarketObject interface {
	InsertContext(CustomMarketContextInterface)
}

type CustomMarketContext struct {
	contractapi.TransactionContext
}

const CHIP = "%s-chip"

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
	// Insert the context
	usingProducer.InsertContext(s)
	return &usingProducer
}

// Get the sellable tokens available to a producer
func (s *CustomMarketContext) GetSellable(producerId string) (int, error) {
	producer := s.GetProducer(producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get sellable with name: %v", producerId)
	}
	sellable, err := producer.GetSellable()
	if err != nil {
		return 0, err
	}
	return sellable, nil
}

// Create a chip to store the direct price offered to a user
func (s *CustomMarketContext) CreateChip(producerId string, amount int) error {
	chip := AmountChip{Amount: amount, Valid: true, Owner: producerId}
	chipJSON, err := json.Marshal(chip)
	if err != nil {
		return fmt.Errorf("unable to format a chip %v", err)
	}
	return s.GetStub().PutState(fmt.Sprintf(CHIP, producerId), chipJSON)
}

// Create the offer on the blockchain for a user
func (s *CustomMarketContext) CreateOffer(producerId string, amount int,
	tokens int, offerID string, carbon int) error {
	duplicateOfferJSON, err := s.GetStub().GetState(offerID)
	if err != nil || duplicateOfferJSON != nil {
		return fmt.Errorf("offer with id already exists on the market")
	}
	// Create the offer so that it is always active
	offer := Offer{DocType: OFFER_TYPE, Producer: producerId, Amount: amount,
		Active: true, OfferID: offerID, CarbonReputation: carbon}
	offer.InsertContext(s)
	err = offer.SetTokens(tokens)
	if err != nil {
		return err
	}
	offerJSON, err := json.Marshal(offer)
	if err != nil {
		return fmt.Errorf("unable to format offer conversion %v", err)
	}
	return s.GetStub().PutState(offerID, offerJSON)
}

// Create production on the blockchain
func (s *CustomMarketContext) CreateProduction(productionId string, carbon int,
	date string, firm string, paid bool, ethical bool) error {
	duplicateProductionJSON, err := s.GetStub().GetState(productionId)
	if err != nil || duplicateProductionJSON != nil {
		return fmt.Errorf("production with id already exists on the market")
	}
	// Create the production
	production := Production{DocType: PRODUCTION_TYPE, ProductionID: productionId,
		Produced: carbon, Date: date, Firm: firm, Paid: paid, Ethical: ethical}
	productionJSON, err := json.Marshal(production)
	if err != nil {
		return fmt.Errorf("unable to format production: %v", err)
	}
	return s.GetStub().PutState(productionId, productionJSON)
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

// Update a high throughput variable - implementation using delta streams
func (s *CustomMarketContext) UpdateHighThrough(name string, op string,
	val int) error {
	if op != "+" && op != "-" {
		return fmt.Errorf("operator %s op is not supported", op)
	}
	txid := s.GetStub().GetTxID()
	compositeValueKey := "varName~op~value~txID"

	// Create the composite key allowing for the future query
	compositeKey, compositeErr := s.GetStub().CreateCompositeKey(
		compositeValueKey, []string{name, op, fmt.Sprint(val), txid})
	if compositeErr != nil {
		return fmt.Errorf("could not create a composite key %v", compositeErr)
	}

	// Save the composite key
	compositePutErr := s.GetStub().PutState(compositeKey, []byte{0x00})
	if compositePutErr != nil {
		return compositePutErr
	}
	return nil
}

// Get a high throughput variable - implementation using delta streams
func (s *CustomMarketContext) GetHighThrough(name string) (int, error) {
	deltaResultsIterator, err := s.GetStub().GetStateByPartialCompositeKey(
		"varName~op~value~txID", []string{name})
	if err != nil {
		return 0, fmt.Errorf("could not retrieve value for %s: %v", name, err)
	}
	defer deltaResultsIterator.Close()

	// Check if the variable even existed
	if !deltaResultsIterator.HasNext() {
		return 0, fmt.Errorf("no variable with name %s exists", name)
	}
	finalVal := 0
	for deltaResultsIterator.HasNext() {
		// Get the next row
		responseRange, nextErr := deltaResultsIterator.Next()
		if nextErr != nil {
			return 0, nextErr
		}
		_, keyParts, splitKeyErr := s.GetStub().SplitCompositeKey(responseRange.Key)
		if splitKeyErr != nil {
			return 0, err
		}

		// Retrieve the delta value and operation
		operation := keyParts[1]
		valueStr := keyParts[2]
		value, err := strconv.Atoi(valueStr)
		fmt.Printf("%d", value)
		if err != nil {
			return 0, fmt.Errorf("err converting to int %v", err)
		}
		switch operation {
		case "+":
			finalVal += value
		case "-":
			finalVal -= value
		default:
			return 0, fmt.Errorf("unexpected operation %s", operation)
		}
	}
	return finalVal, nil
}

func (s *CustomMarketContext) OfferStringGenerator(field string,
	direction bool) string {
	var sorting string
	if direction {
		sorting = `asc`
	} else {
		sorting = `desc`
	}
	queryString := `{"selector":{"docType":"offer", "active": true}%s}`
	if field != "reputation" {
		return fmt.Sprintf(queryString, "")
	} else {
		parameter := `,"sort":[{"docType": "%s"},{"active": "%s"},{"carbonReputation":"%s"}]`
		final := fmt.Sprintf(parameter, sorting, sorting, sorting)
		return fmt.Sprintf(queryString, final)
	}
}
