package chaincode

import (
	"fmt"
	"sort"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const CHANNEL = "mychannel"
const PRODUCER = "producer"
const CERTIFIER = "certifier"
const REGISTER = "register"

type SmartContract struct {
	contractapi.Contract
}

// The result from a query
type PaginatedQueryResult struct {
	Records             []*OfferModel `json:"records"`
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"`
	Bookmark            string        `json:"bookmark"`
}

type NormalQueryResult struct {
	Records []*OfferModel `json:"records"`
}

type PaginatedQueryResultProd struct {
	Records             []*Production `json:"records"`
	FetchedRecordsCount int32         `json:"fetchedRecordsCount"`
	Bookmark            string        `json:"bookmark"`
}

// Public Functions //
// The Public exposable on-chain functions - for dealing with producers //

func (s *SmartContract) AddProducer(ctx CustomMarketContextInterface, producerId string) error {
	exists, err := ctx.CheckProducer(producerId)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("producer with name %v exists", producerId)
	} else {
		// Here we know the producer does not exists
		permissions, err := ctx.GetUserType()
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
		produce, err := NewProducer(producerId, size, ctx)
		if err != nil {
			return err
		}
		return produce.ChainFlush(ctx)
	}
}

// Get the balance of a given producer
func (s *SmartContract) GetBalance(ctx CustomMarketContextInterface, producerId string) (int, error) {
	producer := ctx.GetProducer(producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get producer with name: %v", producerId)
	}
	tokens, err := producer.GetTokens()
	if err != nil {
		return 0, err
	}
	return tokens, nil
}

// Add an offer for the sale of tokens on chain
// the offerID is required to be unique
func (s *SmartContract) AddOffer(ctx CustomMarketContextInterface,
	producerId string, amountGiven int, tokensGiven int, offerID string) error {
	exists := ctx.GetProducer(producerId)
	if exists == nil {
		return fmt.Errorf("failed to determine the existence of producer")
	} else {
		userType, err := ctx.GetUserType()
		if err != nil {
			return err
		}
		userId, err := ctx.GetUserId()
		if err != nil {
			return err
		}
		if userType != "producer" || userId != producerId {
			return fmt.Errorf("carboncoin offers only allowed by valid producers")
		}
		err = exists.DeductSellable(tokensGiven)
		if err != nil {
			return fmt.Errorf("error deducting: %v", err)
		}
		carbon, err := exists.GetCarbon()
		if err != nil {
			return fmt.Errorf("err: getting carbon %v", err)
		}
		return ctx.CreateOffer(producerId, amountGiven, tokensGiven, offerID,
			carbon)
	}
}

// Get the offers contained within a given budget
func (s *SmartContract) GetBudgetOffer(ctx CustomMarketContextInterface,
	reputationMatch bool, target int) (*NormalQueryResult, error) {
	queryString := `{"selector":{"docType":"offer", "active": true}}`
	userId, err := ctx.GetUserId()
	if err != nil {
		return nil, err
	}
	stub := ctx.GetStub()
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("err: error calling blockchain %v", err)
	}
	defer resultsIterator.Close()

	offers := make([]*Offer, 0)
	err = ctx.IteratorResults(resultsIterator, func(offer *Offer) {
		if userId != offer.Producer {
			offers = append(offers, offer)
		}
	})
	if err != nil {
		return nil, fmt.Errorf("err: error with iterator %v", err)
	}
	// Pick an option for sorting
	if reputationMatch {
		sort.Slice(offers[:], func(i, j int) bool {
			return offers[i].CarbonReputation > offers[j].CarbonReputation
		})
	} else {
		sort.Slice(offers[:], func(i, j int) bool {
			return offers[i].Amount < offers[j].Amount
		})
	}
	// Iterate over and greedily pick offers
	runningFound := 0
	foundTarget := false
	returningOffers := make([]*OfferModel, 0)
	for _, offer := range offers {
		offer.InsertContext(ctx)
		model, err := offer.ReturnModel()
		if err != nil {
			return nil, err
		}
		runningFound += model.Tokens
		returningOffers = append(returningOffers, model)
		if runningFound >= target {
			foundTarget = true
			break
		}
	}
	if foundTarget {
		return &NormalQueryResult{
			Records: returningOffers,
		}, nil
	} else {
		emptyOffers := make([]*OfferModel, 0)
		return &NormalQueryResult{
			Records: emptyOffers,
		}, nil
	}
}

// Get the direct price the market should be offering for a carboncoin - price dollar
func (s *SmartContract) GetDirectPrice(ctx CustomMarketContextInterface) (int,
	error) {
	queryString := `{"selector":{"docType":"offer", "active": true}}`
	stub := ctx.GetStub()
	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return 0, fmt.Errorf("error with query: %v", err)
	}
	userType, err := ctx.GetUserType()
	if err != nil {
		return 0, err
	}
	userId, err := ctx.GetUserId()
	if err != nil {
		return 0, err
	}
	if userType != PRODUCER {
		return 0, fmt.Errorf("err: only producers can get a price")
	}
	// Wait until the function terminates before doing the actual close
	defer resultsIterator.Close()

	// Simple max function
	current := 0
	err = ctx.IteratorResults(resultsIterator, func(offer *Offer) {
		if offer.Amount > current {
			current = offer.Amount
		}
	})
	if err != nil {
		return 0, err
	}
	// Return whatever is the max offer + an additional offer of 50 bonus
	val := current + 50
	if err = ctx.CreateChip(userId, val); err != nil {
		return 0, fmt.Errorf("err: cannot create chip %v", err)
	}
	return val, nil
}

