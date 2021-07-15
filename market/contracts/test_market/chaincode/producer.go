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
	Carbon   int    `json:"carbonProduced"`
}

const TOKEN = "%s-tokens"

// Create a new producer
func NewProducer(identification string, size string,
	ctx CustomMarketContextInterface) (*Producer, error) {
	var tokens int
	switch size {
	case MEDIUM:
		tokens = 200
	case LARGE:
		tokens = 300
	default:
		tokens = 100
	}
	highthrough := fmt.Sprintf(TOKEN, identification)
	if err := ctx.UpdateHighThrough(highthrough, "+", tokens); err != nil {
		return nil, err
	}
	return &Producer{ID: identification, Tokens: tokens, Sellable: tokens,
		Carbon: 0}, nil
}

func (pro *Producer) GetTokens(ctx CustomMarketContextInterface) (int, error) {
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	tokens, err := ctx.GetHighThrough(highthrough)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}

func (pro *Producer) AddCarbon(amount int,
	ctx CustomMarketContextInterface) error {
	if amount < 0 {
		return fmt.Errorf("err negative amount of carbon")
	}
	pro.Carbon += amount
	highthrough := fmt.Sprintf("%s-carbon", pro.ID)
	if err := ctx.UpdateHighThrough(highthrough, "+", amount); err != nil {
		return err
	}
	return nil
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
func (pro *Producer) DeductTokens(tokenDeduction int,
	ctx CustomMarketContextInterface) error {
	if tokenDeduction < 0 {
		return fmt.Errorf("err: cannot deduct negative tokens")
	}
	if pro.Tokens-tokenDeduction < 0 {
		return fmt.Errorf("err producer does not have enough tokens")
	}
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	if err := ctx.UpdateHighThrough(highthrough, "-", tokenDeduction); err != nil {
		return err
	}
	pro.Tokens -= tokenDeduction
	return nil
}

// Increment the amount of tokens
func (pro *Producer) IncrementTokens(tokenIncrease int,
	ctx CustomMarketContextInterface) error {
	if tokenIncrease < 0 {
		return fmt.Errorf("err: cannot increase tokens by negative amount")
	}
	pro.Tokens += tokenIncrease
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	if err := ctx.UpdateHighThrough(highthrough, "+", tokenIncrease); err != nil {
		return err
	}
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
