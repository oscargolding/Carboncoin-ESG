import React from 'react';
import Offer from '../components/CreateOffer/Offer';
import API from '../utils/API';

/**
 * Create the offer of a token sale.
 * @returns the Offer
 */
const CreateOfferPage = () => {
  return (
    <>
      <Offer createOffer={API.createOffer} />
    </>
  );
};

export default CreateOfferPage;
