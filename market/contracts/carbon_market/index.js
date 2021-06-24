/**
 * Smart contracts for the carbon market
 */
const carbonMarket = require('./lib/carbonMarket');
const energyCertifier = require('./lib/energyCertifier');

module.exports.CarbonMarket = carbonMarket;
module.exports.EnergyCertifier = energyCertifier;
module.exports.contracts = [carbonMarket, energyCertifier];
