package chaincode_test

import (
	"encoding/json"
	"fmt"
	"market/chaincode"
	"market/chaincode/chaincodefakes"
	"reflect"
	"testing"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/stretchr/testify/require"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name ClientIdentity . clientIdentity
type clientIdentity interface {
	cid.ClientIdentity
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name QueryIterator . queryIterator
type queryIterator interface {
	shim.StateQueryIteratorInterface
}

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -fake-name CustomContex . customContext
type customContext interface {
	chaincode.CustomMarketContextInterface
}

// Get all the stubs used and return
func GetTestStubs() (*chaincodefakes.ChaincodeStub,
	*chaincodefakes.CustomContex,
	chaincode.SmartContract) {
	chaincodeStub := &chaincodefakes.ChaincodeStub{}
	transactionContext := &chaincodefakes.CustomContex{}
	transactionContext.GetStubReturns(chaincodeStub)

	certifier := chaincode.SmartContract{}
	return chaincodeStub, transactionContext, certifier
}

func Test_WHEN_adminAdd_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("admin", nil)
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
	_, ctx, contract := GetTestStubs()
	ctx.CheckProducerReturns(true, nil)

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "producer with name oscar exists")
}

func Test_WHEN_nonAdminAdd_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.CheckProducerReturns(false, nil)
	ctx.GetUserTypeReturns("producer", nil)

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "cannot access unless admin")
}

func Test_WHEN_adminAddsFailedChain_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("admin", nil)
	stub.GetStateReturns(nil, nil)
	stub.InvokeChaincodeReturns(peer.Response{Payload: []byte("small")})
	stub.PutStateReturns(fmt.Errorf("failed"))

	// WHEN
	err := contract.AddProducer(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "failed")
}

func Test_WHEN_getBalance_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{Tokens: 500}
	ctx.GetProducerReturns(expectedFirm)

	// WHEN
	tokens, err := contract.GetBalance(ctx, "oscar")

	// THEN
	require.Equal(t, 500, tokens)
	require.Nil(t, err)
}

func Test_WHEN_getBalanceNotValid_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetProducerReturns(nil)

	// WHEN
	tokens, err := contract.GetBalance(ctx, "oscar")

	// THEN
	require.Equal(t, 0, tokens)
	require.EqualError(t, err, "unable to get producer with name: oscar")
}

func Test_WHEN_addOfferValid_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 500}
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.CreateOfferReturns(nil)
	stub.PutStateReturns(nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 300, "1")

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_addOfferNoProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetProducerReturns(nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 200, "1")

	require.EqualError(t, err, "failed to determine the existence of producer")
}

func Test_WHEN_addOfferNotProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 500}
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("", nil)
	ctx.GetSellableReturns(500, nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 200, "1")

	// THEN
	require.EqualError(t, err, "carboncoin offers only allowed by valid producers")
}

func Test_WHEN_addOfferNotEnoughTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 100}
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 300, "1")

	// THEN
	require.EqualError(t, err,
		"error deducting: err producer does not have enough sellable tokens")
}

func mock_function_offer(call interface{}, name string) (*reflect.Value,
	*reflect.Value) {
	funcType := reflect.TypeOf(call)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(call)
	expectedOffer := &chaincode.Offer{DocType: "offer", Producer: name,
		Amount: 30, Tokens: 30, Active: true, OfferID: "1"}
	bytes, _ := json.Marshal(expectedOffer)
	doc := reflect.New(docType)
	docInterface := doc.Interface()
	json.Unmarshal(bytes, &docInterface)
	return &callBackFunc, &doc
}

func Test_WHEN_getOffers_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer(call, "oscar")
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	queryResultIterator := &chaincodefakes.QueryIterator{}
	responseData := &peer.QueryResponseMetadata{FetchedRecordsCount: 1,
		Bookmark: ""}
	stub.GetQueryResultWithPaginationReturns(queryResultIterator, responseData, nil)

	// WHEN
	result, err := contract.GetOffers(ctx, 5, "")

	// THEN
	require.NoError(t, err)
	require.Equal(t, "", result.Bookmark)
	require.Equal(t, int32(1), result.FetchedRecordsCount)
	offers := result.Records[0]
	require.Equal(t, "oscar", offers.Producer)
	require.Equal(t, 30, offers.Tokens)
	require.Equal(t, 1, len(result.Records))
}

func Test_WHEN_getOffersNone_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		_ interface{}) error {
		return nil
	}
	queryResultIterator := &chaincodefakes.QueryIterator{}
	responseData := &peer.QueryResponseMetadata{FetchedRecordsCount: 0,
		Bookmark: ""}
	stub.GetQueryResultWithPaginationReturns(queryResultIterator, responseData, nil)

	// WHEN
	result, err := contract.GetOffers(ctx, 5, "")

	// THEN
	require.NoError(t, err)
	require.Equal(t, int32(0), result.FetchedRecordsCount)
	require.Equal(t, 0, len(result.Records))
}

func Test_WHEN_getOffersFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	stub.GetQueryResultWithPaginationReturns(nil, nil, fmt.Errorf("failure"))

	// WHEN
	result, err := contract.GetOffers(ctx, 5, "")

	// THEN
	require.EqualError(t, err, "error with query: failure")
	require.Nil(t, result)
}

func IdealPurchaseStub(stub *chaincodefakes.ChaincodeStub,
	ctx *chaincodefakes.CustomContex) {
	ctx.GetUserIdReturns("oscar", nil)
	ctx.GetUserTypeReturns("producer", nil)
	buyingFirm := &chaincode.Producer{ID: "oscar", Tokens: 500, Sellable: 100}
	sellingFirm := &chaincode.Producer{ID: "john", Tokens: 300, Sellable: 300}
	ctx.GetProducerReturnsOnCall(0, buyingFirm)
	ctx.GetProducerReturnsOnCall(1, sellingFirm)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_offer(i, "john")
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	stub.PutStateReturns(nil)
}
func Test_WHEN_purchaseTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)
	// WHEN
	amount, err := contract.PurchaseOfferTokens(ctx, "1", 5)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 505, amount)
}

func Test_WHEN_purchaseTokensTooMuch_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 1000)

	// THEN
	require.EqualError(t, err, "cannot purchase more tokens than offered")
}

func Test_WHEN_purchaseTokensNotProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("oscar", nil)
	ctx.GetUserTypeReturns("certifier", nil)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 10)

	// THEN
	require.EqualError(t, err, "must be a producer to purchase")
}

func Test_WHEN_producerNotEnoughTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)
	sellingFirm := &chaincode.Producer{ID: "john", Tokens: 5, Sellable: 300}
	ctx.GetProducerReturnsOnCall(1, sellingFirm)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 10)

	// THEN
	require.EqualError(t, err, "err producer does not have enough tokens")
}

func Test_WHEN_sellerIsSelf_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)
	sellingFirm := &chaincode.Producer{ID: "oscar", Tokens: 5, Sellable: 300}
	ctx.GetProducerReturnsOnCall(1, sellingFirm)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 10)

	// THEN
	require.EqualError(t, err, "err: cannot purchase tokens from self")
}
