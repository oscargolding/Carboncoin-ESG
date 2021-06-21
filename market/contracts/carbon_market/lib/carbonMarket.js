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
  async AddProducer(ctx, producerId, defaultTokens) {
    const producerExists = await CarbonMarket.checkProducer(ctx, producerId);
    if (producerExists) {
      throw new Error(`The producer with name ${producerId} exists`);
    } else {
      // Here a producer does not exists
      const userType = await CarbonMarket.getCurrentUserType(ctx);
      if (userType !== 'admin') {
        throw new Error('Only admins can add a producer');
      }
      const producer = {
        producerId,
        tokens: defaultTokens,
      };
      await ctx.stub.putState(producer.producerId,
        Buffer.from(JSON.stringify(producer)));
      console.log(`Producer ${producer.producerId} Initialised`);
    }
  }

  /**
   * Get the balance of a particular producer - open with a purpose.
   * @param {context} ctx the context associated with transaction.
   * @param {producerId} producerId the id of the producer querying.
   */
  async GetBalance(ctx, producerId) {
    const retrievedProducer = await ctx.stub.getState(producerId);
    if (!retrievedProducer || retrievedProducer.json === 0) {
      throw new Error('The producer requested does not exist');
    }
    // Get the tokens associated with the user
    const ret = JSON.parse(retrievedProducer.toString());
    return ret.tokens;
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
}

module.exports = CarbonMarket;
