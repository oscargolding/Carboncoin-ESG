import chalk from 'chalk';
import { v4 as uuid4 } from 'uuid';
import moment from 'moment';
import utils from '../utils.js';

const productionList = [
  {
    name: 'Female Employee Rate',
    value: 80,
    firm: 'oscarIndustry',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
    statistic: '80%',
  },
  {
    name: 'Water Score in Hydrogen Production',
    value: 200,
    firm: 'oscarIndustry',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
    statistic: '200 Points',
  },
  {
    name: 'Board Diversity and Structure',
    value: 75,
    firm: 'oscarIndustry',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
    statistic: '75%',
  },
  {
    name: 'Female Employee Rate',
    value: 55,
    firm: 'BHPPetrol',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
    statistic: '55%',
  },
];

const createCertificate = {
  name: 'Female Employee Rate',
  category: 'Social',
  min: 0,
  max: 100,
  Multiplier: 2,
};

const waterScore = {
  name: 'Water Score in Hydrogen Production',
  category: 'Environmental',
  min: 0,
  max: 20000,
  Multiplier: 1,
};

const governance = {
  name: 'Board Diversity and Structure',
  category: 'Governance',
  min: 0,
  max: 100,
  Multiplier: 2,
};

const blockchainCall = async (production, contract, certificate) => {
  await contract.submitTransaction(certificate, production.name,
    production.value, production.firm, production.date, production.id,
    production.statistic);
};
/**
 * Represents the second organisation and connecting to the hyplerdger network
 */
const main = async () => {
  try {
    // Represent a certifier joining the network
    await utils.connectGateway('org2', 'Org2MSP', 'register');
    console.log(chalk.green('Enrolled org2 with a register'));
    utils.setChaincode('Register');
    const production = productionList;
    const { contract, gateway } = await utils.getContract('admin');
    try {
      console.log(chalk.green('>>> Making calls to blockchain for register'));
      await contract.submitTransaction('CreateCertificate',
        createCertificate.name, createCertificate.category,
        createCertificate.min, createCertificate.max,
        createCertificate.Multiplier);
      await contract.submitTransaction('CreateCertificate',
        waterScore.name, waterScore.category,
        waterScore.min, waterScore.max,
        waterScore.Multiplier);
      await contract.submitTransaction('CreateCertificate',
        governance.name, governance.category,
        governance.min, governance.max,
        governance.Multiplier);
      for (let i = 0; i < production.length; i += 1) {
        const produce = production[i];
        console.log(`Making call for ${produce.firm} (Register)`);
        // eslint-disable-next-line no-await-in-loop
        await blockchainCall(produce, contract, 'RegisterUserCertificate');
        console.log(chalk.green('Made call for firm'));
      }
    } catch (err) {
      console.log(chalk.red(`*** error -> ${err.message}`));
    }
    gateway.disconnect();
    console.log(chalk.green('Finished calling the chaincode'));
    process.exit(1);
  } catch (err) {
    console.log(
      chalk.red(`*** Failed to run the main application and report: ${err.message}`),
    );
  }
};
// Perform the main setup and connect
main();
