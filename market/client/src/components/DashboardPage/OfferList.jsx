import React, { useState, useEffect, } from 'react';
import OfferCard from './OfferCard';
import { CentralLoading, } from './styles/DashboardStyles';
import { CircularProgress, } from '@material-ui/core';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';
import API from '../../utils/API';

/**
 * Represent an offer list being provided.
 * @returns the offerlist
 */
const OfferList = () => {
  const [loading, setLoading] = useState(true);
  const [gameList, setGameList] = useState([]);
  const [errorMessage, setErrorMessage] = useState('');
  const { authToken: [authToken], } = storeContext();
  useEffect(() => {
    const getOffers = async () => {
      try {
        const response = await API.getOffers(authToken);
        const apiGameList = response.records.map((offer, i) => {
          return (
            <OfferCard
              key={i}
              producer={offer.producer}
              price={offer.amount}
              quantity={offer.tokens}
              active={offer.active}
            />
          );
        });
        setGameList(apiGameList);
      } catch (err) {
        setErrorMessage(err.message);
      }
      setLoading(false);
    };
    getOffers();
  }, []);
  return (
    <>
      {loading
        ? <CentralLoading> <CircularProgress /> </CentralLoading>
        : gameList}
      {errorMessage !== ''
        ? <Alert severity='error'>{errorMessage}</Alert>
        : <></>}
    </>
  );
};

export default OfferList;
