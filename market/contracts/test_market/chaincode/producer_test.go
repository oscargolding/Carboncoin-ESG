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

	// WHEN
	require.EqualValues(t, 300, producer.Tokens)
	require.EqualValues(t, 300, producer.Sellable)
	require.EqualValues(t, "oscar", producer.ID)
}

func Test_WHEN_createMediumProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN, WHEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.MEDIUM, ctx)
	require.Nil(t, err)

	// THEN
	require.EqualValues(t, 200, producer.Tokens)
	require.EqualValues(t, 200, producer.Sellable)
}

func Test_WHEN_deductSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	require.Nil(t, err)

	// WHEN
	producer.DeductSellable(50)

	// THEN
	require.EqualValues(t, 250, producer.Sellable)
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

	// WHEN
	err = producer.DeductTokens(50, transactionContext)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Tokens, 250)
}

func Test_WHEN_deductTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, err := chaincode.NewProducer("oscar", chaincode.SMALL, ctx)
	require.Nil(t, err)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)

	// WHEN
	err = producer.DeductTokens(400, transactionContext)

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
	err := producer.IncrementTokens(50, transactionContext)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Tokens, 250)
}

func Test_WHEN_increaseSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)

	// WHEN
	err := producer.IncrementSellable(50)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Sellable, 350)
}

func Test_WHEN_addCarbon_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)

	// WHEN
	err := producer.AddCarbon(50, transactionContext)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Carbon, 50)
}

func Test_WHEN_addCarbonExtreme_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.UpdateHighThroughReturns(nil)

	// WHEN
	err := producer.AddCarbon(-50, transactionContext)

	// THEN
	require.EqualError(t, err, "err negative amount of carbon")
}

func Test_WHEN_getTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetHighThroughReturns(5, nil)

	// WHEN
	tokens, err := producer.GetTokens(transactionContext)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 5, tokens)
}

func Test_WHEN_getTokensErr_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx := GetStubProd()
	producer, _ := chaincode.NewProducer("oscar", chaincode.LARGE, ctx)
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetHighThroughReturns(0, fmt.Errorf("err"))

	// WHEN
	tokens, err := producer.GetTokens(transactionContext)

	// THEN
	require.EqualError(t, err, "err")
	require.EqualValues(t, 0, tokens)
}
