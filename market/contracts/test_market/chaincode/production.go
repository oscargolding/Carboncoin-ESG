package chaincode

type Production struct {
	DocType      string `json:"docType"`
	ProductionID string `json:"productionID"`
	Produced     int    `json:"produced"`
	Date         string `json:"date"`
	Firm         string `json:"producingFirm"`
	Paid         bool   `json:"paid"`
}
