package chaincode_test

import (
	"encoding/json"
	"fmt"
	"register/chaincode"
	"register/chaincode/registerfakes"
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

func GetTestStubs() (*registerfakes.ChaincodeStub,
	*registerfakes.TransactionContext, chaincode.SmartContract) {
	chaincodeStub := &registerfakes.ChaincodeStub{}
	transactionContext := &registerfakes.TransactionContext{}
	transactionContext.GetStubReturns(chaincodeStub)

	// WHEN
	certifier := chaincode.SmartContract{}
	return chaincodeStub, transactionContext, certifier
}

func Test_WHEN_goodCertificate_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodestub, transactionContext, register := GetTestStubs()
	chaincodestub.PutStateReturns(nil)

	// WHEN
	err := register.CreateCertificate(transactionContext, "women", "social",
		0, 100, 2)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_badCertificate_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodestub, transactionContext, register := GetTestStubs()
	chaincodestub.PutStateReturns(fmt.Errorf("error"))

	// WHEN
	err := register.CreateCertificate(transactionContext, "women", "social",
		0, 100, 2)

	// THEN
	require.EqualError(t, err, "err: could not put the firm on the chain error")

}

func Test_WHEN_reportCertificate_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	chaincodestub, transactionContext, register := GetTestStubs()
	chaincodestub.InvokeChaincodeReturns(peer.Response{Status: 200})
	certificate := chaincode.ESGCertificate{Name: "women", Category: "social",
		Min: 1, Max: 5, Multiplier: 2}
	chaincodestub.GetStateReturns(json.Marshal(certificate))

	// WHEN
	err := register.RegisterUserCertificate(transactionContext,
		"oscar", 3, "oscar", "12", "a", "12%")

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_reportingOutBounds_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodestub, transactionContext, register := GetTestStubs()
	chaincodestub.InvokeChaincodeReturns(peer.Response{Status: 200})
	certificate := chaincode.ESGCertificate{Name: "women", Category: "social",
		Min: 1, Max: 5, Multiplier: 2}
	chaincodestub.GetStateReturns(json.Marshal(certificate))

	// WHEN
	err := register.RegisterUserCertificate(transactionContext,
		"oscar", 10, "oscar", "12", "a", "12%")

	// THEN
	require.EqualError(t, err, "err: value does not have the right format")
}

func Test_WHEN_reportingBadNetwork_THEN_FAILURE(t *testing.T) {
	// GIVEN
	chaincodestub, transactionContext, register := GetTestStubs()
	chaincodestub.InvokeChaincodeReturns(peer.Response{Status: 400,
		Message: "failed"})
	certificate := chaincode.ESGCertificate{Name: "women", Category: "social",
		Min: 1, Max: 5, Multiplier: 2}
	chaincodestub.GetStateReturns(json.Marshal(certificate))

	// WHEN
	err := register.RegisterUserCertificate(transactionContext,
		"oscar", 3, "oscar", "12", "a", "12%")

	// THEN
	require.EqualError(t, err, "err: failed calling chaincode: failed")
}
