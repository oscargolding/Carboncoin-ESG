/**
 * Stubbing functions
 */
const { ChaincodeStub, ClientIdentity } = require('fabric-shim');
const { Context } = require('fabric-contract-api');
const sinon = require('sinon');

/**
 * Stub out the chaincode functions.
 * @returns stubs
 */
const createStubs = () => {
  const transactionContext = new Context();
  const chaincodeStub = sinon.createStubInstance(ChaincodeStub);
  const clientIdentity = sinon.createStubInstance(ClientIdentity);
  transactionContext.setChaincodeStub(chaincodeStub);
  transactionContext.setClientIdentity(clientIdentity);

  chaincodeStub.putState.callsFake((key, value) => {
    if (!chaincodeStub.states) {
      chaincodeStub.states = {};
    }
    chaincodeStub.states[key] = value;
  });

  chaincodeStub.getState.callsFake(async (key) => {
    let ret;
    if (chaincodeStub.states) {
      ret = chaincodeStub.states[key];
    }
    return Promise.resolve(ret);
  });

  chaincodeStub.deleteState.callsFake(async (key) => {
    if (chaincodeStub.states) {
      delete chaincodeStub.states[key];
    }
    return Promise.resolve(key);
  });

  chaincodeStub.getStateByRange.callsFake(async () => {
    function* internalGetStateByRange() {
      if (chaincodeStub.states) {
        // Shallow copy
        const copied = { ...chaincodeStub.states };
        const keys = Object.keys(copied);
        const values = Object.values(copied);
        for (let i = 0; i < keys.length; i += 1) {
          yield { value: values[i] };
        }
      }
    }

    return Promise.resolve(internalGetStateByRange());
  });
  return { transactionContext, chaincodeStub, clientIdentity };
};

module.exports.stub = createStubs;
