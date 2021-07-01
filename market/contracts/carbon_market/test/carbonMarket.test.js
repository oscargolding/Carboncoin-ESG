/*
 * Oscar Golding
 *
 * SPDX-License-Identifier: Apache-2.0
*/
const sinon = require('sinon');
const chai = require('chai');
const sinonChai = require('sinon-chai');
const stubs = require('./stub');

const { expect } = chai;

const CarbonMarket = require('../lib/carbonMarket');

const { assert } = sinon;
chai.use(sinonChai);

describe('Carbon Market basic tests', () => {
  let transactionContext;
  let chaincodeStub;
  let clientIdentity;

  // Run before all
  beforeEach(() => {
    ({ transactionContext, chaincodeStub, clientIdentity } = stubs.stub());
  });

  // For setting the firm size
  const firmSize = (size) => {
    chaincodeStub.invokeChaincode.callsFake(async () => Promise.resolve(
      {
        payload: Buffer.from(JSON.stringify({ size })),
      },
    ));
  };

  /**
   * Admin permissions are being tested
   */
  const withRole = (role) => {
    clientIdentity.getID
      .callsFake(() => `/CN=${role}::/C=`);
    clientIdentity.getAttributeValue.callsFake(() => role);
  };

  // add a producer
  const addProducer = async () => {
    withRole('admin');
    firmSize('large');
    const carbonMarket = new CarbonMarket();
    await carbonMarket.AddProducer(transactionContext, 'oscar');
    return carbonMarket;
  };

  // Start the tests
  describe('Test Producer Added', () => {
    it('Should allow for a producer to be added', async () => {
      // GIVEN
      withRole('admin');
      firmSize('small');
      const carbonMarket = new CarbonMarket();
      // WHEN
      await carbonMarket.AddProducer(transactionContext, 'oscar');
      const ret = JSON.parse(
        (await chaincodeStub.getState('oscar')).toString(),
      );
      // THEN
      expect(ret).to.eql({ producerId: 'oscar', tokens: 100, sellable: 100 });
    });

    it('Should fail when the producer exists', async () => {
      // GIVEN
      withRole('admin');
      const carbonMarket = new CarbonMarket();
      firmSize('medium');
      await carbonMarket.AddProducer(transactionContext, 'oscar');
      // WHEN
      try {
        await carbonMarket.AddProducer(transactionContext, 'oscar');
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
      firmSize('medium');
      // WHEN
      try {
        await carbonMarket.AddProducer(transactionContext, 'oscar');
        assert.fail('Should have failed due to the presence of producer');
      } catch (err) {
        // THEN
        expect(err.message).to.equal('Only admins can add a producer');
      }
    });
  });

  describe('Allow for the retrieval of balance', () => {
    it('Tests a simple balance retrieval', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      // WHEN
      const balance = await carbonMarket.GetBalance(transactionContext,
        'oscar');
      // THEN
      expect(balance).to.equal(300);
    });

    it('Allows a producer to perform the access', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('producer');
      // WHEN
      const balance = await carbonMarket.GetBalance(transactionContext,
        'oscar');
      // THEN
      expect(balance).to.equal(300);
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

  describe('Allows for the retrieval of sellable', () => {
    it('Allows for a producer to check sellable', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      // WHEN
      const sellable = await carbonMarket.GetSellable(transactionContext,
        'oscar');
      // THEN
      expect(sellable).to.equal(300);
    });
  });

  describe('Allows for the participant to make offer', () => {
    it('Allows for the adding of a token deal', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('oscar');
      clientIdentity.getAttributeValue.callsFake(() => 'producer');
      // WHEN
      await carbonMarket.AddOffer(transactionContext, 'oscar', 100, 100);
      // THEN
      const sellable = await carbonMarket.GetSellable(transactionContext,
        'oscar');
      console.log(sellable);
      expect(sellable).to.equal(200);
    });
    it('Does not allow more than holding', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('oscar');
      clientIdentity.getAttributeValue.callsFake(() => 'producer');
      try {
        // WHEN
        await carbonMarket.AddOffer(transactionContext, 'oscar', 50, 1000);
        assert.fail('The offer is too large');
      } catch (err) {
        // THEN
        expect(err.message).to.equal(
          'Producer oscar does not have enough tokens',
        );
      }
    });
    it('Does not allow with incorrect access', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('oscar');
      try {
        // WHEN
        await carbonMarket.AddOffer(transactionContext, 'oscar', 50, 1000);
        assert.fail('Admin called');
      } catch (err) {
        // THEN
        expect(err.message).to.equal(
          'Incorrect credentials for selling Carboncoin',
        );
      }
    });
    it('Does not allow offer for bad producer', async () => {
      // GIVEN
      const carbonMarket = await addProducer();
      withRole('oscar');
      try {
        // WHEN
        await carbonMarket.AddOffer(transactionContext, 'james', 10, 10);
        assert.fail('User does not exist');
      } catch (err) {
        // THEN
        expect(err.message).to.equal(
          'The producer with the name james does not exist',
        );
      }
    });
  });
});
