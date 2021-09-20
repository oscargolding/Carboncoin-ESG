/**
 * Utils for server connecting to Fabric
 */
import fs from 'fs';
import path from 'path';
import FabricCAServices from 'fabric-ca-client';
import { Wallets, Gateway } from 'fabric-network';
import { fileURLToPath } from 'url';
import { v4 as uuidv4 } from 'uuid';
import { login } from './auth.js';

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

// Credential details for an administrator on the organisation
const adminUserId = 'admin';
const adminUserPasswd = 'adminpw';

// Channel / chaincode details
const channel = 'mychannel';
let chaincode = 'basic';

// The wallet - note the wallet can often be out of sync and might need deleting
let wallet;
let caClient;
let ccpConfig;

// The msp credentials for organisation one
let mspOrg1 = 'Org1MSP';

// Utilities available when accessing the blockchain.
const utils = {};

/**
 * Populate the chaincode with a new name
 * @param {chaincodeName} chaincodeName setting the chaincode
 */
utils.setChaincode = (chaincodeName) => {
  chaincode = chaincodeName;
};

/**
 * Get the in memory object representing the network configuration
 */
const buildNetworkConfig = (orgName) => {
  const ccpPath = path.resolve(dirname, '..', '..', '..', '..', 'test-network',
    'organizations', 'peerOrganizations',
    `${orgName}.example.com`, `connection-${orgName}.json`);
  const fileExists = fs.existsSync(ccpPath);
  if (!fileExists) {
    throw new Error(`No such file or directory: ${ccpPath}`);
  }
  const contents = fs.readFileSync(ccpPath, 'utf-8');

  // create the json object
  const ccp = JSON.parse(contents);

  console.log(`Loaded the network configuration at ${ccpPath}`);
  return ccp;
};

/**
 * Build the Certificate Authority client.
 * @param {ccp} ccp network configuration path
 * @param {caHostName} caHostName the certificate authorith host name
 * @returns
 */
const buildCaClient = (ccp, caHostName) => {
  const caInfo = ccp.certificateAuthorities[caHostName]; // CA host details
  const caTLSCACerts = caInfo.tlsCACerts.pem;
  const usingCaClient = new FabricCAServices(caInfo.url,
    { trustedRoots: caTLSCACerts, verify: false }, caInfo.caName);

  console.log(`Built a CA Client named ${caInfo.caName}`);
  return usingCaClient;
};

