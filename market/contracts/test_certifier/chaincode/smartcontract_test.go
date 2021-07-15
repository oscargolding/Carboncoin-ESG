package chaincode_test

import (
	"certifier/chaincode"
	"certifier/chaincode/chaincodefakes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/require"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

func GetTestStubs() (*chaincodefakes.ChaincodeStub, *chaincodefakes.TransactionContext, chaincode.SmartContract) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	// WHEN
	certifier := chaincode.SmartContract{}
	return chaincodeStub, transactionContext, certifier
}

func TestInitLedger(t *testing.T) {
	// GIVEN
	_, transactionContext, certifier := GetTestStubs()

	// WHEN
	err := certifier.InitLedger(transactionContext)

	// THEN
	require.NoError(t, err)
}

func TestInitFailed(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()

	// WHEN
	chaincodeStub.PutStateReturns(fmt.Errorf("failed inserting key"))

	// THEN
	err := certifier.InitLedger(transactionContext)
	require.EqualError(t, err, "failed to put to world state. failed inserting key")
}

func TestSizeAsset(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	expectedFirm := &chaincode.Firm{ID: "firm", Size: "large"}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	chaincodeStub.GetStateReturns(bytes, nil)

	// WHEN
	size, err := certifier.FirmSize(transactionContext, "firm")

	// THEN
	require.NoError(t, err)
	require.Equal(t, size, "large")
}

func TestFirmMissing(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	chaincodeStub.GetStateReturns(nil, nil)

	// WHEN
	size, err := certifier.FirmSize(transactionContext, "oscar")

	// THEN
	require.NoError(t, err)
	require.Equal(t, "small", size)
}

func TestSizeError(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	chaincodeStub.GetStateReturns(nil, fmt.Errorf("failed inserting key"))

	// WHEN
	size, err := certifier.FirmSize(transactionContext, "oscar")

	// THEN
	require.EqualError(t, err, "failed to read the world state. failed inserting key")
	require.Equal(t, "", size)
}

func TestMarshalError(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	numbers := make([]byte, 0)
	chaincodeStub.GetStateReturns(numbers, nil)

	// WHEN
	size, err := certifier.FirmSize(transactionContext, "oscar")

	// THEN
	require.EqualError(t, err, "unexpected end of JSON input")
	require.Equal(t, "", size)
}

func Test_WHEN_reportProduction_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	chaincodeStub.InvokeChaincodeReturns(peer.Response{Status: 200})

	// WHEN
	err := certifier.ReportProduction(transactionContext, "oscar", 5, "1", "2")

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_reportProductionError_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodeStub, transactionContext, certifier := GetTestStubs()
	chaincodeStub.InvokeChaincodeReturns(peer.Response{Status: 400,
		Message: "failed"})

	// WHEN
	err := certifier.ReportProduction(transactionContext, "oscar", 5, "1", "2")

	// THEN
	require.EqualError(t, err, "err: failed calling chaincode: failed")
}
