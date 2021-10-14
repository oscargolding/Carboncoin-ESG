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
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetHighThroughReturns(500, nil)

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
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.CreateOfferReturns(nil)
	ctx.GetHighThroughReturns(400, nil)
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
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("", nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 200, "1")

	// THEN
	require.EqualError(t, err, "carboncoin offers only allowed by valid producers")
}

func Test_WHEN_addOfferNotEnoughTokens_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 300, "1")

	// THEN
	require.EqualError(t, err,
		"error deducting: err producer does not have enough sellable tokens")
}

func Test_WHEN_addOfferCarbonFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.CreateOfferReturns(nil)
	ctx.GetHighThroughReturnsOnCall(0, 400, nil)
	ctx.GetHighThroughReturnsOnCall(1, 0, fmt.Errorf("err"))

	// WHEN
	err := contract.AddOffer(ctx, "oscar", 200, 300, "1")

	// THEN
	require.EqualError(t, err, "err: getting carbon err")
}

func setup_mock_offer(call interface{}, offer *chaincode.Offer) (*reflect.Value,
	*reflect.Value) {
	funcType := reflect.TypeOf(call)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(call)
	bytes, _ := json.Marshal(offer)
	doc := reflect.New(docType)
	docInterface := doc.Interface()
	json.Unmarshal(bytes, &docInterface)
	return &callBackFunc, &doc
}

func mock_function_offer(call interface{}, name string) (*reflect.Value,
	*reflect.Value) {
	expectedOffer := &chaincode.Offer{DocType: "offer", Producer: name,
		Amount: 30, Active: true, OfferID: "1", CarbonReputation: 10}
	return setup_mock_offer(call, expectedOffer)
}

func mock_function_offer_detailed(call interface{}, name string, rep int,
	amount int) (*reflect.Value, *reflect.Value) {
	expectedOffer := &chaincode.Offer{DocType: "offer", Producer: name,
		Amount: amount, Active: true, OfferID: "1", CarbonReputation: rep}
	return setup_mock_offer(call, expectedOffer)
}

func mock_function_production(call interface{}, id string, paid bool) (*reflect.Value,
	*reflect.Value) {
	funcType := reflect.TypeOf(call)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(call)
	expectedProd := &chaincode.Production{ProductionID: id, DocType: "production",
		Paid: paid, Produced: 4, Firm: id}
	bytes, _ := json.Marshal(expectedProd)
	doc := reflect.New(docType)
	docInterface := doc.Interface()
	json.Unmarshal(bytes, &docInterface)
	return &callBackFunc, &doc
}

func mock_amount_chip(call interface{}, id string, valid bool) (*reflect.Value,
	*reflect.Value) {
	funcType := reflect.TypeOf(call)
	docType := funcType.In(0).Elem()
	callBackFunc := reflect.ValueOf(call)
	expectedProd := &chaincode.AmountChip{Amount: 80, Owner: id, Valid: valid}
	bytes, _ := json.Marshal(expectedProd)
	doc := reflect.New(docType)
	docInterface := doc.Interface()
	json.Unmarshal(bytes, &docInterface)
	return &callBackFunc, &doc
}

func Test_WHEN_getProduction_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("oscar", nil)
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_production(call, "oscar", false)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetHighThroughReturns(500, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	responseData := &peer.QueryResponseMetadata{FetchedRecordsCount: 1,
		Bookmark: ""}
	stub.GetQueryResultWithPaginationReturns(queryResultIterator, responseData, nil)

	// WHEN
	result, err := contract.GetProduction(ctx, 5, "")

	// THEN
	require.Nil(t, err)
	require.Equal(t, "", result.Bookmark)
	require.Equal(t, int32(1), result.FetchedRecordsCount)
	require.Equal(t, 1, len(result.Records))
	require.Equal(t, 500, result.Reputation)
	production := result.Records[0]
	require.Equal(t, false, production.Paid)
}

func Test_WHEN_getProductionNone_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		_ interface{}) error {
		return nil
	}
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetHighThroughReturns(500, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	responseData := &peer.QueryResponseMetadata{FetchedRecordsCount: 0,
		Bookmark: ""}
	stub.GetQueryResultWithPaginationReturns(queryResultIterator, responseData, nil)

	// WHEN
	result, err := contract.GetProduction(ctx, 5, "")

	// THEN
	require.Nil(t, err)
	require.Equal(t, int32(0), result.FetchedRecordsCount)
	require.Equal(t, 0, len(result.Records))
}

