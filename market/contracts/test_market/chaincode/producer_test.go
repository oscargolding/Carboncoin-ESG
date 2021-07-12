package chaincode_test

import (
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_WHEN_createProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	require.EqualValues(t, 300, producer.Tokens)
	require.EqualValues(t, 300, producer.Sellable)
	require.EqualValues(t, "oscar", producer.ID)
}

func Test_WHEN_createMediumProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN, WHEN
	producer := chaincode.NewProducer("oscar", chaincode.MEDIUM)

	// THEN
	require.EqualValues(t, 200, producer.Tokens)
	require.EqualValues(t, 200, producer.Sellable)
}

func Test_WHEN_deductSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	producer.DeductSellable(50)

	// THEN
	require.EqualValues(t, 250, producer.Sellable)
}

func Test_WHEN_deductNegativeSellable_THEN_FAILURE(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.DeductSellable(-50)

	// THEN
	require.EqualError(t, err, "err negative sellable")
}

func Test_WHEN_deductingTooMuch_THEN_FAILURE(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.DeductSellable(500)

	// THEN
	require.EqualError(t, err, "err producer does not have enough sellable tokens")
}

func Test_WHEN_chainFlush_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(nil)
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.ChainFlush(transactionContext)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_chainFlushError_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)
	chaincodeStub.PutStateReturns(fmt.Errorf("err"))
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.ChainFlush(transactionContext)

	// THEN
	require.EqualError(t, err, "err")
}

func Test_WHEN_deductTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.DeductTokens(50)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Tokens, 250)
}

func Test_WHEN_deductTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.SMALL)

	// WHEN
	err := producer.DeductTokens(400)

	// THEN
	require.EqualError(t, err, "err producer does not have enough tokens")
}

func Test_WHEN_increaseTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.MEDIUM)

	// WHEN
	err := producer.IncrementTokens(50)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Tokens, 250)
}

func Test_WHEN_increaseSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	producer := chaincode.NewProducer("oscar", chaincode.LARGE)

	// WHEN
	err := producer.IncrementSellable(50)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, producer.Sellable, 350)
}
