/**
 * Handling the auth logic of the server.
 */
import fs from 'fs';
import jwt from 'jsonwebtoken';
import AsyncLock from 'async-lock';
import InputError from './errors/inputError.js';
import AccessError from './errors/accessError.js';

// Locking the database
const lock = new AsyncLock();

// The token secrets
const JWT_SECRET = 'llamallamaduck';
const DATABASE_FILE = './database.json';

/**
 * Admins are the users in the system
 */
let admins = {};

const update = (savingAdmins) => new Promise((resolve, reject) => {
  lock.acquire('saveData', () => {
    try {
      fs.writeFileSync(DATABASE_FILE, JSON.stringify({
        admins: savingAdmins,
      }, null, 2));
      resolve();
    } catch {
      reject(new Error('Writing to database failed'));
    }
  });
});

/**
 * Save the admins/users as part of the process.
 * @returns nothing
 */
export const save = () => update(admins);

// Try and read the users
try {
  const data = JSON.parse(fs.readFileSync(DATABASE_FILE));
  admins = data.admins;
} catch {
  console.log('WARNING: No database found, create a new one');
  save();
}

const userLock = (callback) => new Promise((resolve, reject) => {
  lock.acquire('userAuthLock', callback(resolve, reject));
});

// The authorisation functions.
// The use of JWT ensures significantly faster turn-back time to the user.

/**
 * Get the email from the authorization provided.
 * @param {authorization} authorization the jwt token associated with the user.
 * @returns the email.
 */
export const getEmailFromAuthorization = (authorization) => {
  try {
    console.log(`>>> Token->${authorization}`);
    const token = authorization.replace('Bearer ', '');
    console.log(`>>> Token->${token}`);
    const { email } = jwt.verify(token, JWT_SECRET);
    if (!(email in admins)) {
      throw new AccessError('Invalid Token');
    }
    return email;
  } catch {
    throw new AccessError('Invalid token');
  }
};

/**
 * Login a user into the system.
 * @param {email} email being logged
 * @param {password} password being used
 * @returns whether user is logged or not
 */
export const login = (email, password) => userLock((resolve, reject) => {
  if (email in admins) {
    if (admins[email].password === password) {
      admins[email].sessionActive = true;
      resolve(jwt.sign({ email }, JWT_SECRET, { algorithm: 'HS256' }));
    }
  }
  reject(new InputError('Invalid username or password'));
});

/**
 * Register the user off-chain.
 * @param {email} email to register.
 * @param {*} password password to use.
 * @param {*} name name to use.
 * @returns the token on no errors.
 */
export const register = (email, password, name) => userLock((resolve, reject) => {
  if (email in admins) {
    reject(new InputError('Email address already registered'));
  }
  admins[email] = {
    name,
    password,
    sessionActive: true,
  };
  console.log(admins[email]);
  const token = jwt.sign({ email }, JWT_SECRET, { algorithm: 'HS256' });
  resolve(token);
});
