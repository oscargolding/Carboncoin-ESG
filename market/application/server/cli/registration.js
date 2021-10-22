import { Command } from 'commander';
import chalk from 'chalk';
import moment from 'moment';
import { v4 as uuid4 } from 'uuid';
import utils from '../utils.js';

/**
 * For reporting production.
 * @param {*} firm to report
 * @param {*} statistic underlying
 */
const reportProduction = async (firm, statistic) => {
  const toReport = {
    firm,
    carbonProduction: statistic,
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
  };
  utils.setChaincode('EnergyCertifier');
  const { contract, gateway } = await utils.getContract('admin');
  await contract.submitTransaction('ReportProduction', toReport.firm,
    toReport.carbonProduction, toReport.date, toReport.id);
  gateway.disconnect();
  console.log(chalk.green('>>> Succesfully reported production'));
};

/**
 * Create a certificate with such values
 * @param {name} name to create
 * @param {*} category category to use
 * @param {*} min value
 * @param {*} max value
 * @param {*} weight given weight
 */
const createCertificate = async (name, category, min, max, weight) => {
  const certificate = {
    name,
    category,
    min,
    max,
    Multiplier: weight,
  };
  utils.setChaincode('Register');
  const { contract, gateway } = await utils.getContract('admin');
  console.log(chalk.green('>>> Making a call to register ESG Certificate'));
  await contract.submitTransaction('CreateCertificate',
    certificate.name, certificate.category, certificate.min,
    certificate.max, certificate.Multiplier);
  gateway.disconnect();
  console.log(chalk.green('>>> Successfully registered the ESG'
    + ' certificate with the carbon market'));
};

/**
 * Register the following information with the system related to ESG.
 * @param {*} name of ESG
 * @param {*} firm to use
 * @param {*} value to use
 * @param {*} statistic underlying
 */
const registerESGData = async (name, firm, value, statistic) => {
  const ESG = {
    name,
    value,
    firm,
    date: moment().format('MMMM Do YYYY, h:mm:ss a'),
    id: uuid4(),
    statistic,
  };
  utils.setChaincode('Register');
  const { contract, gateway } = await utils.getContract('admin');
  console.log(chalk.green('>>> Making a call to register ESG data'));
  await contract.submitTransaction('RegisterUserCertificate', ESG.name,
    ESG.value, ESG.firm, ESG.date, ESG.id, ESG.statistic);
  gateway.disconnect();
  console.log(chalk.green('>>> Successfully registered raw ESG data'));
};

/**
 * Entry into the application
 */
const main = async () => {
  await utils.connectGateway('org2', 'Org2MSP', 'certifier');
  console.log(chalk.green('Enrolled with the certificate authority!'));
  const program = new Command();
  program.option('-p, --production', 'register production')
    .option('-s, --statistic <type>', 'underlying statistic')
    .option('-f, --firm <type>', 'firm to register')
    .option('-c, --certificate', 'create certificate')
    .option('-n, --name <type>', 'certificate name')
    .option('-ca, --category <type>', 'category type')
    .option('-mi, --min <type>', 'min value')
    .option('-ma, --max <type>', 'max value')
    .option('-w, --weight <type>', 'weight given')
    .option('-r, --raw', 'esg raw data')
    .option('-v, --value <type>', 'value');
  program.parse(process.argv);
  const options = program.opts();
  if (options.production) {
    console.log('>>> Registering production for firm');
    if (!(options.statistic && options.firm)) {
      console.log(chalk.red('Specify production and firm'));
      process.exit(0);
    }
    const { firm, statistic } = options;
    await reportProduction(firm, statistic);
    process.exit(1);
  } else if (options.certificate) {
    if (!(options.name && options.category && options.min && options.max
      && options.weight)) {
      console.log(chalk.red('Wrong parameters for registering weights'));
      process.exit(0);
    }
    const {
      name, category, min, max, weight,
    } = options;
    await createCertificate(name, category, min, max, weight);
    process.exit(1);
  } else if (options.raw) {
    if (!(options.name && options.firm && options.value && options.statistic)) {
      console.log(chalk.red('Wrong parameters for registering ESG data'));
      process.exit(0);
    }
    const {
      name, firm, value, statistic,
    } = options;
    await registerESGData(name, firm, value, statistic);
    process.exit(0);
  }
};
// Perform the main setup and connect
main();
