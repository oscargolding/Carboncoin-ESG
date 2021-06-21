/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
*/
const sinon = require('sinon');
const chai = require('chai');
const sinonChai = require('sinon-chai');

const { expect } = chai;

const { Context } = require('fabric-contract-api');
const { ChaincodeStub, ClientIdentity } = require('fabric-shim');

const CarbonMarket = require('../lib/carbonMarket');

const { assert } = sinon;
chai.use(sinonChai);

describe('Carbon Market basic tests', () => {
  let transactionContext;
  let chaincodeStub;
  let clientIdentity;

  // Run before all
  beforeEach(() => {
    transactionContext = new Context();

    chaincodeStub = sinon.createStubInstance(ChaincodeStub);
    clientIdentity = sinon.createStubInstance(ClientIdentity);
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
  });

  /**
   * Admin permissions are being tested
   */
  const withRole = (role) => {
    clientIdentity.getID
      .callsFake(() => `/CN=${role}::/C=`);
    clientIdentity.getAttributeValue.callsFake(() => role);
  };

  // Start the tests
  describe('Test Producer Added', () => {
    it('Should allow for a producer to be added', async () => {
      // GIVEN
      withRole('admin');
      const carbonMarket = new CarbonMarket();
      // WHEN
      await carbonMarket.AddProducer(transactionContext, 'oscar', 10);
      const ret = JSON.parse(
        (await chaincodeStub.getState('oscar')).toString(),
      );
      // THEN
      expect(ret).to.eql({ producerId: 'oscar', tokens: 10 });
    });

    it('Should fail when the producer exists', async () => {
      // GIVEN
      withRole('admin');
      const carbonMarket = new CarbonMarket();
      await carbonMarket.AddProducer(transactionContext, 'oscar', 30);
      // WHEN
      try {
        await carbonMarket.AddProducer(transactionContext, 'oscar', 50);
        assert.fail('Should have failed with same producer');
      } catch (err) {
        // THEN
        expect(err.message).to.equal('The producer with name oscar exists');
      }
    });

    it('Should not allow a producer to call', async () => {
      // GIVEN
      withRole('producer');
      const carbonMarket = new CarbonMarket();
      // WHEN
      try {
        await carbonMarket.AddProducer(transactionContext, 'oscar', 50);
        assert.fail('Should have failed due to the presence of producer');
      } catch (err) {
        // THEN
        expect(err.message).to.equal('Only admins can add a producer');
      }
    });
  });

  describe('Allow for the retrieval of balance', () => {
    // For adding a producer to the market
    const addProducer = async () => {
      withRole('admin');
      const carbonMarket = new CarbonMarket();
      await carbonMarket.AddProducer(transactionContext, 'oscar', 100);
      return carbonMarket;
    };

    it('Tests a simple balance retrieval', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      // WHEN
      const balance = await carbonMarket.GetBalance(transactionContext,
        'oscar');
      // THEN
      expect(balance).to.equal(100);
    });

    it('Allows a producer to perform the access', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('producer');
      // WHEN
      const balance = await carbonMarket.GetBalance(transactionContext,
        'oscar');
      // THEN
      expect(balance).to.equal(100);
    });

    it('Throws error on wrong user supplied', async () => {
      // GIVEN
      const carbonMarket = new CarbonMarket();
      try {
        // WHEN
        await carbonMarket.GetBalance(transactionContext, 'james');
        assert.fail('The user should not exist');
      } catch (err) {
        // THEN
        expect(err.message).to.equal('The producer requested does not exist');
      }
    });
  });
});
