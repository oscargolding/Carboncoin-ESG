import React from 'react';
import Purchase from '../components/PurchaseOfferPage/Purchase';
import API from '../utils/API';

/**
 * Purchase the given number of tokens
 * @returns the page to purchase offers from
 */
const PurchaseOfferPage = () => {
  return (
    <>
      <Purchase acceptOffer={API.acceptOffer} />
    </>
  );
};

export default PurchaseOfferPage;
