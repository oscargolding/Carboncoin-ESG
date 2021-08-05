package chaincode_test

import (
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/stretchr/testify/require"
)

func PerformChipStubs() (*chaincodefakes.CustomContex,
	*chaincodefakes.ChaincodeStub) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	context := &chaincodefakes.CustomContex{}
	context.GetStubReturns(chaincodeStub)
	return context, chaincodeStub
}

func Test_WHEN_invalid_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, stub := PerformChipStubs()
	stub.PutStateReturns(nil)
	ac := chaincode.AmountChip{Amount: 10, Valid: true, Owner: "oscar"}
	ac.InsertContext(ctx)

	// WHEN
	err := ac.MarkInvalid()

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, false, ac.Valid)
}

func Test_WHEN_noContext_THEN_ERROR(t *testing.T) {
	// GIVEN
	_, stub := PerformChipStubs()
	stub.PutStateReturns(nil)
	ac := chaincode.AmountChip{Amount: 10, Valid: true, Owner: "oscar"}

	// WHEN
	err := ac.MarkInvalid()

	// THEN
	require.EqualError(t, err, "err: the blockchain context is not set on chip")
	require.EqualValues(t, true, ac.Valid)
}
