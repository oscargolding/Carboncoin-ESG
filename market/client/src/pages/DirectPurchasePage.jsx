import React from 'react';
import PurchaseDirect from '../components/DirectPurchase/PurchaseDirect';
import API from '../utils/API';

/**
 * A page to directly purchase from.
 * @returns the DirectPurchasePage
 */
const DirectPurchasePage = () => {
  return (
    <>
      <PurchaseDirect
        getOffer={API.getDirectPrice}
        redeemOffer={API.acceptDirectPrice}
      />
    </>
  );
};

export default DirectPurchasePage;
