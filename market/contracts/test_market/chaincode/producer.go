// Custom package for defining the producer object
package chaincode

import (
	"encoding/json"
	"fmt"
)

const SMALL = "small"
const MEDIUM = "medium"
const LARGE = "large"

// Represents a producer being used
type Producer struct {
	ID       string `json:"ID"`
	Tokens   int    `json:"tokens"`
	Sellable int    `json:"sellable"`
}

// Create a new producer
func NewProducer(identification string, size string) *Producer {
	var tokens int
	switch size {
	case MEDIUM:
		tokens = 200
	case LARGE:
		tokens = 300
	default:
		tokens = 100
	}
	return &Producer{ID: identification, Tokens: tokens, Sellable: tokens}
}

// Deduct the amount of sellable for a producer
func (pro *Producer) DeductSellable(tokenDeduction int) error {
	if tokenDeduction < 0 {
		return fmt.Errorf("err negative sellable")
	}
	if pro.Sellable-tokenDeduction < 0 {
		return fmt.Errorf("err producer does not have enough sellable tokens")
	}
	pro.Sellable -= tokenDeduction
	return nil
}

// Increment the amount of sellable tokens
func (pro *Producer) IncrementSellable(tokenIncrease int) error {
	if tokenIncrease < 0 {
		return fmt.Errorf("err: cannot increase by negative amount")
	}
	pro.Sellable += tokenIncrease
	return nil
}

// Deduct the amount of tokens offered
func (pro *Producer) DeductTokens(tokenDeduction int) error {
	if tokenDeduction < 0 {
		return fmt.Errorf("err: cannot deduct negative tokens")
	}
	if pro.Tokens-tokenDeduction < 0 {
		return fmt.Errorf("err producer does not have enough tokens")
	}
	pro.Tokens -= tokenDeduction
	return nil
}

// Increment the amount of tokens
func (pro *Producer) IncrementTokens(tokenIncrease int) error {
	if tokenIncrease < 0 {
		return fmt.Errorf("err: cannot increase tokens by negative amount")
	}
	pro.Tokens += tokenIncrease
	return nil
}

// Flush the producer to the blockchain
//
// WARNING do not use after call
func (prop *Producer) ChainFlush(ctx CustomMarketContextInterface) error {
	jsonProducer, err := json.Marshal(prop)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(prop.ID, jsonProducer)
}
