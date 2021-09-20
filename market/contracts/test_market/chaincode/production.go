package chaincode

import "encoding/json"

type Production struct {
	DocType      string `json:"docType"`
	ProductionID string `json:"productionID"`
	Produced     int    `json:"produced"`
	Date         string `json:"date"`
	Firm         string `json:"producingFirm"`
	Paid         bool   `json:"paid"`
	Ethical      bool   `json:"ethical"`
	Category     string `json:"category"`
	Description  string `json:"description"`
}

// Flush the production to the blockchain
//
// WARNING do not use after calling
func (prod *Production) ChainFlush(ctx CustomMarketContextInterface) error {
	jsonProduction, err := json.Marshal(prod)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(prod.ProductionID, jsonProduction)
}
