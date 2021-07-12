package chaincode_test

import (
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WHEN_removeTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Tokens: 400, Active: true, OfferID: "1"}

	// WHEN
	err := offer.RemoveTokens(40)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, offer.Tokens, 360)
}

func Test_WHEN_removeTokensStale_THEN_FAILURE(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Active: false}

	// WHEN
	err := offer.RemoveTokens(40)

	// THEN
	require.EqualError(t, err, "offer not active")
}

func Test_WHEN_removeTokensLarge_THEN_FAILURE(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Active: true, Tokens: 400}

	// WHEN
	err := offer.RemoveTokens(500)

	// THEN
	require.EqualError(t, err, "cannot purchase more tokens than offered")
}

func Test_WHEN_makeOfferStale_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Active: true}

	// WHEN
	offer.MakeOfferStale()

	// THEN
	require.EqualValues(t, offer.Active, false)
}

func Test_WHEN_checkState_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Tokens: 0}

	// WHEN
	res := offer.IsStale()

	// THEN
	require.EqualValues(t, res, true)
}

func Test_WHEN_chainFlushOffer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(nil)
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Tokens: 400, Active: true, OfferID: "1"}

	// WHEN
	err := offer.ChainFlush(transactionContext)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_chainFlushOfferChainDown_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(fmt.Errorf("err"))
	offer := chaincode.Offer{Active: true}

	// WHEN
	err := offer.ChainFlush(transactionContext)

	// THEN
	require.EqualError(t, err, "err")
}
