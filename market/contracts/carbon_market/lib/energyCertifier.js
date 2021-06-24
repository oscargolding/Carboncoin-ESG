/**
 * For having an energy certifier with the exchange of tokens.
 * Written by Oscar Golding - UNSW
 * -- Not the primary objective of the thesis (just an exploration)
 */

const { Contract } = require('fabric-contract-api');

const SMALL = 'small';
const MEDIUM = 'medium';
const LARGE = 'large';

/**
 * Energy certifier to interact with the blockchain.
 */
class EnergyCertifier extends Contract {
  /**
   * Initialise the ledger with a given state.
   * @param {ctx} ctx the context for the smart contract
   */
  async InitLedger(ctx) {
    const assets = [
      {
        ID: 'oscarIndustry',
        size: LARGE,
      },
      {
        ID: 'rioTinto',
        size: MEDIUM,
      },
      {
        ID: 'smallFirm',
        size: SMALL,
      },
      {
        ID: 'largeFirm',
        size: LARGE,
      },
    ];
    await Promise.all(assets.map(async (asset) => {
      await ctx.stub.putState(asset.ID, Buffer.from(JSON.stringify(asset)));
      console.info(`Firm ${asset.ID} initialised`);
    }));
  }

  /**
   * Get the size of the firm - verified on the chain.
   * @param {context} ctx the context of the firm
   * @param {*} firm to call
   */
  async FirmSize(ctx, firm) {
    const recoveredFirm = await ctx.stub.getState(firm);
    if (!recoveredFirm || recoveredFirm.length === 0) {
      return { size: SMALL };
    }
    const json = JSON.parse(recoveredFirm.toString());
    return { size: json.size };
  }
}

module.exports = EnergyCertifier;
