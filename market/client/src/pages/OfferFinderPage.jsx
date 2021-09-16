import React from 'react';
import OfferFinder from '../components/OfferFinderPage/OfferFinder';

/**
 * Page to represent the finding of offers on the blockchain.
 * @returns the OfferFinderPage
 */
const OfferFinderPage = () => {
  return (
    <>
      <h1> Carboncoin Offer Finder </h1>
      <OfferFinder />
    </>
  );
};

export default OfferFinderPage;