const buildWallet = async (walletPath) => {
  // Create a new wallet to be used for managing identities
  let usingWallet;
  if (walletPath) {
    usingWallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Built a file system wallet at ${walletPath}`);
  } else {
    usingWallet = await Wallets.newInMemoryWallet();
    console.log('Built a new memory wallet');
  }
  return usingWallet;
};

const enrolAdmin = async (usableCaClient, usableWallet, orgMspId, adValue) => {
  try {
    // Is the admin user present
    const adminIdentity = await usableWallet.get(adminUserId);
    if (adminIdentity) {
      console.log('An identity for the user already exists in the wallet');
      console.log('Hint: if errors are persisting delete the file system wallet');
      return;
    }
    const enrollment = await usableCaClient.enroll({
      enrollmentID: adminUserId,
      enrollmentSecret: adminUserPasswd,
      attrs: [
        {
          name: 'usertype', // application role
          value: adValue,
          ecert: true,
        }],
    });
    const x509Identity = {
      credentials: {
        certificate: enrollment.certificate,
        privateKey: enrollment.key.toBytes(),
      },
      mspId: orgMspId,
      type: 'X.509',
    };
    await usableWallet.put(adminUserId, x509Identity);
    console.log('Successfully enrolled admin user and imported it into the wallet');
  } catch (error) {
    console.log(`Failed to enroll admin user: ${error}`);
  }
};

/**
 * @param {orgName} orgName organisation ame to use when connecting
 * @param {mspId} mspId the mspId to use when connecting
 * @param {adValue} adValue the admin value to use when connecting
 * Launch the middleware layer for creating the Hydrogen market.
 */
utils.connectGateway = async (orgName, mspId, adValue) => {
  console.log(`>>> Starting a connection to the gateway with msp ${mspId}`);
  mspOrg1 = mspId;
  // Generate the network configuration from the file system
  ccpConfig = buildNetworkConfig(orgName);
  // Build the certificate authority for the application
  caClient = buildCaClient(ccpConfig, `ca.${orgName}.example.com`);
  // Create a wallet
  wallet = await buildWallet('wallet');
  console.log('Built a certificate authority and a wallet');
  // Now enroll the admin
  await enrolAdmin(caClient, wallet, mspOrg1, adValue);
  console.log('>>> The admin is now enrolled');
};

/**
 * Register a user into the system.
 * @param {username} username of the hydrogen producer.
 * @param {password} password of the hydrogen producer.
 * @throws {error} on any failure.
 */
utils.registerProducer = async (username, password) => {
  console.log(`>>> Attempting register for ${username} and ${password}`);
  // Check if the user exists in the wallet
  const userIdentity = await wallet.get(username);
  if (userIdentity) {
    console.log(`An identity for the user ${username} already exists`);
    return;
  }

  // Retrieve the admin identity
  const adminIdentity = await wallet.get(adminUserId);
  if (!adminIdentity) {
    console.log('Identity for the admin does not exist');
    return;
  }

  // create the user object dor the CA
  const provider = wallet.getProviderRegistry().getProvider(adminIdentity.type);
  const adminUser = await provider.getUserContext(adminIdentity, adminUserId);
  const attributes = [{
    name: 'usertype',
    value: 'producer',
    ecert: true,
  }];

  // register and enroll the user
  const secret = await caClient.register({
    enrollmentID: username,
    enrollmentSecret: password,
    attrs: attributes,
  }, adminUser);
  const enrollment = await caClient.enroll({
    enrollmentID: username,
    enrollmentSecret: secret,
    attrs: attributes,
  });
  const x509Identity = {
    credentials: {
      certificate: enrollment.certificate,
      privateKey: enrollment.key.toBytes(),
    },
    mspId: mspOrg1,
    type: 'X.509',
  };
  await wallet.put(username, x509Identity);
  console.log(`>>> Successfully enrolled ${username} and put in the wallet`);
};

/**
 * Login a producer into the system.
 * @param {username} username of the producer to login.
 * @param {*} password of the producer to login.
 */
utils.loginProducer = async (username, password) => {
  console.log(`>>> Attempting login for ${username} and ${password}`);
  const userIdentity = await wallet.get(username);
  if (!userIdentity) {
    throw new Error('Credentials do not exist - system error');
  }
  // Perform the login off-chain
  const token = await login(username, password);
  console.log(`>>> User is enrolled ${token}`);
  return token;
};

/**
 * Get the associated credentials.
 * @param {credentilas} credentials given credentilas
 * @returns the credentials
 */
utils.getContract = async (credentials) => {
  // Create the gateway and connect to it
  const gateway = new Gateway();
  await gateway.connect(ccpConfig, {
    wallet,
    identity: credentials,
    discovery: { enabled: true, asLocalhost: true },
  });
  // Here have connected to the gateway
  const network = await gateway.getNetwork(channel);

  // Get the smart contract on the channel
  const contract = network.getContract(chaincode);
  return { contract, gateway };
};

/**
 * Register tokens for a given user on the gateway.
 * @param {userId} userId the userId to register
 */
utils.registerTokens = async (userId) => {
  console.log('>>> Starting the process of providing tokens');
  const { contract, gateway } = await utils.getContract(adminUserId);

  // Add the producer and give them a number of tokens
  await contract.submitTransaction('AddProducer', userId);

  // Leave the application
  gateway.disconnect();
  console.log('>>> Successfully provided tokens to the user');
};

/**
 * Retrieving the balance of a particular user in the system.
 * @param {userId} userId the id to retrieve the balance of.
 */
utils.retrieveBalance = async (userId) => {
  console.log('>>> Retrieving the balance for a user');
  const { contract, gateway } = await utils.getContract(adminUserId);

  // Submit the transaction
  const balance = await contract.submitTransaction('GetBalance', userId);

  // Return and disconnect from the gateway
  gateway.disconnect();
  console.log(`>>> Retrieved ${balance} for user ${userId}`);
  // Note - have it as a string to get proper data returned (is in a buffer)
  return balance.toString();
};

utils.addOffer = async (userId, amount, tokens) => {
  console.log('>>> Adding an offer for user');
  console.log(`${userId}, ${amount}, ${tokens}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the offer
  const offeruuid = uuidv4();
  await contract.submitTransaction('AddOffer', userId, amount, tokens, offeruuid);

  // Return and disconnect
  gateway.disconnect();
  console.log(`>>> Added the sale offer for ${tokens},${amount},${userId}`);
};

utils.acceptOffer = async (userId, amount, offerId) => {
  console.log('>>> Purchasing an offer for a user');
  console.log(`${userId}, ${amount}, ${offerId}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the offer
  const balance = await contract.submitTransaction('PurchaseOfferTokens',
    offerId, amount);

  // Return and disconnect
  gateway.disconnect();
  console.log(`>>> Purchased tokens for ${amount} with offer id ${offerId}`);
  return balance.toString();
};

/**
 * Get the offers for a particular user
 * @param {token} paginationToken the pagination token
 * @param {number} number the number of results
 */
utils.getOffers = async (paginationToken, number, field, direction, email,
  username) => {
  console.log('>>> Getting offers for a user');
  console.log(`${paginationToken} ${number}`);
  console.log(`${field} and dir -> ${direction}`);
  console.log(`Username present -> ${username}`);
  const { contract, gateway } = await utils.getContract(email);

  // Submit the offer
  const result = await contract.submitTransaction('GetOffers',
    number, paginationToken, field, direction, username);

  const jsonResult = JSON.parse(result);
  // Return and disconnect
  gateway.disconnect();
  console.log(`>>> Retrieved listing for ${jsonResult.fetchedRecordsCount} number`);
  return jsonResult;
};

utils.getProduction = async (userId, paginationToken, number) => {
  console.log('>>> Getting production for a user');
  console.log(`${paginationToken} ${number}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the request to get production for a user
  const result = await contract.submitTransaction('GetProduction', number,
    paginationToken);

  const jsonResult = JSON.parse(result);
  // Return and disconnect
  gateway.disconnect();
  console.log(`>>> Retrieved listing for ${jsonResult.fetchedRecordsCount}`);
  return jsonResult;
};

utils.payProduction = async (userId, prodId) => {
  console.log('>>> Paying for user production');
  console.log(`${userId} ${prodId}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the transaction and do the payment
  const balance = await contract.submitTransaction('PayForProduction', prodId);
  // Return and disconnect
  gateway.disconnect();
  console.log('>>> Paid for production');
  return balance.toString();
};

utils.getDirectPrice = async (userId) => {
  console.log('>>> Determining the direct purchase price');
  console.log(`${userId}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the transaction and then wait
  const directPrice = await contract.submitTransaction('GetDirectPrice');
  // Return and disconnect
  gateway.disconnect();
  console.log(`Got the price >>> ${directPrice}`);
  return directPrice.toString();
};

/**
 * Redeeming the chip for a user - chips represent a direct offer on blockchain
 * @param {*} userId user wanting to redeem the chip
 * @returns the balance available to the user
 */
utils.redeemChip = async (userId, amount) => {
  console.log('>>> Allowing a user to redeem a chip');
  console.log(`${userId}`);
  console.log(`${amount}`);
  const { contract, gateway } = await utils.getContract(userId);

  // Submit the transaction and then wait
  const balance = await contract.submitTransaction('RedeemChip', amount);
  // Return and disconnect
  gateway.disconnect();
  console.log(`Got the following balance >>> ${balance}`);
  return balance.toString();
};

/**
 * Helper to query target offers within a given range
 * @param {*} reputationSort sorting on reputation or not
 * @param {*} target the target value to use
 * @returns the balance available to the user
 */
utils.targetOffers = async (userId, reputationSort, target) => {
  console.log('>>> Allowing a user to do a target offer search');
  console.log(`Sort present -> ${reputationSort}`);
  console.log(`Target present -> ${target}`);

  const { contract, gateway } = await utils.getContract(userId);

  // Submit the transaction
  const result = await contract.submitTransaction('GetBudgetOffer',
    reputationSort, target);
  const jsonResult = JSON.parse(result);
  // Return and disconnect
  gateway.disconnect();
  return jsonResult;
};

utils.deleteOffer = async (userId, offerId) => {
  console.log('>>> Allowing a user to delete an offer');
  console.log(`The username ${userId} and offer Id ${offerId}`);

  const { contract, gateway } = await utils.getContract(userId);

  // Submit the transaction
  await contract.submitTransaction('MakeOfferStale', offerId);

  gateway.disconnect();
  console.log('>>> Deleted the offer on chain');
};

export default utils;
