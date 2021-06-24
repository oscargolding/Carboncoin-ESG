/**
 * Smart contracts for the carbon market
 */
const energyCertifier = require('./lib/energyCertifier');

module.exports.EnergyCertifier = energyCertifier;
module.exports.contracts = [energyCertifier];
