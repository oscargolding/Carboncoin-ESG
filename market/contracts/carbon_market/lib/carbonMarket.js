/**
 * For creating a carbon market with the exchange of tokens.
 * Written by Oscar Golding - UNSW
 */

'use strict';

const { Contract } = require('fabric-contract-api');

class CarbonMarket extends Contract {

  /**
   * Add a hydrogen producer to the market with a given number of tokens.
   * @param {context} ctx the transaction context.
   * @param {defaultTokens} defaultTokens for the market.
   */
  async AddProducer(ctx, producerId, defaultTokens) {
    const producerExists = await this.checkProducer(ctx, producerId);
    if (producerExists) {
      throw new Error(`The producer with name ${producerId} exists`);
    } else {
      // Here a producer does not exists
      let producer = {
        producerId: producerId,
        tokens: defaultTokens
      };
      await ctx.stub.putState(producer.producerId,
        Buffer.from(JSON.stringify(producer)));
      console.log(`Producer ${producer.producerId} Initialised`);
    }
  }

  async checkProducer(ctx, producerId) {
    const producerJson = await ctx.stub.getState(producerId);
    if (!producerJson || producerJson.length === 0) {
      return false;
    }
    return true;
  }
}

module.exports = CarbonMarket;