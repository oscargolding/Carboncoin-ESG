package chaincode

import (
	"encoding/json"
	"fmt"
)

type AmountChip struct {
	Amount int                          `json:"amount"`
	Valid  bool                         `json:"valid"`
	Owner  string                       `json:"owner"`
	Ctx    CustomMarketContextInterface `json:"-"`
}

func (ac *AmountChip) EnforceCtx() error {
	if ac.Ctx == nil {
		return fmt.Errorf("err: the blockchain context is not set on chip")
	}
	return nil
}

func (ac *AmountChip) InsertContext(ctx CustomMarketContextInterface) {
	ac.Ctx = ctx
}

// Marks a chip as invalid
func (ac *AmountChip) MarkInvalid() error {
	if err := ac.EnforceCtx(); err != nil {
		return err
	}
	ac.Valid = false
	return ac.ChainFlush()
}

func (ac *AmountChip) ChainFlush() error {
	achip, err := json.Marshal(ac)
	if err != nil {
		return err
	}
	return ac.Ctx.GetStub().PutState(fmt.Sprintf(CHIP, ac.Owner), achip)
}
