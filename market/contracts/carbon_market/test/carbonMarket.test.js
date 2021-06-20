/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
*/

'use strict';
const sinon = require('sinon');
const chai = require('chai');
const sinonChai = require('sinon-chai');
const expect = chai.expect;

const { Context } = require('fabric-contract-api');
const { ChaincodeStub } = require('fabric-shim');

const CarbonMarket = require('../lib/carbonMarket.js');

let assert = sinon.assert;
chai.use(sinonChai);

describe('Carbon Market basic tests', () => {
  let transactionContext, chaincodeStub;

  // Run before all
  beforeEach(() => {
    transactionContext = new Context();

    chaincodeStub = sinon.createStubInstance(ChaincodeStub);
    transactionContext.setChaincodeStub(chaincodeStub);

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
          const copied = Object.assign({}, chaincodeStub.states);

          for (let key in copied) {
            yield { value: copied[key] };
          }
        }
      }

      return Promise.resolve(internalGetStateByRange());
    });
  });

  // Start the tests
  describe('Test Producer Added', () => {
    it('Should allow for a producer to be added', async () => {
      // GIVEN
      const carbonMarket = new CarbonMarket();
      // WHEN
      await carbonMarket.AddProducer(transactionContext, 'oscar', 10);
      let ret = JSON.parse((await chaincodeStub.getState('oscar')).toString());
      // THEN
      expect(ret).to.eql(Object.assign({},
        { producerId: 'oscar', tokens: 10 }));
    });

    it('Should fail when the producer exists', async () => {
      const carbonMarket = new CarbonMarket();
      await carbonMarket.AddProducer(transactionContext, 'oscar', 30);
      try {
        await carbonMarket.AddProducer(transactionContext, 'oscar', 50);
        assert.fail('Should have failed with same producer');
      } catch (err) {
        expect(err.message).to.equal('The producer with name oscar exists');
      }
    });
  });
});
