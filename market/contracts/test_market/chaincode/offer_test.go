package chaincode_test

import (
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/stretchr/testify/require"
)

func PerformTestStubs() (*chaincodefakes.CustomContex,
	*chaincodefakes.ChaincodeStub) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	context := &chaincodefakes.CustomContex{}
	context.GetStubReturns(chaincodeStub)
	return context, chaincodeStub
}

func Test_WHEN_removeTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _ := PerformTestStubs()
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Active: true, OfferID: "1"}
	offer.InsertContext(ctx)
	ctx.GetHighThroughReturns(500, nil)

	// WHEN
	err := offer.RemoveTokens(40)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_removeTokensStale_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _ := PerformTestStubs()
	ctx.GetHighThroughReturns(0, nil)
	offer := chaincode.Offer{Active: false}

	// WHEN
	err := offer.RemoveTokens(40)

	// THEN
	require.EqualError(t, err, "offer not active")
}

func Test_WHEN_removeTokensLarge_THEN_FAILURE(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Active: true}
	ctx, _ := PerformTestStubs()
	offer.InsertContext(ctx)
	ctx.GetHighThroughReturns(300, nil)

	// WHEN
	err := offer.RemoveTokens(500)

	// THEN
	require.EqualError(t, err, "cannot purchase more tokens than offered")
}

func Test_WHEN_getModel_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Active: true, OfferID: "1"}
	ctx, _ := PerformTestStubs()
	ctx.UpdateHighThroughReturns(nil)
	ctx.GetHighThroughReturns(50, nil)
	ctx.GetUserIdReturns("admin", nil)
	offer.InsertContext(ctx)

	// WHEN
	model, err := offer.ReturnModel()

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 50, model.Tokens)
	require.EqualValues(t, true, model.Active)
	require.EqualValues(t, "1", model.OfferID)
	require.EqualValues(t, "oscar", model.Producer)
	require.EqualValues(t, 500, model.Amount)
	require.EqualValues(t, false, model.Owned)
}

func Test_WHEN_getModelUser_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Active: true, OfferID: "1"}
	ctx, _ := PerformTestStubs()
	ctx.UpdateHighThroughReturns(nil)
	ctx.GetHighThroughReturns(50, nil)
	offer.InsertContext(ctx)
	ctx.GetUserIdReturns("oscar", nil)

	// WHEN
	model, err := offer.ReturnModel()
	require.Nil(t, err)
	require.EqualValues(t, true, model.Owned)
}

func Test_WHEN_blockchainCtxFailed_THEN_FAILURE(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{Active: true}

	// WHEN
	err := offer.RemoveTokens(50)

	// THEN
	require.EqualError(t, err, "err: the blockchain context is not set on offer")
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
	offer := chaincode.Offer{}
	ctx, _ := PerformTestStubs()
	ctx.GetHighThroughReturns(0, nil)
	offer.InsertContext(ctx)

	// WHEN
	res, err := offer.IsStale()
	require.Nil(t, err)

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
		Active: true, OfferID: "1"}

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

func Test_WHEN_setTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	offer := chaincode.Offer{DocType: "offer", Producer: "oscar", Amount: 500,
		Active: true, OfferID: "1"}
	ctx, _ := PerformTestStubs()
	ctx.UpdateHighThroughReturns(nil)
	offer.InsertContext(ctx)

	// WHEN
	err := offer.SetTokens(500)

	// THEN
	require.Nil(t, err)
}
