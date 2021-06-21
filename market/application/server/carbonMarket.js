/**
 * Carbon market as written by Oscar Golding
 */
const express = require('express');
const utils = require('./utils');

const carbonMarketRouter = express.Router();
const { InputError } = require('./errors/inputError');
const { AccessError } = require('./errors/accessError');

const catchErrors = (fn) => async (req, res) => {
  try {
    await fn(req, res);
  } catch (err) {
    if (err instanceof InputError) {
      res.status(400).send({ error: err.message });
    } else if (err instanceof AccessError) {
      res.status(403).send({ error: err.message });
    } else {
      console.log(err);
      res.status(500).send({ error: 'A system error ocurred' });
    }
  }
};

/**
 * When a producer wants to register for the first time.
 */
carbonMarketRouter.post('/admin/auth/register', catchErrors(async (req, res) => {
  const { email, password } = req.body;
  await utils.registerProducer(email, password);
  await utils.registerTokens(email);
  return res.json({ mssg: 'Successful register' });
}));

/**
 * For retrieving the balance of a particular user
 */
carbonMarketRouter.get('/token/balance/:id', catchErrors(async (req, res) => {
  const user = req.params.id;
  const balance = await utils.retrieveBalance(user);
  return res.json({ balance });
}));

// Export the carbon market router
module.exports = carbonMarketRouter;
