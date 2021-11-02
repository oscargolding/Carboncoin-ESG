import React, { useState, } from 'react';
import { useLocation, } from 'react-router';
import OfferFinder from '../components/OfferFinderPage/OfferFinder';

/**
 * Page to represent the finding of offers on the blockchain.
 * @returns the OfferFinderPage
 */
const OfferFinderPage = () => {
  const location = useLocation();
  const [offers, setOffers] = useState(location.state.offers);
  return (
    <>
      <h1> Carboncoin Offer Finder </h1>
      <OfferFinder setOffers={setOffers} offers={offers} />
    </>
  );
};

export default OfferFinderPage;
