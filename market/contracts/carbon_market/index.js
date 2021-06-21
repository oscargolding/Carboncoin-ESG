/**
 * Smart contracts for the carbon market
 */
const carbonMarket = require('./lib/carbonMarket');

module.exports.CarbonMarket = carbonMarket;
module.exports.contracts = [carbonMarket];
