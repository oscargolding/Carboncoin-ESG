import chalk from 'chalk';
import { v4 as uuid4 } from 'uuid';
import moment from 'moment';
import utils from '../utils.js';

const productionList = [
  {
    firm: 'oscarIndustry',
    carbonProduction: '-5',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
  {
    firm: 'oscarIndustry',
    carbonProduction: '-305',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
  {
    firm: 'oscarIndustry',
    carbonProduction: '-310',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
  {
    firm: 'BHPPetrol',
    carbonProduction: '-150',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
  {
    firm: 'BHPPetrol',
    carbonProduction: '10',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
];

const simpleList = [
  {
    firm: 'modernenergy@hydrogen.com',
    carbonProduction: '-50',
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  },
];

const blockchainCall = async (production, contract) => {
  await contract.submitTransaction('ReportProduction', production.firm,
    production.carbonProduction, production.date, production.id);
};
/**
 * Represents the second organisation and connecting to the hyplerdger network
 */
const main = async () => {
  try {
    // Represent a certifier joining the network
    await utils.connectGateway('org2', 'Org2MSP', 'certifier');
    console.log(chalk.green('Enrolled org2 with a certifier'));
    utils.setChaincode('EnergyCertifier');
    const production = process.argv.length > 2 ? simpleList : productionList;
    const { contract, gateway } = await utils.getContract('admin');
    try {
      console.log('>>> Making calls to blockchain');
      for (let i = 0; i < production.length; i += 1) {
        const produce = production[i];
        console.log(`Making call for ${produce.firm}`);
        // eslint-disable-next-line no-await-in-loop
        await blockchainCall(produce, contract);
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
