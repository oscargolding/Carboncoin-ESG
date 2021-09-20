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
	expectedFirm := &chaincode.Producer{ID: "oscar"}
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
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	bytes, err := json.Marshal(expectedFirm)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	producer := ctx.GetProducer("oscar")

	// THEN
	require.NotNil(t, producer)
	require.NotNil(t, producer.Ctx)
	require.Equal(t, "oscar", producer.ID)
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

func Test_WHEN_createOffer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.PutStateReturns(nil)
	stub.GetStateReturns(nil, nil)
	stub.GetTxIDReturns("1")

	// WHEN
	err := ctx.CreateOffer("oscar", 5, 5, "1", 5)

	// THEN
	require.NoError(t, err)
}

func Test_WHEN_createOfferFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)
	stub.PutStateReturns(fmt.Errorf("error"))

	// WHEN
	err := ctx.CreateOffer("oscar", 5, 5, "oscar", 5)

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
	err = ctx.CreateOffer("oscar", 5, 5, "oscar", 0)

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
		Amount: 30, Active: true}
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

func Test_WHEN_GetHighThrough_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	queryResultIterator := chaincodefakes.QueryIterator{}
	queryResultIterator.HasNextReturnsOnCall(0, true)
	queryResultIterator.HasNextReturnsOnCall(1, true)
	queryResultIterator.HasNextReturnsOnCall(2, false)
	str := "oscarIndustry~+~2~1"
	queryResultIterator.NextReturns(&queryresult.KV{Key: str}, nil)
	stub.GetStateByPartialCompositeKeyReturns(&queryResultIterator, nil)
	stub.SplitCompositeKeyReturns("", []string{"", "+", "2"}, nil)

	// WHEN
	val, err := ctx.GetHighThrough("oscarIndustry")

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 2, val)
}

func Test_WHEN_GetHighThroughBadOper_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	queryResultIterator := chaincodefakes.QueryIterator{}
	queryResultIterator.HasNextReturnsOnCall(0, true)
	queryResultIterator.HasNextReturnsOnCall(1, true)
	queryResultIterator.HasNextReturnsOnCall(2, false)
	str := "oscarIndustry~+~2~1"
	queryResultIterator.NextReturns(&queryresult.KV{Key: str}, nil)
	stub.GetStateByPartialCompositeKeyReturns(&queryResultIterator, nil)
	stub.SplitCompositeKeyReturns("", []string{"", "/", "2"}, nil)

	// WHEN
	val, err := ctx.GetHighThrough("oscarIndustry")

	// THEN
	require.EqualValues(t, val, 0)
	require.EqualError(t, err, "unexpected operation /")
}

func Test_WHEN_GetHighThroughCompErr_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateByPartialCompositeKeyReturns(nil, fmt.Errorf("err"))

	// WHEN
	val, err := ctx.GetHighThrough("oscarIndustry")

	// THEN
	require.EqualError(t, err, "could not retrieve value for oscarIndustry: err")
	require.EqualValues(t, val, 0)
}

func Test_WHEN_getResult_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	offer := &chaincode.Offer{DocType: "offer", Producer: "oscar",
		Amount: 30, Active: true}
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

func Test_WHEN_createProduction_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)
	stub.PutStateReturns(nil)

	// WHEN
	err := ctx.CreateProduction("1", 2, "12/2", "oscar", true, false, "Energy",
		"greenhouse")

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_createProductionExists_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	producer := &chaincode.Production{DocType: "production"}
	bytes, err := json.Marshal(producer)
	require.NoError(t, err)
	stub.GetStateReturns(bytes, nil)

	// WHEN
	err = ctx.CreateProduction("1", 2, "12/2", "oscar", true, false, "Energy",
		"greenhouse")

	// THEN
	require.EqualError(t, err, "production with id already exists on the market")
}

func Test_WHEN_createProductionBlcokError_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetStateReturns(nil, nil)
	stub.PutStateReturns(fmt.Errorf("error"))

	// WHEN
	err := ctx.CreateProduction("1", 2, "12/2", "oscar", true, false, "Energy",
		"greenhouse")

	// THEN
	require.EqualError(t, err, "error")
}

func Test_WHEN_updateHighThroughput_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetTxIDReturns("1")
	stub.PutStateReturns(nil)

	// WHEN
	err := ctx.UpdateHighThrough("oscarproduction", "+", 5)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_updateHighThroughputMinus_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.GetTxIDReturns("1")

	// WHEN
	err := ctx.UpdateHighThrough("oscarproduction", "/", 6)

	// THEN
	require.EqualError(t, err, "operator / op is not supported")
}

func Test_WHEN_generateOfferString_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, _ := PerformTestStub()

	// WHEN
	str := ctx.OfferStringGenerator("reputation", true, "")

	// THEN
	expected := `{"selector":{"docType":"offer", "active": true},` +
		`"sort":[{"carbonReputation":"asc"}]}`
	require.EqualValues(t, expected, str)
}

func Test_WHEN_generateNone_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, _ := PerformTestStub()

	// WHEN
	str := ctx.OfferStringGenerator("", false, "")

	// THEN
	expected := `{"selector":{"docType":"offer", "active": true}}`
	require.EqualValues(t, expected, str)
}

func Test_WHEN_generateOfferStringDesc_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, _ := PerformTestStub()

	// WHEN
	str := ctx.OfferStringGenerator("reputation", false, "")

	// THEN
	expected := `{"selector":{"docType":"offer", "active": true},` +
		`"sort":[{"carbonReputation":"desc"}]}`
	require.EqualValues(t, expected, str)
}

func Test_WHEN_generateOfferStringDescUser_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, _ := PerformTestStub()

	// WHEN
	str := ctx.OfferStringGenerator("", false, "oscarIndustry")

	// THEN
	expected := `{"selector":{"docType":"offer", "active": true,` +
		` "producer": "oscarIndustry"}}`
	require.EqualValues(t, expected, str)
}

func Test_WHEN_createChip_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.PutStateReturns(nil)

	// WHEN
	err := ctx.CreateChip("oscar", 50)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_createChip_THEN_FAILURE(t *testing.T) {
	// GIVEN
	ctx, _, stub := PerformTestStub()
	stub.PutStateReturns(fmt.Errorf("error"))

	// WHEN
	err := ctx.CreateChip("oscar", 30)

	// THEN
	require.EqualError(t, err, "error")
}
