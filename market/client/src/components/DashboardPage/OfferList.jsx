import React, { useState, useRef, useCallback, } from 'react';
import OfferCard from './OfferCard';
import { CentralLoading, } from './styles/DashboardStyles';
import { CircularProgress, } from '@material-ui/core';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';
import useOfferSearch from './useOfferSearch';

/**
 * Represent an offer list being provided.
 * @returns the offerlist
 */
const OfferList = () => {
  const { authToken: [authToken], } = storeContext();
  const [pageToken, setPageToken] = useState('');
  const {
    loading,
    error,
    offers,
    hasMore,
    paginationToken,
  } = useOfferSearch(pageToken, authToken);
  const observer = useRef();
  const lastElementRef = useCallback(node => {
    if (loading === true) {
      return;
    }
    if (observer.current) {
      observer.current.disconnect();
    }
    observer.current = new IntersectionObserver(entries => {
      if (entries[0].isIntersecting && hasMore) {
        setPageToken(paginationToken);
      }
    });
    if (node) {
      observer.current.observe(node);
    }
  }, [loading, hasMore]);
  return (
    <>
      {offers.map((offer, i) => {
        if (offers.length === i + 1) {
          return (
            <OfferCard
              key={i}
              producer={offer.producer}
              price={offer.amount}
              quantity={offer.tokens}
              active={offer.active}
              offerid={offer.offerId}
              usingRef={lastElementRef}
            />
          );
        } else {
          return (
            <OfferCard
              key={i}
              producer={offer.producer}
              price={offer.amount}
              quantity={offer.tokens}
              active={offer.active}
              offerid={offer.offerId}
            />
          );
        }
      })}
      {loading
        ? <CentralLoading> <CircularProgress /> </CentralLoading>
        : <></>}
      {error !== ''
        ? <Alert severity='error'>{error}</Alert>
        : <></>}
    </>
  );
};

export default OfferList;
