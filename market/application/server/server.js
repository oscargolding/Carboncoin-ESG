/**
 * Blockchain carbon market - Oscar Golding
 */
import express from 'express';
import cors from 'cors';
import utils from './utils.js';
import carbonMarketRouter from './carbonMarket.js';

const app = express();

// Create the express functions
app.use(express.json());
app.use((_, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization');
  res.header('Access-Control-Allow-Methods', 'PUT, POST, GET, DELETE, OPTIONS');
  next();
});
app.use(cors());

app.use('/api', carbonMarketRouter);

/**
 * Standard ping from the server
 */
app.get('/ping', (_, res) => {
  res.send('Response from the blockchain carbon market!');
});

/**
 * Entry into the application - launch middleware to create the application
 */
const main = async () => {
  const port = process.env.PORT || 3000;
  try {
    await utils.connectGateway('org1', 'Org1MSP', 'admin');
  } catch (error) {
    console.log(`Error in connecting to Fabric network ${error}`);
  }
  app.listen(port, (error) => {
    if (error) {
      console.log(`Error ${error}`);
    }
    console.log(`Carbon market server listening on port: ${port}`);
  });
};
main();