func Test_WHEN_getProductionErrors_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	stub.GetQueryResultWithPaginationReturns(nil, nil, fmt.Errorf("failure"))
	expectedFirm := &chaincode.Producer{ID: "oscar"}
	expectedFirm.InsertContext(ctx)
	ctx.GetProducerReturns(expectedFirm)
	ctx.GetHighThroughReturns(500, nil)

	// WHEN
	result, err := contract.GetProduction(ctx, 5, "")

	// THEN
	require.EqualError(t, err, "error with query: failure")
	require.Nil(t, result)
}

func Test_WHEN_getDirectPrice_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer(call, "oscar")
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.CreateChipReturns(nil)
	queryResultsIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultsIterator, nil)

	// WHEN
	result, err := contract.GetDirectPrice(ctx)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 80, result)
}

func Test_WHEN_getDirectPriceEmpty_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		_ interface{}) error {
		return nil
	}
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.CreateChipReturns(nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	result, err := contract.GetDirectPrice(ctx)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 50, result)
}

func Test_WHEN_getDirectPriceFailure_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	stub.GetQueryResultReturns(nil, fmt.Errorf("error"))

	// WHEN
	result, err := contract.GetDirectPrice(ctx)

	// THEN
	require.EqualValues(t, 0, result)
	require.EqualError(t, err, "error with query: error")
}

func Test_WHEN_getDirectPriceNotProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		_ interface{}) error {
		return nil
	}
	ctx.GetUserTypeReturns("regulator", nil)
	ctx.GetUserIdReturns("asic", nil)
	ctx.CreateChipReturns(nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	result, err := contract.GetDirectPrice(ctx)

	// THEN
	require.EqualValues(t, 0, result)
	require.EqualError(t, err, "err: only producers can get a price")
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
	ctx.GetHighThroughReturns(300, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	responseData := &peer.QueryResponseMetadata{FetchedRecordsCount: 1,
		Bookmark: ""}
	stub.GetQueryResultWithPaginationReturns(queryResultIterator, responseData, nil)

	// WHEN
	result, err := contract.GetOffers(ctx, 5, "", "", false, "")

	// THEN
	require.NoError(t, err)
	require.Equal(t, "", result.Bookmark)
	require.Equal(t, int32(1), result.FetchedRecordsCount)
	offers := result.Records[0]
	require.Equal(t, "oscar", offers.Producer)
	require.Equal(t, 300, offers.Tokens)
	require.Equal(t, 1, len(result.Records))
	require.Equal(t, 10, offers.Reputation)
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
	result, err := contract.GetOffers(ctx, 5, "", "", false, "")

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
	result, err := contract.GetOffers(ctx, 5, "", "", false, "")

	// THEN
	require.EqualError(t, err, "error with query: failure")
	require.Nil(t, result)
}

func Test_WHEN_getBudgetOffer_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("john", nil)
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer(call, "oscar")
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	ctx.GetHighThroughReturns(300, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	results, err := contract.GetBudgetOffer(ctx, false, 250)

	// THEN
	require.NoError(t, err)
	require.Equal(t, 1, len(results.Records))
}

func Test_WHEN_getBudgetOfferMany_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("john", nil)
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer_detailed(call, "oscar", 10, 12)
		callBackFunc.Call([]reflect.Value{*doc})
		secondCall, secondDoc := mock_function_offer_detailed(call, "james", 10, 12)
		secondCall.Call([]reflect.Value{*secondDoc})
		return nil
	}
	ctx.GetHighThroughReturns(300, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	results, err := contract.GetBudgetOffer(ctx, false, 250)

	// THEN
	require.NoError(t, err)
	require.Equal(t, 1, len(results.Records))
}

func Test_WHEN_getBudgetOfferManyReputation_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("john", nil)
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer_detailed(call, "oscar", 12, 12)
		callBackFunc.Call([]reflect.Value{*doc})
		secondCall, secondDoc := mock_function_offer_detailed(call, "james", 8, 12)
		secondCall.Call([]reflect.Value{*secondDoc})
		return nil
	}
	ctx.GetHighThroughReturns(300, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	results, err := contract.GetBudgetOffer(ctx, true, 5)

	// THEN
	require.NoError(t, err)
	require.Equal(t, 1, len(results.Records))
	require.Equal(t, "oscar", results.Records[0].Producer)
}

