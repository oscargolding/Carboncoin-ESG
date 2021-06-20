/**
 * Smart contracts for the carbon market
 */
'use strict';

const carbonMarket = require('./lib/carbonMarket');

module.exports.CarbonMarket = carbonMarket;
module.exports.contracts = [carbonMarket];