func (s *SmartContract) RedeemChip(ctx CustomMarketContextInterface,
	amount int) (int, error) {
	userType, err := ctx.GetUserType()
	if err != nil {
		return 0, err
	}
	if userType != PRODUCER {
		return 0, fmt.Errorf("err: only producers can redeem offer chips")
	}
	userId, err := ctx.GetUserId()
	if err != nil {
		return 0, err
	}
	buyer := ctx.GetProducer(userId)
	var offerChip *AmountChip
	err = ctx.GetResult(fmt.Sprintf(CHIP, userId), func(chip *AmountChip) {
		offerChip = chip
		offerChip.InsertContext(ctx)
	})
	if err != nil {
		return 0, fmt.Errorf("err: error getting offerchip %v", err)
	}
	if !offerChip.Valid {
		return 0, fmt.Errorf("err: the offer chip is no longer valid")
	}
	if err = buyer.IncrementSellable(amount); err != nil {
		return 0, err
	}
	if err = buyer.IncrementTokens(amount); err != nil {
		return 0, err
	}
	// Want to mark the chip is invalid for now
	if err := offerChip.MarkInvalid(); err != nil {
		return 0, err
	}
	tokens, err := buyer.GetTokens()
	if err != nil {
		return 0, err
	}
	return tokens + amount, nil
}

