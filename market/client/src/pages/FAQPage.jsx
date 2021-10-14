import React from 'react';
import { Typography, } from '@material-ui/core';

/**
 * Details on the website.
 * @returns the FAQPage for the site
 */
const FAQPage = () => {
  return <>
    <h1> Frequently Asked Questions </h1>
    <h3> What is reputation? </h3>
    <Typography varaint="body2" component="p">
      Reputation is a score given to users on the carbon market
      as a way to measure the market reputation of a carbon
      producer. Carboncoin is linked to the Environmental, Social
      and Governance performance of firms who trade on the platform.
      ESG certificates reported on the blockchain are automatically
      monitored by Carboncoin and added to the reputation score.
    </Typography>
    <h3> How is carbon production measured? </h3>
    <Typography varaint="body2" component="p">
      The standard unit for carbon production on the blockchain is
      Carbon dioxide equivalent (CO2e). The unit is expressed in
      terms parts per million by volume, ppmv.
    </Typography>
    <h3> What currency is used to buy Carboncoin? </h3>
    <Typography varaint="body2" component="p">
      Australia Dollars (AUD). Carboncoin can be extended in the future
      to support cryptographic coins such as Ether and Bitcoin.
    </Typography>
  </>;
};

export default FAQPage;
