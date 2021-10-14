// Code generated by counterfeiter. DO NOT EDIT.
package chaincodefakes

import (
	"market/chaincode"
	"sync"

	"github.com/hyperledger/fabric-chaincode-go/pkg/cid"
	"github.com/hyperledger/fabric-chaincode-go/shim"
)

type CustomContex struct {
	CheckProducerStub        func(string) (bool, error)
	checkProducerMutex       sync.RWMutex
	checkProducerArgsForCall []struct {
		arg1 string
	}
	checkProducerReturns struct {
		result1 bool
		result2 error
	}
	checkProducerReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	CreateChipStub        func(string, int) error
	createChipMutex       sync.RWMutex
	createChipArgsForCall []struct {
		arg1 string
		arg2 int
	}
	createChipReturns struct {
		result1 error
	}
	createChipReturnsOnCall map[int]struct {
		result1 error
	}
	CreateOfferStub        func(string, int, int, string, int) error
	createOfferMutex       sync.RWMutex
	createOfferArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
		arg4 string
		arg5 int
	}
	createOfferReturns struct {
		result1 error
	}
	createOfferReturnsOnCall map[int]struct {
		result1 error
	}
	CreateProductionStub        func(string, int, string, string, bool, bool, string, string, string, int) error
	createProductionMutex       sync.RWMutex
	createProductionArgsForCall []struct {
		arg1  string
		arg2  int
		arg3  string
		arg4  string
		arg5  bool
		arg6  bool
		arg7  string
		arg8  string
		arg9  string
		arg10 int
	}
	createProductionReturns struct {
		result1 error
	}
	createProductionReturnsOnCall map[int]struct {
		result1 error
	}
	GetClientIdentityStub        func() cid.ClientIdentity
	getClientIdentityMutex       sync.RWMutex
	getClientIdentityArgsForCall []struct {
	}
	getClientIdentityReturns struct {
		result1 cid.ClientIdentity
	}
	getClientIdentityReturnsOnCall map[int]struct {
		result1 cid.ClientIdentity
	}
	GetHighThroughStub        func(string) (int, error)
	getHighThroughMutex       sync.RWMutex
	getHighThroughArgsForCall []struct {
		arg1 string
	}
	getHighThroughReturns struct {
		result1 int
		result2 error
	}
	getHighThroughReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	GetProducerStub        func(string) *chaincode.Producer
	getProducerMutex       sync.RWMutex
	getProducerArgsForCall []struct {
		arg1 string
	}
	getProducerReturns struct {
		result1 *chaincode.Producer
	}
	getProducerReturnsOnCall map[int]struct {
		result1 *chaincode.Producer
	}
	GetResultStub        func(string, interface{}) error
	getResultMutex       sync.RWMutex
	getResultArgsForCall []struct {
		arg1 string
		arg2 interface{}
	}
	getResultReturns struct {
		result1 error
	}
	getResultReturnsOnCall map[int]struct {
		result1 error
	}
	GetStubStub        func() shim.ChaincodeStubInterface
	getStubMutex       sync.RWMutex
	getStubArgsForCall []struct {
	}
	getStubReturns struct {
		result1 shim.ChaincodeStubInterface
	}
	getStubReturnsOnCall map[int]struct {
		result1 shim.ChaincodeStubInterface
	}
	GetUserIdStub        func() (string, error)
	getUserIdMutex       sync.RWMutex
	getUserIdArgsForCall []struct {
	}
	getUserIdReturns struct {
		result1 string
		result2 error
	}
	getUserIdReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	GetUserTypeStub        func() (string, error)
	getUserTypeMutex       sync.RWMutex
	getUserTypeArgsForCall []struct {
	}
	getUserTypeReturns struct {
		result1 string
		result2 error
	}
	getUserTypeReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	IteratorResultsStub        func(shim.StateQueryIteratorInterface, interface{}) error
	iteratorResultsMutex       sync.RWMutex
	iteratorResultsArgsForCall []struct {
		arg1 shim.StateQueryIteratorInterface
		arg2 interface{}
	}
	iteratorResultsReturns struct {
		result1 error
	}
	iteratorResultsReturnsOnCall map[int]struct {
		result1 error
	}
	OfferStringGeneratorStub        func(string, bool, string) string
	offerStringGeneratorMutex       sync.RWMutex
	offerStringGeneratorArgsForCall []struct {
		arg1 string
		arg2 bool
		arg3 string
	}
	offerStringGeneratorReturns struct {
		result1 string
	}
	offerStringGeneratorReturnsOnCall map[int]struct {
		result1 string
	}
	UpdateHighThroughStub        func(string, string, int) error
	updateHighThroughMutex       sync.RWMutex
	updateHighThroughArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 int
	}
	updateHighThroughReturns struct {
		result1 error
	}
	updateHighThroughReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CustomContex) CheckProducer(arg1 string) (bool, error) {
	fake.checkProducerMutex.Lock()
	ret, specificReturn := fake.checkProducerReturnsOnCall[len(fake.checkProducerArgsForCall)]
	fake.checkProducerArgsForCall = append(fake.checkProducerArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.CheckProducerStub
	fakeReturns := fake.checkProducerReturns
	fake.recordInvocation("CheckProducer", []interface{}{arg1})
	fake.checkProducerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CustomContex) CheckProducerCallCount() int {
	fake.checkProducerMutex.RLock()
	defer fake.checkProducerMutex.RUnlock()
	return len(fake.checkProducerArgsForCall)
}

func (fake *CustomContex) CheckProducerCalls(stub func(string) (bool, error)) {
	fake.checkProducerMutex.Lock()
	defer fake.checkProducerMutex.Unlock()
	fake.CheckProducerStub = stub
}

func (fake *CustomContex) CheckProducerArgsForCall(i int) string {
	fake.checkProducerMutex.RLock()
	defer fake.checkProducerMutex.RUnlock()
	argsForCall := fake.checkProducerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *CustomContex) CheckProducerReturns(result1 bool, result2 error) {
	fake.checkProducerMutex.Lock()
	defer fake.checkProducerMutex.Unlock()
	fake.CheckProducerStub = nil
	fake.checkProducerReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) CheckProducerReturnsOnCall(i int, result1 bool, result2 error) {
	fake.checkProducerMutex.Lock()
	defer fake.checkProducerMutex.Unlock()
	fake.CheckProducerStub = nil
	if fake.checkProducerReturnsOnCall == nil {
		fake.checkProducerReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.checkProducerReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) CreateChip(arg1 string, arg2 int) error {
	fake.createChipMutex.Lock()
	ret, specificReturn := fake.createChipReturnsOnCall[len(fake.createChipArgsForCall)]
	fake.createChipArgsForCall = append(fake.createChipArgsForCall, struct {
		arg1 string
		arg2 int
	}{arg1, arg2})
	stub := fake.CreateChipStub
	fakeReturns := fake.createChipReturns
	fake.recordInvocation("CreateChip", []interface{}{arg1, arg2})
	fake.createChipMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) CreateChipCallCount() int {
	fake.createChipMutex.RLock()
	defer fake.createChipMutex.RUnlock()
	return len(fake.createChipArgsForCall)
}

func (fake *CustomContex) CreateChipCalls(stub func(string, int) error) {
	fake.createChipMutex.Lock()
	defer fake.createChipMutex.Unlock()
	fake.CreateChipStub = stub
}

func (fake *CustomContex) CreateChipArgsForCall(i int) (string, int) {
	fake.createChipMutex.RLock()
	defer fake.createChipMutex.RUnlock()
	argsForCall := fake.createChipArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *CustomContex) CreateChipReturns(result1 error) {
	fake.createChipMutex.Lock()
	defer fake.createChipMutex.Unlock()
	fake.CreateChipStub = nil
	fake.createChipReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) CreateChipReturnsOnCall(i int, result1 error) {
	fake.createChipMutex.Lock()
	defer fake.createChipMutex.Unlock()
	fake.CreateChipStub = nil
	if fake.createChipReturnsOnCall == nil {
		fake.createChipReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createChipReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) CreateOffer(arg1 string, arg2 int, arg3 int, arg4 string, arg5 int) error {
	fake.createOfferMutex.Lock()
	ret, specificReturn := fake.createOfferReturnsOnCall[len(fake.createOfferArgsForCall)]
	fake.createOfferArgsForCall = append(fake.createOfferArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
		arg4 string
		arg5 int
	}{arg1, arg2, arg3, arg4, arg5})
	stub := fake.CreateOfferStub
	fakeReturns := fake.createOfferReturns
	fake.recordInvocation("CreateOffer", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.createOfferMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) CreateOfferCallCount() int {
	fake.createOfferMutex.RLock()
	defer fake.createOfferMutex.RUnlock()
	return len(fake.createOfferArgsForCall)
}

func (fake *CustomContex) CreateOfferCalls(stub func(string, int, int, string, int) error) {
	fake.createOfferMutex.Lock()
	defer fake.createOfferMutex.Unlock()
	fake.CreateOfferStub = stub
}

func (fake *CustomContex) CreateOfferArgsForCall(i int) (string, int, int, string, int) {
	fake.createOfferMutex.RLock()
	defer fake.createOfferMutex.RUnlock()
	argsForCall := fake.createOfferArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *CustomContex) CreateOfferReturns(result1 error) {
	fake.createOfferMutex.Lock()
	defer fake.createOfferMutex.Unlock()
	fake.CreateOfferStub = nil
	fake.createOfferReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) CreateOfferReturnsOnCall(i int, result1 error) {
	fake.createOfferMutex.Lock()
	defer fake.createOfferMutex.Unlock()
	fake.CreateOfferStub = nil
	if fake.createOfferReturnsOnCall == nil {
		fake.createOfferReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createOfferReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) CreateProduction(arg1 string, arg2 int, arg3 string, arg4 string, arg5 bool, arg6 bool, arg7 string, arg8 string, arg9 string, arg10 int) error {
	fake.createProductionMutex.Lock()
	ret, specificReturn := fake.createProductionReturnsOnCall[len(fake.createProductionArgsForCall)]
	fake.createProductionArgsForCall = append(fake.createProductionArgsForCall, struct {
		arg1  string
		arg2  int
		arg3  string
		arg4  string
		arg5  bool
		arg6  bool
		arg7  string
		arg8  string
		arg9  string
		arg10 int
	}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10})
	stub := fake.CreateProductionStub
	fakeReturns := fake.createProductionReturns
	fake.recordInvocation("CreateProduction", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10})
	fake.createProductionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9, arg10)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) CreateProductionCallCount() int {
	fake.createProductionMutex.RLock()
	defer fake.createProductionMutex.RUnlock()
	return len(fake.createProductionArgsForCall)
}

func (fake *CustomContex) CreateProductionCalls(stub func(string, int, string, string, bool, bool, string, string, string, int) error) {
	fake.createProductionMutex.Lock()
	defer fake.createProductionMutex.Unlock()
	fake.CreateProductionStub = stub
}

func (fake *CustomContex) CreateProductionArgsForCall(i int) (string, int, string, string, bool, bool, string, string, string, int) {
	fake.createProductionMutex.RLock()
	defer fake.createProductionMutex.RUnlock()
	argsForCall := fake.createProductionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6, argsForCall.arg7, argsForCall.arg8, argsForCall.arg9, argsForCall.arg10
}

func (fake *CustomContex) CreateProductionReturns(result1 error) {
	fake.createProductionMutex.Lock()
	defer fake.createProductionMutex.Unlock()
	fake.CreateProductionStub = nil
	fake.createProductionReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) CreateProductionReturnsOnCall(i int, result1 error) {
	fake.createProductionMutex.Lock()
	defer fake.createProductionMutex.Unlock()
	fake.CreateProductionStub = nil
	if fake.createProductionReturnsOnCall == nil {
		fake.createProductionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.createProductionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) GetClientIdentity() cid.ClientIdentity {
	fake.getClientIdentityMutex.Lock()
	ret, specificReturn := fake.getClientIdentityReturnsOnCall[len(fake.getClientIdentityArgsForCall)]
	fake.getClientIdentityArgsForCall = append(fake.getClientIdentityArgsForCall, struct {
	}{})
	stub := fake.GetClientIdentityStub
	fakeReturns := fake.getClientIdentityReturns
	fake.recordInvocation("GetClientIdentity", []interface{}{})
	fake.getClientIdentityMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) GetClientIdentityCallCount() int {
	fake.getClientIdentityMutex.RLock()
	defer fake.getClientIdentityMutex.RUnlock()
	return len(fake.getClientIdentityArgsForCall)
}

func (fake *CustomContex) GetClientIdentityCalls(stub func() cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = stub
}

func (fake *CustomContex) GetClientIdentityReturns(result1 cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = nil
	fake.getClientIdentityReturns = struct {
		result1 cid.ClientIdentity
	}{result1}
}

func (fake *CustomContex) GetClientIdentityReturnsOnCall(i int, result1 cid.ClientIdentity) {
	fake.getClientIdentityMutex.Lock()
	defer fake.getClientIdentityMutex.Unlock()
	fake.GetClientIdentityStub = nil
	if fake.getClientIdentityReturnsOnCall == nil {
		fake.getClientIdentityReturnsOnCall = make(map[int]struct {
			result1 cid.ClientIdentity
		})
	}
	fake.getClientIdentityReturnsOnCall[i] = struct {
		result1 cid.ClientIdentity
	}{result1}
}

func (fake *CustomContex) GetHighThrough(arg1 string) (int, error) {
	fake.getHighThroughMutex.Lock()
	ret, specificReturn := fake.getHighThroughReturnsOnCall[len(fake.getHighThroughArgsForCall)]
	fake.getHighThroughArgsForCall = append(fake.getHighThroughArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetHighThroughStub
	fakeReturns := fake.getHighThroughReturns
	fake.recordInvocation("GetHighThrough", []interface{}{arg1})
	fake.getHighThroughMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CustomContex) GetHighThroughCallCount() int {
	fake.getHighThroughMutex.RLock()
	defer fake.getHighThroughMutex.RUnlock()
	return len(fake.getHighThroughArgsForCall)
}

func (fake *CustomContex) GetHighThroughCalls(stub func(string) (int, error)) {
	fake.getHighThroughMutex.Lock()
	defer fake.getHighThroughMutex.Unlock()
	fake.GetHighThroughStub = stub
}

func (fake *CustomContex) GetHighThroughArgsForCall(i int) string {
	fake.getHighThroughMutex.RLock()
	defer fake.getHighThroughMutex.RUnlock()
	argsForCall := fake.getHighThroughArgsForCall[i]
	return argsForCall.arg1
}

func (fake *CustomContex) GetHighThroughReturns(result1 int, result2 error) {
	fake.getHighThroughMutex.Lock()
	defer fake.getHighThroughMutex.Unlock()
	fake.GetHighThroughStub = nil
	fake.getHighThroughReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) GetHighThroughReturnsOnCall(i int, result1 int, result2 error) {
	fake.getHighThroughMutex.Lock()
	defer fake.getHighThroughMutex.Unlock()
	fake.GetHighThroughStub = nil
	if fake.getHighThroughReturnsOnCall == nil {
		fake.getHighThroughReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.getHighThroughReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) GetProducer(arg1 string) *chaincode.Producer {
	fake.getProducerMutex.Lock()
	ret, specificReturn := fake.getProducerReturnsOnCall[len(fake.getProducerArgsForCall)]
	fake.getProducerArgsForCall = append(fake.getProducerArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.GetProducerStub
	fakeReturns := fake.getProducerReturns
	fake.recordInvocation("GetProducer", []interface{}{arg1})
	fake.getProducerMutex.Unlock()
	if stub != nil {
		return stub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) GetProducerCallCount() int {
	fake.getProducerMutex.RLock()
	defer fake.getProducerMutex.RUnlock()
	return len(fake.getProducerArgsForCall)
}

func (fake *CustomContex) GetProducerCalls(stub func(string) *chaincode.Producer) {
	fake.getProducerMutex.Lock()
	defer fake.getProducerMutex.Unlock()
	fake.GetProducerStub = stub
}

func (fake *CustomContex) GetProducerArgsForCall(i int) string {
	fake.getProducerMutex.RLock()
	defer fake.getProducerMutex.RUnlock()
	argsForCall := fake.getProducerArgsForCall[i]
	return argsForCall.arg1
}

func (fake *CustomContex) GetProducerReturns(result1 *chaincode.Producer) {
	fake.getProducerMutex.Lock()
	defer fake.getProducerMutex.Unlock()
	fake.GetProducerStub = nil
	fake.getProducerReturns = struct {
		result1 *chaincode.Producer
	}{result1}
}

func (fake *CustomContex) GetProducerReturnsOnCall(i int, result1 *chaincode.Producer) {
	fake.getProducerMutex.Lock()
	defer fake.getProducerMutex.Unlock()
	fake.GetProducerStub = nil
	if fake.getProducerReturnsOnCall == nil {
		fake.getProducerReturnsOnCall = make(map[int]struct {
			result1 *chaincode.Producer
		})
	}
	fake.getProducerReturnsOnCall[i] = struct {
		result1 *chaincode.Producer
	}{result1}
}

func (fake *CustomContex) GetResult(arg1 string, arg2 interface{}) error {
	fake.getResultMutex.Lock()
	ret, specificReturn := fake.getResultReturnsOnCall[len(fake.getResultArgsForCall)]
	fake.getResultArgsForCall = append(fake.getResultArgsForCall, struct {
		arg1 string
		arg2 interface{}
	}{arg1, arg2})
	stub := fake.GetResultStub
	fakeReturns := fake.getResultReturns
	fake.recordInvocation("GetResult", []interface{}{arg1, arg2})
	fake.getResultMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) GetResultCallCount() int {
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	return len(fake.getResultArgsForCall)
}

func (fake *CustomContex) GetResultCalls(stub func(string, interface{}) error) {
	fake.getResultMutex.Lock()
	defer fake.getResultMutex.Unlock()
	fake.GetResultStub = stub
}

func (fake *CustomContex) GetResultArgsForCall(i int) (string, interface{}) {
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	argsForCall := fake.getResultArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *CustomContex) GetResultReturns(result1 error) {
	fake.getResultMutex.Lock()
	defer fake.getResultMutex.Unlock()
	fake.GetResultStub = nil
	fake.getResultReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) GetResultReturnsOnCall(i int, result1 error) {
	fake.getResultMutex.Lock()
	defer fake.getResultMutex.Unlock()
	fake.GetResultStub = nil
	if fake.getResultReturnsOnCall == nil {
		fake.getResultReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.getResultReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) GetStub() shim.ChaincodeStubInterface {
	fake.getStubMutex.Lock()
	ret, specificReturn := fake.getStubReturnsOnCall[len(fake.getStubArgsForCall)]
	fake.getStubArgsForCall = append(fake.getStubArgsForCall, struct {
	}{})
	stub := fake.GetStubStub
	fakeReturns := fake.getStubReturns
	fake.recordInvocation("GetStub", []interface{}{})
	fake.getStubMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) GetStubCallCount() int {
	fake.getStubMutex.RLock()
	defer fake.getStubMutex.RUnlock()
	return len(fake.getStubArgsForCall)
}

func (fake *CustomContex) GetStubCalls(stub func() shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = stub
}

func (fake *CustomContex) GetStubReturns(result1 shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = nil
	fake.getStubReturns = struct {
		result1 shim.ChaincodeStubInterface
	}{result1}
}

func (fake *CustomContex) GetStubReturnsOnCall(i int, result1 shim.ChaincodeStubInterface) {
	fake.getStubMutex.Lock()
	defer fake.getStubMutex.Unlock()
	fake.GetStubStub = nil
	if fake.getStubReturnsOnCall == nil {
		fake.getStubReturnsOnCall = make(map[int]struct {
			result1 shim.ChaincodeStubInterface
		})
	}
	fake.getStubReturnsOnCall[i] = struct {
		result1 shim.ChaincodeStubInterface
	}{result1}
}

func (fake *CustomContex) GetUserId() (string, error) {
	fake.getUserIdMutex.Lock()
	ret, specificReturn := fake.getUserIdReturnsOnCall[len(fake.getUserIdArgsForCall)]
	fake.getUserIdArgsForCall = append(fake.getUserIdArgsForCall, struct {
	}{})
	stub := fake.GetUserIdStub
	fakeReturns := fake.getUserIdReturns
	fake.recordInvocation("GetUserId", []interface{}{})
	fake.getUserIdMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CustomContex) GetUserIdCallCount() int {
	fake.getUserIdMutex.RLock()
	defer fake.getUserIdMutex.RUnlock()
	return len(fake.getUserIdArgsForCall)
}

func (fake *CustomContex) GetUserIdCalls(stub func() (string, error)) {
	fake.getUserIdMutex.Lock()
	defer fake.getUserIdMutex.Unlock()
	fake.GetUserIdStub = stub
}

func (fake *CustomContex) GetUserIdReturns(result1 string, result2 error) {
	fake.getUserIdMutex.Lock()
	defer fake.getUserIdMutex.Unlock()
	fake.GetUserIdStub = nil
	fake.getUserIdReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) GetUserIdReturnsOnCall(i int, result1 string, result2 error) {
	fake.getUserIdMutex.Lock()
	defer fake.getUserIdMutex.Unlock()
	fake.GetUserIdStub = nil
	if fake.getUserIdReturnsOnCall == nil {
		fake.getUserIdReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getUserIdReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) GetUserType() (string, error) {
	fake.getUserTypeMutex.Lock()
	ret, specificReturn := fake.getUserTypeReturnsOnCall[len(fake.getUserTypeArgsForCall)]
	fake.getUserTypeArgsForCall = append(fake.getUserTypeArgsForCall, struct {
	}{})
	stub := fake.GetUserTypeStub
	fakeReturns := fake.getUserTypeReturns
	fake.recordInvocation("GetUserType", []interface{}{})
	fake.getUserTypeMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *CustomContex) GetUserTypeCallCount() int {
	fake.getUserTypeMutex.RLock()
	defer fake.getUserTypeMutex.RUnlock()
	return len(fake.getUserTypeArgsForCall)
}

func (fake *CustomContex) GetUserTypeCalls(stub func() (string, error)) {
	fake.getUserTypeMutex.Lock()
	defer fake.getUserTypeMutex.Unlock()
	fake.GetUserTypeStub = stub
}

func (fake *CustomContex) GetUserTypeReturns(result1 string, result2 error) {
	fake.getUserTypeMutex.Lock()
	defer fake.getUserTypeMutex.Unlock()
	fake.GetUserTypeStub = nil
	fake.getUserTypeReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) GetUserTypeReturnsOnCall(i int, result1 string, result2 error) {
	fake.getUserTypeMutex.Lock()
	defer fake.getUserTypeMutex.Unlock()
	fake.GetUserTypeStub = nil
	if fake.getUserTypeReturnsOnCall == nil {
		fake.getUserTypeReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.getUserTypeReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *CustomContex) IteratorResults(arg1 shim.StateQueryIteratorInterface, arg2 interface{}) error {
	fake.iteratorResultsMutex.Lock()
	ret, specificReturn := fake.iteratorResultsReturnsOnCall[len(fake.iteratorResultsArgsForCall)]
	fake.iteratorResultsArgsForCall = append(fake.iteratorResultsArgsForCall, struct {
		arg1 shim.StateQueryIteratorInterface
		arg2 interface{}
	}{arg1, arg2})
	stub := fake.IteratorResultsStub
	fakeReturns := fake.iteratorResultsReturns
	fake.recordInvocation("IteratorResults", []interface{}{arg1, arg2})
	fake.iteratorResultsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) IteratorResultsCallCount() int {
	fake.iteratorResultsMutex.RLock()
	defer fake.iteratorResultsMutex.RUnlock()
	return len(fake.iteratorResultsArgsForCall)
}

func (fake *CustomContex) IteratorResultsCalls(stub func(shim.StateQueryIteratorInterface, interface{}) error) {
	fake.iteratorResultsMutex.Lock()
	defer fake.iteratorResultsMutex.Unlock()
	fake.IteratorResultsStub = stub
}

func (fake *CustomContex) IteratorResultsArgsForCall(i int) (shim.StateQueryIteratorInterface, interface{}) {
	fake.iteratorResultsMutex.RLock()
	defer fake.iteratorResultsMutex.RUnlock()
	argsForCall := fake.iteratorResultsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *CustomContex) IteratorResultsReturns(result1 error) {
	fake.iteratorResultsMutex.Lock()
	defer fake.iteratorResultsMutex.Unlock()
	fake.IteratorResultsStub = nil
	fake.iteratorResultsReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) IteratorResultsReturnsOnCall(i int, result1 error) {
	fake.iteratorResultsMutex.Lock()
	defer fake.iteratorResultsMutex.Unlock()
	fake.IteratorResultsStub = nil
	if fake.iteratorResultsReturnsOnCall == nil {
		fake.iteratorResultsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.iteratorResultsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) OfferStringGenerator(arg1 string, arg2 bool, arg3 string) string {
	fake.offerStringGeneratorMutex.Lock()
	ret, specificReturn := fake.offerStringGeneratorReturnsOnCall[len(fake.offerStringGeneratorArgsForCall)]
	fake.offerStringGeneratorArgsForCall = append(fake.offerStringGeneratorArgsForCall, struct {
		arg1 string
		arg2 bool
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.OfferStringGeneratorStub
	fakeReturns := fake.offerStringGeneratorReturns
	fake.recordInvocation("OfferStringGenerator", []interface{}{arg1, arg2, arg3})
	fake.offerStringGeneratorMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) OfferStringGeneratorCallCount() int {
	fake.offerStringGeneratorMutex.RLock()
	defer fake.offerStringGeneratorMutex.RUnlock()
	return len(fake.offerStringGeneratorArgsForCall)
}

func (fake *CustomContex) OfferStringGeneratorCalls(stub func(string, bool, string) string) {
	fake.offerStringGeneratorMutex.Lock()
	defer fake.offerStringGeneratorMutex.Unlock()
	fake.OfferStringGeneratorStub = stub
}

func (fake *CustomContex) OfferStringGeneratorArgsForCall(i int) (string, bool, string) {
	fake.offerStringGeneratorMutex.RLock()
	defer fake.offerStringGeneratorMutex.RUnlock()
	argsForCall := fake.offerStringGeneratorArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CustomContex) OfferStringGeneratorReturns(result1 string) {
	fake.offerStringGeneratorMutex.Lock()
	defer fake.offerStringGeneratorMutex.Unlock()
	fake.OfferStringGeneratorStub = nil
	fake.offerStringGeneratorReturns = struct {
		result1 string
	}{result1}
}

func (fake *CustomContex) OfferStringGeneratorReturnsOnCall(i int, result1 string) {
	fake.offerStringGeneratorMutex.Lock()
	defer fake.offerStringGeneratorMutex.Unlock()
	fake.OfferStringGeneratorStub = nil
	if fake.offerStringGeneratorReturnsOnCall == nil {
		fake.offerStringGeneratorReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.offerStringGeneratorReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *CustomContex) UpdateHighThrough(arg1 string, arg2 string, arg3 int) error {
	fake.updateHighThroughMutex.Lock()
	ret, specificReturn := fake.updateHighThroughReturnsOnCall[len(fake.updateHighThroughArgsForCall)]
	fake.updateHighThroughArgsForCall = append(fake.updateHighThroughArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.UpdateHighThroughStub
	fakeReturns := fake.updateHighThroughReturns
	fake.recordInvocation("UpdateHighThrough", []interface{}{arg1, arg2, arg3})
	fake.updateHighThroughMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *CustomContex) UpdateHighThroughCallCount() int {
	fake.updateHighThroughMutex.RLock()
	defer fake.updateHighThroughMutex.RUnlock()
	return len(fake.updateHighThroughArgsForCall)
}

func (fake *CustomContex) UpdateHighThroughCalls(stub func(string, string, int) error) {
	fake.updateHighThroughMutex.Lock()
	defer fake.updateHighThroughMutex.Unlock()
	fake.UpdateHighThroughStub = stub
}

func (fake *CustomContex) UpdateHighThroughArgsForCall(i int) (string, string, int) {
	fake.updateHighThroughMutex.RLock()
	defer fake.updateHighThroughMutex.RUnlock()
	argsForCall := fake.updateHighThroughArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *CustomContex) UpdateHighThroughReturns(result1 error) {
	fake.updateHighThroughMutex.Lock()
	defer fake.updateHighThroughMutex.Unlock()
	fake.UpdateHighThroughStub = nil
	fake.updateHighThroughReturns = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) UpdateHighThroughReturnsOnCall(i int, result1 error) {
	fake.updateHighThroughMutex.Lock()
	defer fake.updateHighThroughMutex.Unlock()
	fake.UpdateHighThroughStub = nil
	if fake.updateHighThroughReturnsOnCall == nil {
		fake.updateHighThroughReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateHighThroughReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CustomContex) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkProducerMutex.RLock()
	defer fake.checkProducerMutex.RUnlock()
	fake.createChipMutex.RLock()
	defer fake.createChipMutex.RUnlock()
	fake.createOfferMutex.RLock()
	defer fake.createOfferMutex.RUnlock()
	fake.createProductionMutex.RLock()
	defer fake.createProductionMutex.RUnlock()
	fake.getClientIdentityMutex.RLock()
	defer fake.getClientIdentityMutex.RUnlock()
	fake.getHighThroughMutex.RLock()
	defer fake.getHighThroughMutex.RUnlock()
	fake.getProducerMutex.RLock()
	defer fake.getProducerMutex.RUnlock()
	fake.getResultMutex.RLock()
	defer fake.getResultMutex.RUnlock()
	fake.getStubMutex.RLock()
	defer fake.getStubMutex.RUnlock()
	fake.getUserIdMutex.RLock()
	defer fake.getUserIdMutex.RUnlock()
	fake.getUserTypeMutex.RLock()
	defer fake.getUserTypeMutex.RUnlock()
	fake.iteratorResultsMutex.RLock()
	defer fake.iteratorResultsMutex.RUnlock()
	fake.offerStringGeneratorMutex.RLock()
	defer fake.offerStringGeneratorMutex.RUnlock()
	fake.updateHighThroughMutex.RLock()
	defer fake.updateHighThroughMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CustomContex) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