func Test_WHEN_getBudgetOfferWithDollar_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("john", nil)
	ctx.IteratorResultsStub = func(_ shim.StateQueryIteratorInterface,
		call interface{}) error {
		callBackFunc, doc := mock_function_offer_detailed(call, "oscar", 12, 15)
		callBackFunc.Call([]reflect.Value{*doc})
		secondCall, secondDoc := mock_function_offer_detailed(call, "james", 8, 10)
		secondCall.Call([]reflect.Value{*secondDoc})
		return nil
	}
	ctx.GetHighThroughReturns(300, nil)
	queryResultIterator := &chaincodefakes.QueryIterator{}
	stub.GetQueryResultReturns(queryResultIterator, nil)

	// WHEN
	results, err := contract.GetBudgetOffer(ctx, false, 600)

	// THEN
	require.NoError(t, err)
	require.Equal(t, 2, len(results.Records))
	require.Equal(t, "james", results.Records[0].Producer)
	require.Equal(t, "oscar", results.Records[1].Producer)

}

func Test_WHEN_getBudgetError_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserIdReturns("oscar", nil)
	stub.GetQueryResultReturns(nil, fmt.Errorf("error"))

	// WHEN
	results, err := contract.GetBudgetOffer(ctx, false, 50)

	// THEN
	require.EqualError(t, err, "err: error calling blockchain error")
	require.Nil(t, results)
}

func IdealPurchaseStub(stub *chaincodefakes.ChaincodeStub,
	ctx *chaincodefakes.CustomContex) {
	ctx.GetUserIdReturns("oscar", nil)
	ctx.GetUserTypeReturns("producer", nil)
	buyingFirm := &chaincode.Producer{ID: "oscar"}
	sellingFirm := &chaincode.Producer{ID: "john"}
	buyingFirm.InsertContext(ctx)
	sellingFirm.InsertContext(ctx)
	ctx.GetProducerReturnsOnCall(0, buyingFirm)
	ctx.GetProducerReturnsOnCall(1, sellingFirm)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_offer(i, "john")
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	stub.PutStateReturns(nil)
}

func StaleStub(stub *chaincodefakes.ChaincodeStub,
	ctx *chaincodefakes.CustomContex, name string) {
	ctx.GetUserIdReturns(name, nil)
	ctx.GetUserTypeReturns("producer", nil)
	buyingFirm := &chaincode.Producer{ID: name}
	buyingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(buyingFirm)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_offer(i, name)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	stub.PutStateReturns(nil)
}

func Test_WHEN_makeOfferStaleUser_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	StaleStub(stub, ctx, "oscar")
	ctx.GetHighThroughReturns(500, nil)

	// WHEN
	err := contract.MakeOfferStale(ctx, "1")

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_makeOfferStaleNotMatch_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	StaleStub(stub, ctx, "oscar")
	ctx.GetHighThroughReturns(500, nil)
	ctx.GetUserIdReturns("james", nil)

	// WHEN
	err := contract.MakeOfferStale(ctx, "1")

	// THEN
	require.EqualError(t, err,
		"err: user does not own the offer and cannot cancel it")

}
func Test_WHEN_purchaseTokens_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)
	ctx.GetHighThroughReturns(500, nil)
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
	ctx.GetHighThroughReturns(500, nil)

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
	ctx.GetHighThroughReturnsOnCall(0, 10, nil)
	ctx.GetHighThroughReturnsOnCall(1, 5, nil)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 10)

	// THEN
	require.EqualError(t, err, "err producer does not have enough tokens")
}

func Test_WHEN_sellerIsSelf_THEN_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	IdealPurchaseStub(stub, ctx)
	sellingFirm := &chaincode.Producer{ID: "oscar"}
	ctx.GetProducerReturnsOnCall(1, sellingFirm)

	// WHEN
	_, err := contract.PurchaseOfferTokens(ctx, "1", 10)

	// THEN
	require.EqualError(t, err, "err: cannot purchase tokens from self")
}

