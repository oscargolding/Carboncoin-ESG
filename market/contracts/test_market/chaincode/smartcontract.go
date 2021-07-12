package chaincode

import (
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const CHANNEL = "mychannel"
const PRODUCER = "producer"

type SmartContract struct {
	contractapi.Contract
}

// The result from a query
type PaginatedQueryResult struct {
	Records             []*Offer `json:"records"`
	FetchedRecordsCount int32    `json:"fetchedRecordsCount"`
	Bookmark            string   `json:"bookmark"`
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
		produce := NewProducer(producerId, size)
		return produce.ChainFlush(ctx)
	}
}

// Get the balance of a given producer
func (s *SmartContract) GetBalance(ctx CustomMarketContextInterface, producerId string) (int, error) {
	producer := ctx.GetProducer(producerId)
	if producer == nil {
		return 0, fmt.Errorf("unable to get producer with name: %v", producerId)
	}
	return producer.Tokens, nil
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
		err = exists.ChainFlush(ctx)
		if err != nil {
			return err
		}
		return ctx.CreateOffer(producerId, amountGiven, tokensGiven, offerID)
	}
}

// Get all the offers with the following bookmark, pageSize and string
func (s *SmartContract) GetOffers(ctx CustomMarketContextInterface,
	pageSize int32, bookmark string) (*PaginatedQueryResult, error) {
	queryString := `{"selector":{"docType":"offer", "active": true}}`
	stub := ctx.GetStub()
	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(
		queryString, pageSize, bookmark)
	if err != nil {
		return nil, fmt.Errorf("error with query: %v", err)
	}
	// Wait until the function finishes before closing
	defer resultsIterator.Close()

	offers := make([]*Offer, 0)
	err = ctx.IteratorResults(resultsIterator, func(offer *Offer) {
		offers = append(offers, offer)
	})
	if err != nil {
		return nil, err
	}

	return &PaginatedQueryResult{
		Records:             offers,
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
	if err = seller.IncrementSellable(tokenAmount); err != nil {
		return 0, err
	}
	// Give the tokens to the requesting user
	if err = buyer.IncrementTokens(tokenAmount); err != nil {
		return 0, err
	}
	if usingOffer.IsStale() {
		usingOffer.MakeOfferStale()
	}
	// Now flush
	if err = buyer.ChainFlush(ctx); err != nil {
		return 0, fmt.Errorf("error flushing buyer: %v", err)
	}
	if err = seller.ChainFlush(ctx); err != nil {
		return 0, fmt.Errorf("error flushing seller: %v", err)
	}
	if err = usingOffer.ChainFlush(ctx); err != nil {
		return 0, fmt.Errorf("error flusing offer: %v", err)
	}
	return buyer.Tokens, nil
}
