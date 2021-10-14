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
	ID          string                       `json:"ID"`
	Environment int                          `json:"environment"`
	Social      int                          `json:"int"`
	Governance  int                          `json:"governance"`
	Ctx         CustomMarketContextInterface `json:"-"`
}

const TOKEN = "%s-tokens"
const SELLABLE = "%s-sellable"
const CARBON = "%s-carbon"

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
	highSell := fmt.Sprintf(SELLABLE, identification)
	highcarbon := fmt.Sprintf(CARBON, identification)
	if err := ctx.UpdateHighThrough(highthrough, "+", tokens); err != nil {
		return nil, err
	}
	if err := ctx.UpdateHighThrough(highSell, "+", tokens); err != nil {
		return nil, err
	}
	// Producer starts on a clean slate
	if err := ctx.UpdateHighThrough(highcarbon, "+", 0); err != nil {
		return nil, err
	}
	return &Producer{ID: identification, Ctx: ctx, Environment: 0, Social: 0,
		Governance: 0}, nil
}

func (pro *Producer) EnforceCtx() error {
	if pro.Ctx == nil {
		return fmt.Errorf("err: the blockchain context not set for producer")
	}
	return nil
}

func (pro *Producer) InsertContext(ctx CustomMarketContextInterface) {
	pro.Ctx = ctx
}

func (pro *Producer) GetTokens() (int, error) {
	// As a start - want to enforce the connection to the blockchain
	if err := pro.EnforceCtx(); err != nil {
		return 0, err
	}
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	tokens, err := pro.Ctx.GetHighThrough(highthrough)
	if err != nil {
		return 0, err
	}
	return tokens, nil
}

// Get the carbon associated with a user
func (pro *Producer) GetCarbon() (int, error) {
	if err := pro.EnforceCtx(); err != nil {
		return 0, err
	}
	highthrough := fmt.Sprintf(CARBON, pro.ID)
	carbon, err := pro.Ctx.GetHighThrough(highthrough)
	if err != nil {
		return 0, err
	}
	return carbon, nil
}

// Add carbon production to the blockchain
func (pro *Producer) AddCarbon(amount int, category string) error {
	if err := pro.EnforceCtx(); err != nil {
		return err
	}
	sign := "+"
	highAmount := amount
	if amount < 0 {
		sign = "-"
		highAmount = -amount
	}
	highthrough := fmt.Sprintf(CARBON, pro.ID)
	if err := pro.Ctx.UpdateHighThrough(highthrough, sign, highAmount); err != nil {
		return err
	}
	switch category {
	case "Environmental":
		pro.Environment += amount
	case "Social":
		pro.Social += amount
	case "Governance":
		pro.Governance += amount
	}
	return nil
}

// Get the sellable tokens associated with a user
func (pro *Producer) GetSellable() (int, error) {
	if err := pro.EnforceCtx(); err != nil {
		return 0, err
	}
	sellableTokens, err := pro.Ctx.GetHighThrough(fmt.Sprintf(SELLABLE, pro.ID))
	if err != nil {
		return 0, nil
	}
	return sellableTokens, nil
}

// Deduct the amount of sellable for a producer
func (pro *Producer) DeductSellable(tokenDeduction int) error {
	if err := pro.EnforceCtx(); err != nil {
		return err
	}
	if tokenDeduction < 0 {
		return fmt.Errorf("err negative sellable")
	}
	sellableTokens, err := pro.Ctx.GetHighThrough(fmt.Sprintf(SELLABLE, pro.ID))
	if err != nil {
		return err
	}
	if sellableTokens-tokenDeduction < 0 {
		return fmt.Errorf("err producer does not have enough sellable tokens")
	}
	err = pro.Ctx.UpdateHighThrough(fmt.Sprintf(SELLABLE, pro.ID), "-",
		tokenDeduction)
	if err != nil {
		return err
	}
	return nil
}

// Increment the amount of sellable tokens
func (pro *Producer) IncrementSellable(tokenIncrease int) error {
	if err := pro.EnforceCtx(); err != nil {
		return err
	}
	if tokenIncrease < 0 {
		return fmt.Errorf("err: cannot increase by negative amount")
	}
	err := pro.Ctx.UpdateHighThrough(fmt.Sprintf(SELLABLE, pro.ID), "+",
		tokenIncrease)
	if err != nil {
		return err
	}
	return nil
}

// Deduct the amount of tokens offered
func (pro *Producer) DeductTokens(tokenDeduction int) error {
	if tokenDeduction < 0 {
		return fmt.Errorf("err: cannot deduct negative tokens")
	}
	tokens, err := pro.GetTokens()
	if err != nil {
		return nil
	}
	if tokens-tokenDeduction < 0 {
		return fmt.Errorf("err producer does not have enough tokens")
	}
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	if err := pro.Ctx.UpdateHighThrough(highthrough, "-", tokenDeduction); err != nil {
		return err
	}
	return nil
}

// Increment the amount of tokens
func (pro *Producer) IncrementTokens(tokenIncrease int) error {
	if tokenIncrease < 0 {
		return fmt.Errorf("err: cannot increase tokens by negative amount")
	}
	highthrough := fmt.Sprintf(TOKEN, pro.ID)
	if err := pro.Ctx.UpdateHighThrough(highthrough, "+", tokenIncrease); err != nil {
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
