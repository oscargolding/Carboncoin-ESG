// Custom package for defining the offer object
package chaincode

import (
	"encoding/json"
	"fmt"
)

type OfferModel struct {
	Producer   string `json:"producer"`
	Amount     int    `json:"amount"`
	Active     bool   `json:"active"`
	OfferID    string `json:"offerId"`
	Tokens     int    `json:"tokens"`
	Reputation int    `json:"reputation"`
	Owned      bool   `json:"owned"`
}

type Offer struct {
	DocType          string                       `json:"docType"`
	Producer         string                       `json:"producer"`
	Amount           int                          `json:"amount"`
	Active           bool                         `json:"active"`
	OfferID          string                       `json:"offerId"`
	CarbonReputation int                          `json:"carbonReputation"`
	Ctx              CustomMarketContextInterface `json:"-"`
}

const PROD_OFFER = "%s-offer"

func (off *Offer) EnforceCtx() error {
	if off.Ctx == nil {
		return fmt.Errorf("err: the blockchain context is not set on offer")
	}
	return nil
}

// Remove tokens from the offer
func (off *Offer) RemoveTokens(tokenDecrease int) error {
	if !off.Active {
		return fmt.Errorf("offer not active")
	}
	if err := off.EnforceCtx(); err != nil {
		return err
	}
	balance, err := off.Ctx.GetHighThrough(fmt.Sprintf(PROD_OFFER, off.OfferID))
	if err != nil {
		return err
	}
	if balance < 0 {
		return fmt.Errorf("err: cannot decrease tokens by negative amount")
	}
	if balance-tokenDecrease < 0 {
		return fmt.Errorf("cannot purchase more tokens than offered")
	}
	err = off.Ctx.UpdateHighThrough(PROD_OFFER, "-", tokenDecrease)
	if err != nil {
		return err
	}
	return nil
}

// Get a cleaner representation from the blockchain of an offer
func (off *Offer) ReturnModel() (*OfferModel, error) {
	if err := off.EnforceCtx(); err != nil {
		return nil, err
	}
	userId, err := off.Ctx.GetUserId()
	if err != nil {
		return nil, err
	}
	owned := false
	if userId == off.Producer {
		owned = true
	}
	balance, err := off.Ctx.GetHighThrough(fmt.Sprintf(PROD_OFFER, off.OfferID))
	if err != nil {
		return nil, err
	}
	return &OfferModel{Producer: off.Producer, Amount: off.Amount,
		Active: off.Active, OfferID: off.OfferID, Tokens: balance,
		Reputation: off.CarbonReputation, Owned: owned}, nil
}

// Get the balance on the offer
func (off *Offer) GetTokens() (int, error) {
	balance, err := off.Ctx.GetHighThrough(fmt.Sprintf(PROD_OFFER, off.OfferID))
	if err != nil {
		return 0, err
	}
	return balance, nil
}

// Set the tokens available on offer
func (off *Offer) SetTokens(tokenAmount int) error {
	// Have to add some amount of tokens
	if err := off.EnforceCtx(); err != nil {
		return err
	}
	err := off.Ctx.UpdateHighThrough(fmt.Sprintf(PROD_OFFER, off.OfferID),
		"+", tokenAmount)
	if err != nil {
		return err
	}
	return nil
}

// Make a sale of the offer
func (off *Offer) MakeOfferStale() {
	off.Active = false
}

// Check if an offer is stale
func (off *Offer) IsStale() (bool, error) {
	if err := off.EnforceCtx(); err != nil {
		return false, err
	}
	balance, err := off.Ctx.GetHighThrough(fmt.Sprintf(PROD_OFFER, off.OfferID))
	if err != nil {
		return false, err
	}
	return balance == 0, nil
}

func (off *Offer) InsertContext(ctx CustomMarketContextInterface) {
	off.Ctx = ctx
}

// Flush to the blockchain
func (off *Offer) ChainFlush(ctx CustomMarketContextInterface) error {
	jsonOffer, err := json.Marshal(off)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(off.OfferID, jsonOffer)
}
