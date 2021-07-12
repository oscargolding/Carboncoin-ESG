// Custom package for defining the offer object
package chaincode

import (
	"encoding/json"
	"fmt"
)

type Offer struct {
	DocType  string `json:"docType"`
	Producer string `json:"producer"`
	Amount   int    `json:"amount"`
	Tokens   int    `json:"tokens"`
	Active   bool   `json:"active"`
	OfferID  string `json:"offerId"`
}

// Remove tokens from the offer
func (off *Offer) RemoveTokens(tokenDecrease int) error {
	if !off.Active {
		return fmt.Errorf("offer not active")
	}
	if tokenDecrease < 0 {
		return fmt.Errorf("err: cannot decrease tokens by negative amount")
	}
	if off.Tokens-tokenDecrease < 0 {
		return fmt.Errorf("cannot purchase more tokens than offered")
	}
	off.Tokens -= tokenDecrease
	return nil
}

// Make a sale of the offer
func (off *Offer) MakeOfferStale() {
	off.Active = false
}

// Check if an offer is stale
func (off *Offer) IsStale() bool {
	return off.Tokens == 0
}

// Flush to the blockchain
func (off *Offer) ChainFlush(ctx CustomMarketContextInterface) error {
	jsonOffer, err := json.Marshal(off)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(off.OfferID, jsonOffer)
}