// Get all the offers with the following bookmark, pageSize and string
func (s *SmartContract) GetOffers(ctx CustomMarketContextInterface,
	pageSize int32, bookmark string, field string,
	direction bool, username string) (*PaginatedQueryResult, error) {
	queryString := ctx.OfferStringGenerator(field, direction, username)
	stub := ctx.GetStub()
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(
		queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("error with query: %v", err)
	}
	// Wait until the function finishes before closing
	defer resultsIterator.Close()

	offers := make([]*OfferModel, 0)
	var errorFound error = nil
	err = ctx.IteratorResults(resultsIterator, func(offer *Offer) {
		offer.InsertContext(ctx)
		model, err := offer.ReturnModel()
		if err != nil {
			errorFound = err
		}
		offers = append(offers, model)
	})
	if errorFound != nil {
		return nil, errorFound
	}
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             offers,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

func (s *SmartContract) GetProduction(ctx CustomMarketContextInterface,
	pageSize int32, bookmark string) (*PaginatedQueryResultProd, error) {
	firm, err := ctx.GetUserId()
	if err != nil {
		return nil, err
	}
	queryString := fmt.Sprintf(
		`{"selector":{"docType":"production", "producingFirm": "%s"}}`, firm)
	stub := ctx.GetStub()
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(
		queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("error with query: %v", err)
	}
	// Wait until the function finishes before closing
	defer resultsIterator.Close()
	production := make([]*Production, 0)
	err = ctx.IteratorResults(resultsIterator, func(prod *Production) {
		production = append(production, prod)
	})
	if err != nil {
		return nil, err
	}
	return &PaginatedQueryResultProd{
		Records:             production,
		FetchedRecordsCount: responseMetadata.FetchedRecordsCount,
		Bookmark:            responseMetadata.Bookmark,
	}, nil
}

// Purchase tokens from an offer
func (s *SmartContract) PurchaseOfferTokens(ctx CustomMarketContextInterface,
	purchasingOfferId string, tokenAmount int) (int, error) {
	// Get the producer wanting to buy
	userId, err := ctx.GetUserId()
	if err != nil {
		return 0, err
	}
	userType, err := ctx.GetUserType()
	if err != nil {
		return 0, err
	}
	if userType != PRODUCER {
		return 0, fmt.Errorf("must be a producer to purchase")
	}
	buyer := ctx.GetProducer(userId)
	if buyer == nil {
		return 0, fmt.Errorf("err: buyer could not be found")
	}
	var usingOffer *Offer
	err = ctx.GetResult(purchasingOfferId, func(offer *Offer) {
		usingOffer = offer
		usingOffer.InsertContext(ctx)
	})
	if err != nil {
		return 0, fmt.Errorf("error getting offer: %v", err)
	}
	seller := ctx.GetProducer(usingOffer.Producer)
	if seller == nil {
		return 0, fmt.Errorf("err: seller could not be found")
	}
	if seller.ID == userId {
		return 0, fmt.Errorf("err: cannot purchase tokens from self")
	}
	// Now have the offer, seller and buyer
	// Deduct sellable
	if err = usingOffer.RemoveTokens(tokenAmount); err != nil {
		return 0, err
	}
	if err = seller.DeductTokens(tokenAmount); err != nil {
		return 0, err
	}
	if err = buyer.IncrementSellable(tokenAmount); err != nil {
		return 0, err
	}
	// Give the tokens to the requesting user
	if err = buyer.IncrementTokens(tokenAmount); err != nil {
		return 0, err
	}
	stale, err := usingOffer.IsStale()
	if err != nil {
		return 0, err
	}
	if stale {
		usingOffer.MakeOfferStale()
		// Last entry point for the offer - no more high throughput
		if err = usingOffer.ChainFlush(ctx); err != nil {
			return 0, fmt.Errorf("error flusing offer: %v", err)
		}
	}
	tokens, err := buyer.GetTokens()
	if err != nil {
		return 0, err
	}
	return tokens + tokenAmount, nil
}

// Make an offer stale - no longer active on the blockchain
func (s *SmartContract) MakeOfferStale(ctx CustomMarketContextInterface,
	offerId string) error {
	userId, err := ctx.GetUserId()
	if err != nil {
		return err
	}
	userType, err := ctx.GetUserType()
	if err != nil {
		return err
	}
	if userType != PRODUCER {
		return fmt.Errorf("err: the user must be a producer")
	}
	var usingOffer *Offer
	err = ctx.GetResult(offerId, func(offer *Offer) {
		usingOffer = offer
		usingOffer.InsertContext(ctx)
	})
	if err != nil {
		return fmt.Errorf("err: getting offer %v", err)
	}
	if userId != usingOffer.Producer {
		return fmt.Errorf("err: user does not own the offer and cannot cancel it")
	}
	seller := ctx.GetProducer(userId)
	if seller == nil {
		return fmt.Errorf("err: getting offer seller")
	}
	tokens, err := seller.GetTokens()
	if err != nil {
		return fmt.Errorf("err: getting user tokens %v", err)
	}
	if err = seller.IncrementSellable(tokens); err != nil {
		return err
	}
	usingOffer.MakeOfferStale()
	// Last entry point for the offer - no more high throughput
	if err = usingOffer.ChainFlush(ctx); err != nil {
		return fmt.Errorf("error flusing offer: %v", err)
	}
	return nil
}

// Report the producer production - requires a certifier to call
func (s *SmartContract) ProducerProduction(ctx CustomMarketContextInterface,
	firm string, carbonProduction int,
	day string, id string, prodCategory string, description string) error {
	userType, err := ctx.GetUserType()
	if err != nil {
		return err
	}
	if !(userType == CERTIFIER || userType == "admin" || userType == REGISTER) {
		return fmt.Errorf("err: only certifiers/admins can report carbon production")
	}
	producer := ctx.GetProducer(firm)
	if producer == nil {
		return fmt.Errorf("err: producer does not exist")
	}
	if err = producer.AddCarbon(carbonProduction); err != nil {
		return err
	}
	tokens, err := producer.GetTokens()
	if err != nil {
		return err
	}
	totalCarbonProduction := carbonProduction
	if carbonProduction < 0 {
		totalCarbonProduction = -carbonProduction
	}
	if tokens >= totalCarbonProduction && carbonProduction < 0 {
		if err = producer.DeductTokens(totalCarbonProduction); err != nil {
			return err
		}
		if err = ctx.CreateProduction(id, totalCarbonProduction, day, firm,
			true, false, prodCategory, description); err != nil {
			return err
		}
	} else {
		if err = ctx.CreateProduction(id, totalCarbonProduction, day, firm,
			carbonProduction > 0, carbonProduction > 0, prodCategory,
			description); err != nil {
			return err
		}
	}
	return nil
}

// Pay for the production emitted by a producer
func (s *SmartContract) PayForProduction(ctx CustomMarketContextInterface,
	id string) (int, error) {
	// Make sure only the producer can pay for production
	userType, err := ctx.GetUserType()
	if err != nil {
		return 0, err
	}
	if userType != PRODUCER {
		return 0, fmt.Errorf("err: only producer can pay for production")
	}
	var production *Production
	err = ctx.GetResult(id, func(offer *Production) {
		production = offer
	})
	if err != nil {
		return 0, fmt.Errorf("err: could not get production with id: %s %v", id, err)
	}
	userId, err := ctx.GetUserId()
	if err != nil {
		return 0, err
	}
	if userId != production.Firm {
		return 0, fmt.Errorf("err: do not own carbon debt with id: %s", id)
	}
	producer := ctx.GetProducer(production.Firm)
	if producer == nil {
		return 0, fmt.Errorf("err: could not find the required producer")
	}
	if err = producer.DeductTokens(production.Produced); err != nil {
		return 0, err
	}
	if production.Paid {
		return 0, fmt.Errorf("err: already paid for production")
	}
	production.Paid = true
	if err = production.ChainFlush(ctx); err != nil {
		return 0, err
	}
	tokens, err := producer.GetTokens()
	if err != nil {
		return 0, err
	}
	return tokens - production.Produced, nil
}
