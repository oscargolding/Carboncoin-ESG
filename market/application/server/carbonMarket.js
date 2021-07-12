/**
 * Carbon market as written by Oscar Golding
 */
import express from 'express';
import { register, save, getEmailFromAuthorization } from './auth.js';
import utils from './utils.js';
import InputError from './errors/inputError.js';
import AccessError from './errors/accessError.js';

const carbonMarketRouter = express.Router();

/**
 * Catching errors.
 * @param {*} fn to catch errors and return to the user
 * @returns the result otherwise error
 */
const catchErrors = (fn) => async (req, res) => {
  try {
    await fn(req, res);
    save(); // Persist when using json storage
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
 * Wrapper to ensure auth is provided.
 * @param {*} fn to wrap around the auth request
 * @returns the result if auth is provided.
 */
const authed = (fn) => async (req, res) => {
  const email = getEmailFromAuthorization(req.header('Authorization'));
  await fn(req, res, email);
};

/**
 * When a producer wants to register for the first time.
 */
carbonMarketRouter.post('/admin/auth/register', catchErrors(async (req, res) => {
  const { email, firm, password } = req.body;
  await utils.registerProducer(email, password);
  await utils.registerTokens(email);
  const token = await register(email, password, firm);
  return res.json({ token });
}));

/**
 * Login of a producer into the system.
 */
carbonMarketRouter.post('/admin/auth/login', catchErrors(async (req, res) => {
  const { email, password } = req.body;
  const token = await utils.loginProducer(email, password);
  return res.json({ token });
}));

carbonMarketRouter.post('/offer/create',
  catchErrors(authed(async (req, res, email) => {
    const { amount, tokens } = req.body;
    console.log(amount);
    console.log(tokens);
    await utils.addOffer(email, amount, tokens);
    return res.json({ message: 'success' });
  })));

/**
 * POST to get tokens to purchase from the market
 */
carbonMarketRouter.post('/offer/buy',
  catchErrors(authed(async (req, res, email) => {
    const { purchased, id } = req.body;
    console.log(purchased);
    console.log(id);
    const balance = await utils.acceptOffer(email, purchased, id);
    return res.json({ balance });
  })));

/**
 * For retrieving the balance of a particular user
 */
carbonMarketRouter.get('/token/balance',
  catchErrors(authed(async (_, res, email) => {
    const balance = await utils.retrieveBalance(email);
    return res.json({ balance });
  })));

carbonMarketRouter.get('/offers/list',
  catchErrors(authed(async (req, res) => {
    const token = req.query.token ? req.query.token : '';
    const size = req.query.amount ? req.query.amount : 10;
    const queryResult = await utils.getOffers(token, size);
    return res.json(queryResult);
  })));

export default carbonMarketRouter;