func Test_WHEN_producerProduction_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("certifier", nil)
	sellingFirm := &chaincode.Producer{ID: "oscar"}
	sellingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(sellingFirm)
	ctx.CreateProductionReturns(nil)
	ctx.UpdateHighThroughReturns(nil)
	stub.PutStateReturns(nil)

	// WHEN
	err := contract.ProducerProduction(ctx, "oscar", 10, "1/1", "1",
		"energy", "greenhouse", "12 co2e", 1)

	// THEN
	require.Nil(t, err)
}

func Test_WHEN_producerProductionNotCertifier_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)

	// WHEN
	err := contract.ProducerProduction(ctx, "oscar", 10, "1/1", "1", "energy",
		"greenhouse", "12 co2e", 1)

	// THEN
	require.EqualError(t, err,
		"err: only certifiers/admins can report carbon production")
}

func Test_WHEN_producerProductionDNE_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetProducerReturns(nil)
	ctx.GetUserTypeReturns("certifier", nil)

	// WHEN
	err := contract.ProducerProduction(ctx, "oscar", 15, "1/1", "1", "energy",
		"greenhouse", "12 co2e", 1)

	// THEN
	require.EqualError(t, err, "err: producer does not exist")
}

func Test_WHEN_producerProductionErrorCreating_FAILURE(t *testing.T) {
	// GIVEN
	stub, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("certifier", nil)
	sellingFirm := &chaincode.Producer{ID: "oscar"}
	sellingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(sellingFirm)
	ctx.CreateProductionReturns(fmt.Errorf("failed adding production"))
	stub.PutStateReturns(nil)

	// WHEN
	err := contract.ProducerProduction(ctx, "oscar", 10, "1/1", "1", "energy",
		"greenhouse", "12 co2e", 1)

	// THEN
	require.EqualError(t, err, "failed adding production")
}

func Test_WHEN_payingProduction_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_production(i, "oscar", false)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	producingFirm := &chaincode.Producer{ID: "oscar"}
	producingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(producingFirm)
	ctx.GetHighThroughReturns(200, nil)

	// WHEN
	amount, err := contract.PayForProduction(ctx, "oscar")

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, 196, amount)
}

func Test_WHEN_notProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("certifier", nil)

	// WHEN
	_, err := contract.PayForProduction(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "err: only producer can pay for production")
}

func Test_WHEN_notOwnedDebt_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("john", nil)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_production(i, "oscar", false)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}

	// WHEN
	_, err := contract.PayForProduction(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "err: do not own carbon debt with id: oscar")
}

func Test_WHEN_alreadyPaid_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_function_production(i, "oscar", true)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	producingFirm := &chaincode.Producer{ID: "oscar"}
	producingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(producingFirm)
	ctx.GetHighThroughReturns(200, nil)

	// WHEN
	_, err := contract.PayForProduction(ctx, "oscar")

	// THEN
	require.EqualError(t, err, "err: already paid for production")
}

func Test_WHEN_redeemChip_THEN_SUCCESS(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	producingFirm := &chaincode.Producer{ID: "oscar"}
	producingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(producingFirm)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_amount_chip(i, "oscar", true)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	ctx.GetHighThroughReturns(200, nil)

	// WHEN
	amount, err := contract.RedeemChip(ctx, 15)

	// THEN
	require.Nil(t, err)
	require.EqualValues(t, amount, 215)
}

func Test_WHEN_redeemChipNotProducer_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("certifier", nil)
	ctx.GetUserIdReturns("oscar", nil)

	// WHEN
	amount, err := contract.RedeemChip(ctx, 15)

	// THEN
	require.EqualError(t, err, "err: only producers can redeem offer chips")
	require.EqualValues(t, 0, amount)
}

func Test_WHEN_redeemChipNotValid_THEN_FAILURE(t *testing.T) {
	// GIVEN
	_, ctx, contract := GetTestStubs()
	ctx.GetUserTypeReturns("producer", nil)
	ctx.GetUserIdReturns("oscar", nil)
	producingFirm := &chaincode.Producer{ID: "oscar"}
	producingFirm.InsertContext(ctx)
	ctx.GetProducerReturns(producingFirm)
	ctx.GetResultStub = func(s string, i interface{}) error {
		callBackFunc, doc := mock_amount_chip(i, "oscar", false)
		callBackFunc.Call([]reflect.Value{*doc})
		return nil
	}
	ctx.GetHighThroughReturns(200, nil)

	// WHEN
	_, err := contract.RedeemChip(ctx, 15)

	// THEN
	require.EqualError(t, err, "err: the offer chip is no longer valid")
}
