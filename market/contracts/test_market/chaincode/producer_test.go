package chaincode_test

import (
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/stretchr/testify/require"
)

func GetStubProd() (*chaincodefakes.ChaincodeStub,
	*chaincodefakes.CustomContex) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(nil)
	transactionContext.UpdateHighThroughCalls(nil)
	return chaincodeStub, transactionContext
}

func Test_WHEN_createProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	require.Nil(t, err)
	_, _, tokens := ctx.UpdateHighThroughArgsForCall(0)
	_, _, sellable := ctx.UpdateHighThroughArgsForCall(1)

	// THEN
	require.EqualValues(t, 300, tokens)
	require.EqualValues(t, 300, sellable)
	require.EqualValues(t, "oscar", producer.ID)
}

func Test_WHEN_createMediumProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN, WHEN
	_, ctx := GetStubProd()
	_, err := chaincode.NewProducer("oscar", chaincode.MEDIUM, ctx)
	require.Nil(t, err)
	_, _, tokens := ctx.UpdateHighThroughArgsForCall(0)
	_, _, sellable := ctx.UpdateHighThroughArgsForCall(1)

	// THEN
	require.EqualValues(t, 200, tokens)
	require.EqualValues(t, 200, sellable)
}

func Test_WHEN_deductSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	ctx.GetHighThroughReturns(500, nil)
	require.Nil(t, err)

	// WHEN
	err = producer.DeductSellable(50)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_deductNegativeSellable_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	require.Nil(t, err)

	// WHEN
	err = producer.DeductSellable(-50)

	// THEN
	require.EqualError(t, err, "err negative sellable")
}

func Test_WHEN_deductingTooMuch_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	require.Nil(t, err)

	// WHEN
	err = producer.DeductSellable(500)

	// THEN
	require.EqualError(t, err, "err producer does not have enough sellable tokens")
}

func Test_WHEN_chainFlush_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(nil)
	transactionContext.UpdateHighThroughCalls(nil)
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE,
		transactionContext)
	require.Nil(t, err)

	// WHEN
	err = producer.ChainFlush(transactionContext)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_chainFlushError_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(fmt.Errorf("err"))
	transactionContext.UpdateHighThroughCalls(nil)
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE,
		transactionContext)
	require.Nil(t, err)

	// WHEN
	err = producer.ChainFlush(transactionContext)

	// THEN
	require.EqualError(t, err, "err")
}

func Test_WHEN_deductTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	require.Nil(t, err)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)
	ctx.GetHighThroughReturns(500, nil)

	// WHEN
	err = producer.DeductTokens(50)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_deductTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.SMALL, ctx)
	require.Nil(t, err)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)

	// WHEN
	err = producer.DeductTokens(400)

	// THEN
	require.EqualError(t, err, "err producer does not have enough tokens")
}

func Test_WHEN_increaseTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.MEDIUM, ctx)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)

	// WHEN
	err := producer.IncrementTokens(50)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_increaseSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)

	// WHEN
	err := producer.IncrementSellable(50)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_addCarbon_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	ctx.UpdateHighThroughReturns(nil)

	// WHEN
	err := producer.AddCarbon(50, "Environmental")

	// THEN
	require.Equal(t, producer.Environment, 50)
	require.Nil(t, err)
}

func Test_WHEN_addCarbonExtreme_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	ctx.UpdateHighThroughReturns(nil)

	// WHEN
	err := producer.AddCarbon(-50, "soc")
	_, sign, amount := ctx.UpdateHighThroughArgsForCall(3)

	// THEN
	require.Nil(t, err)
	require.Equal(t, 50, amount)
	require.Equal(t, "-", sign)
}

func Test_WHEN_addCarbonFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.MEDIUM, ctx)
	ctx.UpdateHighThroughReturns(fmt.Errorf("error"))

	// WHEN
	err := producer.AddCarbon(50, "Governance")

	// THEN
	require.EqualError(t, err, "error")
}

func Test_WHEN_getTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	ctx.GetHighThroughReturns(5, nil)

	// WHEN
	tokens, err := producer.GetTokens()

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 5, tokens)
}

func Test_WHEN_getTokensErr_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	ctx.GetHighThroughReturns(0, fmt.Errorf("err"))

	// WHEN
	tokens, err := producer.GetTokens()

	// THEN
	require.EqualError(t, err, "err")
	require.EqualValues(t, 0, tokens)
}

func Test_WHEN_getCarbon_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	producer.InsertContext(ctx)
	ctx.GetHighThroughReturns(0, nil)

	// WHEN
	tokens, err := producer.GetCarbon()

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 0, tokens)
}

func Test_WHEN_getCarbonFailed_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	producer.InsertContext(ctx)
	ctx.GetHighThroughReturns(0, fmt.Errorf("err"))

	// WHEN
	tokens, err := producer.GetCarbon()

	// THEN
	require.EqualError(t, err, "err")
	require.EqualValues(t, 0, tokens)
}
