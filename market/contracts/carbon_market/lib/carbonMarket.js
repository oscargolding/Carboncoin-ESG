/**
 * For creating a carbon market with the exchange of tokens.
 * Written by Oscar Golding - UNSW
 */

const { Contract } = require('fabric-contract-api');

class CarbonMarket extends Contract {
  /**
   * Add a hydrogen producer to the market with a given number of tokens.
   * @param {context} ctx the transaction context.
   * @param {defaultTokens} defaultTokens for the market.
   */
  async AddProducer(ctx, producerId) {
    const producerExists = await CarbonMarket.checkProducer(ctx, producerId);
    if (producerExists) {
      throw new Error(`The producer with name ${producerId} exists`);
    } else {
      // Here a producer does not exists
      const userType = await CarbonMarket.getCurrentUserType(ctx);
      if (userType !== 'admin') {
        throw new Error('Only admins can add a producer');
      }
      const ccArgs = ['FirmSize', producerId];
      const producerSize = await ctx.stub.invokeChaincode('EnergyCertifier',
        ccArgs);
      const jsonSize = JSON.parse(Buffer.from(producerSize.payload));
      const tokens = CarbonMarket.determineTokens(jsonSize.size);
      const producer = {
        producerId,
        tokens,
        sellable: tokens,
      };
      await ctx.stub.putState(producer.producerId,
        Buffer.from(JSON.stringify(producer)));
      console.log(`Producer ${producer.producerId} Initialised`);
    }
  }

  /**
   * Get the relevant producer.
   * @param {contex} ctx context
   * @param {*} producerId producerId
   * @returns the producer
   */
  async GetProducer(ctx, producerId) {
    const retrievedProducer = await ctx.stub.getState(producerId);
    if (!retrievedProducer || retrievedProducer.json === 0) {
      throw new Error('The producer requested does not exist');
    }
    // Get the tokens associated with the user
    const ret = JSON.parse(retrievedProducer.toString());
    return ret;
  }

  /**
   * Get the balance of a particular producer - open with a purpose.
   * @param {context} ctx the context associated with transaction.
   * @param {producerId} producerId the id of the producer querying.
   */
  async GetBalance(ctx, producerId) {
    const ret = await this.GetProducer(ctx, producerId);
    return ret.tokens;
  }

  /**
   * Get the amount of tokens allowed to be sold.
   * @param {contex} ctx transaction contex
   * @param {*} producerId producerId
   * @returns the sellable tokens
   */
  async GetSellable(ctx, producerId) {
    const ret = await this.GetProducer(ctx, producerId);
    return ret.sellable;
  }

  /**
   * Add an offer for the sale of tokens.
   * @param {ctx} ctx the transaction context
   * @param {*} producerId the producerId
   * @param {*} amount the amount to sell token for
   * @param {*} tokens the number of tokens to sell
   */
  async AddOffer(ctx, producerId, amountGiven, tokensGiven) {
    const producerExists = await CarbonMarket.checkProducer(ctx, producerId);
    const amount = Number(amountGiven);
    const tokens = Number(tokensGiven);
    if (!producerExists) {
      throw new Error(
        `The producer with the name ${producerId} does not exist`,
      );
    } else {
      // The producer exists
      const userType = await CarbonMarket.getCurrentUserType(ctx);
      const userId = await CarbonMarket.getCurrentUserId(ctx);
      console.log(userType);
      console.log(userId);
      if (userType !== 'producer' || userId !== producerId) {
        throw new Error('Incorrect credentials for selling Carboncoin');
      }
      const producerBalance = await this.GetSellable(ctx, producerId);
      if (producerBalance < tokens) {
        throw new Error(`Producer ${producerId} does not have enough tokens`);
      }
      console.log('>>> Creating the offer');
      // Can perform the adding of an offer here
      const offer = {
        producer: producerId,
        amount,
        token: tokens,
      };
      const key = `${producerId}~${amount}~${tokens}`;
      await ctx.stub.putState(key, Buffer.from(JSON.stringify(offer)));
      console.log('>>> Created the offer');
      const indexName = 'producer~id';
      console.log('>>> Creating Index');
      const producerIndex = await ctx.stub.createCompositeKey(indexName,
        [producerId, amountGiven, tokensGiven]);
      await ctx.stub.putState(producerIndex, Buffer.from('\u0000'));
      console.log('>>> Created Index');
      console.log('>>> Getting Producer');
      const producer = await this.GetProducer(ctx, producerId);
      producer.sellable -= tokens;
      await ctx.stub.putState(producerId,
        Buffer.from(JSON.stringify(producer)));
      console.log('>>> Got Producer And Added');
    }
  }

  /**
   * Check the producer exists.
   * @param {context} ctx context to check
   * @param {producerId} producerId the producer id to use
   * @returns true if the producer exists, false otherwise
   */
  static async checkProducer(ctx, producerId) {
    const producerJson = await ctx.stub.getState(producerId);
    if (!producerJson || producerJson.length === 0) {
      return false;
    }
    return true;
  }

  /**
  * getCurrentUserId
  * To be called by application to get the type for a user who is logged in
  *
  * @param {Context} ctx the transaction context
  * Usage:  getCurrentUserId ()
 */
  static async getCurrentUserId(ctx) {
    const id = [];
    id.push(ctx.clientIdentity.getID());
    const begin = id[0].indexOf('/CN=');
    const end = id[0].lastIndexOf('::/C=');
    const userid = id[0].substring(begin + 4, end);
    return userid;
  }

  /**
  * getCurrentUserType
  * To be called by application to get the type for a user who is logged in
  *
  * @param {context} ctx the transaction context
  * Usage:  getCurrentUserType ()
 */
  static async getCurrentUserType(ctx) {
    const userid = await CarbonMarket.getCurrentUserId(ctx);

    //  check user id;  if admin, return type = admin;
    //  else return value set for attribute "type" in certificate;
    if (userid === 'admin') {
      return userid;
    }
    return ctx.clientIdentity.getAttributeValue('usertype');
  }

  /**
   * Get the right amount of tokens for the firm.
   * @param {string} size of the firm
   * @returns token amount
   */
  static determineTokens(size) {
    switch (size) {
      case 'medium':
        return 200;
      case 'large':
        return 300;
      default:
        return 100;
    }
  }
}

module.exports = CarbonMarket;
