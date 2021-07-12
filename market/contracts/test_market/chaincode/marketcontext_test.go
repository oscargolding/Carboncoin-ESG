package chaincode_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"testing"

	"github.com/hyperledger/fabric-protos-go/ledger/queryresult"
	"github.com/stretchr/testify/require"
)

func PerformTestStub() (*chaincode.CustomMarketContext,
	*chaincodefakes.ClientIdentity, *chaincodefakes.ChaincodeStub) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	clientIdentity := &chaincodefakes.ClientIdentity{}
	first := &chaincode.CustomMarketContext{}
	first.SetStub(chaincodeStub)
	first.SetClientIdentity(clientIdentity)
	return first, clientIdentity, chaincodeStub
}

func Test_WHEN_getID_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	str := "x509::CN=admin,OU=client::CN=ca.org1.example.com," +
		"O=org1.example.com,L=Durham,ST=North Carolina,C=US"
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	ctx, id, _ := PerformTestStub()
	id.GetIDReturns(sEnc, nil)

	// WHEN
	idString, err := ctx.GetUserId()

	// THEN
	require.NoError(t, err)
	require.Equal(t, "admin", idString)
}

func Test_WHEN_invalidID_THEN_FAILURE(t *testing.T) {
	// GIVEN
	str := "x509::CNoscar,OU=client::CN=ca.org1.example.com,"
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	ctx, id, _ := PerformTestStub()
	id.GetIDReturns(sEnc, nil)

	// WHEN
	idString, err := ctx.GetUserId()

	// THEN
	require.EqualError(t, err, "invalid ID")
	require.Equal(t, "", idString)
}

func Test_WHEN_adminAttribute_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	str := "x509::CN=admin,OU=client::CN=ca.org1.example.com," +
		"O=org1.example.com,L=Durham,ST=North Carolina,C=US"
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	ctx, id, _ := PerformTestStub()
	id.GetIDReturns(sEnc, nil)

	// WHEN
	idString, err := ctx.GetUserType()

	// THEN
	require.NoError(t, err)
	require.Equal(t, "admin", idString)
}

func Test_WHEN_producerGetAttribute_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	str := "x509::CN=oscar,OU=client::CN=ca.org1.example.com," +
		"O=org1.example.com,L=Durham,ST=North Carolina,C=US"
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	ctx, id, _ := PerformTestStub()
	id.GetIDReturns(sEnc, nil)
	id.GetAttributeValueReturns("producer", true, nil)

	// WHEN
	idString, err := ctx.GetUserType()

	// THEN
	require.NoError(t, err)
	require.Equal(t, "producer", idString)
}

func Test_WHEN_attributeNotFound_THEN_FAILURE(t *testing.T) {
	// GIVEN
	str := "x509::CN=oscar,OU=client::CN=ca.org1.example.com," +
		"O=org1.example.com,L=Durham,ST=North Carolina,C=US"
	sEnc := base64.StdEncoding.EncodeToString([]byte(str))
	ctx, id, _ := PerformTestStub()
	id.GetIDReturns(sEnc, nil)
	id.GetAttributeValueReturns("", false, nil)

	// WHEN
	idString, err := ctx.GetUserType()

	// THEN
	require.EqualError(t, err, "no `usertype` attribute value found")
	require.EqualValues(t, idString, "")
}

func Test_WHEN_getProducerFound_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	expectedFirm := &chaincode.Producer{Tokens: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	found, err := ctx.CheckProducer("oscar")

	// THEN
	require.Nil(t, err)
	require.Equal(t, true, found)
}

func Test_WHEN_getProducerNotFound_THEN_FAILURE(t *testing.T) {
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)

	// WHEN
	found, err := ctx.CheckProducer("oscar")

	// THEN
	require.Nil(t, err)
	require.Equal(t, false, found)
}

func Test_WHEN_getProducer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	expectedFirm := &chaincode.Producer{Tokens: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	producer := ctx.GetProducer("oscar")

	// THEN
	require.NotNil(t, producer)
	require.Equal(t, 500, producer.Tokens)
}

func Test_WHEN_getProducerFailed_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)

	// WHEN
	producer := ctx.GetProducer("oscar")

	// THEN
	require.Nil(t, producer)
}

func Test_WHEN_getSellable_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	expectedFirm := &chaincode.Producer{Tokens: 500, Sellable: 500}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	sellable, err := ctx.GetSellable("oscar")

	// THEN
	require.NoError(t, err)
	require.Equal(t, 500, sellable)
}

func Test_WHEN_getSellableNoProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)

	// WHEN
	sellable, err := ctx.GetSellable("oscar")

	// THEN
	require.Equal(t, 0, sellable)
	require.EqualError(t, err, "unable to get sellable with name: oscar")
}

func Test_WHEN_createOffer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.PutStateReturns(nil)
	stub.GetStateReturns(nil, nil)

	// WHEN
	err := ctx.CreateOffer("oscar", 5, 5, "1")

	// THEN
	require.NoError(t, err)
}

func Test_WHEN_createOfferFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)
	stub.PutStateReturns(fmt.Errorf("error"))

	// WHEN
	err := ctx.CreateOffer("oscar", 5, 5, "oscar")

	// THEN
	require.EqualError(t, err, "error")
}

func Test_WHEN_createOfferDuplicateId_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	expectedFirm := &chaincode.Offer{OfferID: "1"}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	err = ctx.CreateOffer("oscar", 5, 5, "oscar")

	// THEN
	require.EqualError(t, err, "offer with id already exists on the market")
}

func Test_WHEN_queryIterator_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, _ := PerformTestStub()
	queryResultIterator := chaincodefakes.QueryIterator{}
	queryResultIterator.HasNextReturnsOnCall(0, true)
	queryResultIterator.HasNextReturnsOnCall(1, false)
	expectedOffer := &chaincode.Offer{DocType: "offer", Producer: "oscar",
		Amount: 30, Tokens: 30, Active: true}
	bytes, err := json.Marshal(expectedOffer)
	require.NoError(t, err)
	queryResultIterator.NextReturns(&queryresult.KV{Value: bytes}, nil)
	queryResultIterator.CloseReturns(fmt.Errorf("done"))

	// WHEN
	ctx.IteratorResults(&queryResultIterator, func(offer *chaincode.Offer) {
		// THEN
		require.Equal(t, offer.Amount, 30)
	})
}

func Test_WHEN_getResult_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	offer := &chaincode.Offer{DocType: "offer", Producer: "oscar",
		Amount: 30, Tokens: 30, Active: true}
	bytes, err := json.Marshal(offer)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	err = ctx.GetResult("oscar", func(offer *chaincode.Offer) {
		// THEN
		var usingOffer *chaincode.Offer = offer
		require.Equal(t, offer.Producer, "oscar")
		require.Equal(t, usingOffer.Producer, "oscar")
	})
	require.Nil(t, err)
}

func Test_WHEN_getResult_THEN_ERROR(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)

	// WHEN
	err := ctx.GetResult("oscar", func(offer *chaincode.Offer) {})
	require.EqualError(t, err, "failed getting state")
}
