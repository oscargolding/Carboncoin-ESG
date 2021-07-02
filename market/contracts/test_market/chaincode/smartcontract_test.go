package chaincode_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
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

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name ClientIdentity . clientIdentity
type clientIdentity interface {
	cid.ClientIdentity
}

// Get all the stubs used and return
func GetTestStubs() (*chaincodefakes.ChaincodeStub, *chaincodefakes.TransactionContext, *chaincodefakes.ClientIdentity, chaincode.SmartContract) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.TransactionContext{}
	clientIdentity := &chaincodefakes.ClientIdentity{}
	transactionContext.GetStubReturns(chaincodeStub)
	transactionContext.GetClientIdentityReturns(clientIdentity)

	certifier := chaincode.SmartContract{}
	return chaincodeStub, transactionContext, clientIdentity, certifier
}

func Test_WHEN_adminAdd_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	str := "x509::CN=admin,OU=client::CN=ca.org1.example.com," +
		"O=org1.example.com,L=Durham,ST=North Carolina,C=US"
	stub, ctx, id, contract := GetTestStubs()
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	id.GetIDReturns(sEnc, nil)
	stub.GetStateReturns(nil, nil)
	stub.InvokeChaincodeReturns(peer.Response{Payload: []byte("small")})
	stub.PutStateReturns(nil)

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.NoError(t, err)
}

func Test_WHEN_adminAddPresentUser_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, _, contract := GetTestStubs()
	stub.GetStateReturns([]byte("producer"), nil)

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "producer with name oscar exists")
}

func Test_WHEN_nonAdminAdd_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, id, contract := GetTestStubs()
	id.GetIDReturns("/CN=producer::/C=", nil)
	id.GetAttributeValueReturns("producer", true, nil)
	stub.GetStateReturns(nil, nil)
	stub.InvokeChaincodeReturns(peer.Response{Payload: []byte("small")})
	stub.PutStateReturns(nil)

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "cannot access unless admin")
}

func Test_WHEN_adminAddsFailedChain_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, id, contract := GetTestStubs()
	id.GetIDReturns("/CN=admin::/C=", nil)
	stub.GetStateReturns(nil, nil)
	stub.InvokeChaincodeReturns(peer.Response{Payload: []byte("small")})
	stub.PutStateReturns(fmt.Errorf("failed"))

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "error putting to world state. failed")
}

func Test_WHEN_getBalance_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, _, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{Tokens: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	tokens, err := contract.GetBalance(ctx, "oscar")

	// THEN
	require.Equal(t, 500, tokens)
	require.Nil(t, err)
}

func Test_WHEN_getBalanceNotValid_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, _, contract := GetTestStubs()
	stub.GetStateReturns(nil, nil)

	// WHEN
	tokens, err := contract.GetBalance(ctx, "oscar")

	// THEN
	require.Equal(t, 0, tokens)
	require.EqualError(t, err, "unable to get producer with name: oscar")
}

func Test_WHEN_addOfferValid_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, id, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)
	id.GetIDReturns("/CN=oscar::/C=", nil)
	id.GetAttributeValueReturns("producer", true, nil)
	stub.PutStateReturns(nil)

	// WHEN
	err = contract.AddOffer(ctx, "oscar", 200, 300)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_addOfferNoProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, _, contract := GetTestStubs()
	stub.GetStateReturns(nil, nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 200)

	require.EqualError(t, err, "failed to determine the existence of producer")
}

func Test_WHEN_notProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, id, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)
	id.GetIDReturns("/CN=oscar::/C=", nil)
	id.GetAttributeValueReturns("user", true, nil)
	stub.PutStateReturns(nil)

	// WHEN
	err = contract.AddOffer(ctx, "oscar", 200, 200)

	// THEN
	require.EqualError(t, err, "carboncoin offers only allowed by valid producers")
}

func Test_WHEN_addOfferNotEnoughTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, id, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 100}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)
	id.GetIDReturns("/CN=oscar::/C=", nil)
	id.GetAttributeValueReturns("producer", true, nil)
	stub.PutStateReturns(nil)

	// WHEN
	err = contract.AddOffer(ctx, "oscar", 200, 300)

	// THEN
	require.EqualError(t, err, "oscar does not have enough sellable tokens")
}
