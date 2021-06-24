const chai = require('chai');
const sinonChai = require('sinon-chai');
const stubs = require('./stub');

const { expect } = chai;

const EnergyCertifier = require('../lib/energyCertifier');

chai.use(sinonChai);

describe('Energy certifier basic tests', () => {
  let transactionContext;
  let chaincodeStub;
  const SMALL = 'small';
  const LARGE = 'large';

  // Run the before all to perform stubbing
  beforeEach(() => {
    ({ transactionContext, chaincodeStub } = stubs.stub());
  });

  describe('Test the large firms are present', () => {
    it('Should have some of the firms', async () => {
      // GIVEN
      const energyCertifier = new EnergyCertifier();
      // WHEN
      await energyCertifier.InitLedger(transactionContext);
      const firstRet = JSON.parse(
        (await chaincodeStub.getState('oscarIndustry')).toString(),
      );
      const secondRet = JSON.parse(
        (await chaincodeStub.getState('smallFirm')).toString(),
      );
      // THEN
      expect(firstRet).to.deep.eq({ ID: 'oscarIndustry', size: 'large' });
      expect(secondRet).to.deep.eq({ ID: 'smallFirm', size: 'small' });
    });
  });

  describe('Should correctly get the size of the firm', () => {
    it('Should allow for a firm that exists', async () => {
      // GIVEN
      const energyCertifier = new EnergyCertifier();
      await energyCertifier.InitLedger(transactionContext);
      // WHEN
      const size = await energyCertifier.FirmSize(transactionContext,
        'oscarIndustry');
      // THEN
      expect(JSON.parse(JSON.stringify(size))).to.deep.eq({ size: LARGE });
    });
    it('Should allow for checking a firm not present', async () => {
      // GIVEN
      const energyCertifier = new EnergyCertifier();
      await energyCertifier.InitLedger(transactionContext);
      // WHEN
      const size = await energyCertifier.FirmSize(transactionContext,
        'elizabethIndustries');
      console.log(size);
      // THEN
      expect(JSON.parse(JSON.stringify(size))).to.deep.eq({ size: SMALL });
    });
  });
});
