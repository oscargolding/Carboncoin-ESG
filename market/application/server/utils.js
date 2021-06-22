/**
 * Utils for server connecting to Fabric
 */
import fs from 'fs';
import path from 'path';
import FabricCAServices from 'fabric-ca-client';
import { Wallets, Gateway } from 'fabric-network';
import { fileURLToPath } from 'url';
import { login } from './auth.js';

const filename = fileURLToPath(import.meta.url);
const dirname = path.dirname(filename);

// Credential details for an administrator on the organisation
const adminUserId = 'admin';
const adminUserPasswd = 'adminpw';

// Channel / chaincode details
const channel = 'mychannel';
const chaincode = 'basic';

// The wallet - note the wallet can often be out of sync and might need deleting
let wallet;
let caClient;
let ccpConfig;

// The msp credentials for organisation one
const mspOrg1 = 'Org1MSP';

// Utilities available when accessing the blockchain.
const utils = {};

/**
 * Get the in memory object representing the network configuration
 */
const buildNetworkConfig = () => {
  const ccpPath = path.resolve(dirname, '..', '..', '..', '..', 'test-network',
    'organizations', 'peerOrganizations',
    'org1.example.com', 'connection-org1.json');
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

const enrolAdmin = async (usableCaClient, usableWallet, orgMspId) => {
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
          value: 'admin',
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
 * Launch the middleware layer for creating the Hydrogen market.
 */
utils.connectGateway = async () => {
  console.log('>>> Starting a connection to the gateway');
  // Generate the network configuration from the file system
  ccpConfig = buildNetworkConfig();
  // Build the certificate authority for the application
  caClient = buildCaClient(ccpConfig, 'ca.org1.example.com');
  // Create a wallet
  wallet = await buildWallet('wallet');
  console.log('Built a certificate authority and a wallet');
  // Now enroll the admin
  await enrolAdmin(caClient, wallet, mspOrg1);
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
const getContract = async (credentials) => {
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
  const { contract, gateway } = await getContract(adminUserId);

  // Add the producer and give them a number of tokens
  await contract.submitTransaction('AddProducer', userId, '100');

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
  const { contract, gateway } = await getContract(adminUserId);

  // Submit the transaction
  const balance = await contract.submitTransaction('GetBalance', userId);

  // Return and disconnect from the gateway
  gateway.disconnect();
  console.log(`>>> Retrieved ${balance} for user ${userId}`);
  // Note - have it as a string to get proper data returned (is in a buffer)
  return balance.toString();
};

export default utils;
