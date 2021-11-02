import React, { useState, useRef, useCallback, useEffect, } from 'react';
import OfferCard from './OfferCard';
import { CentralLoading, } from './styles/DashboardStyles';
import { CircularProgress, } from '@material-ui/core';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';
import useOfferSearch from './useOfferSearch';
import PropTypes from 'prop-types';
import API from '../../utils/API';

/**
 * Represent an offer list being provided.
 * @returns the offerlist
 */
const OfferList = (props) => {
  const { authToken: [authToken], } = storeContext();
  const { sortTerm, direction, username, } = props;
  const [pageToken, setPageToken] = useState('');
  useEffect(() => {
    setPageToken('');
  }, [sortTerm, direction]);
  const {
    loading,
    error,
    offers,
    hasMore,
    paginationToken,
    setOffers,
  } = useOfferSearch(pageToken, authToken, API.getOffers, sortTerm, direction,
    username);
  const observer = useRef();
  const deleteOffer = offerId => {
    const copyArray = JSON.parse(JSON.stringify(offers));
    const index = copyArray.findIndex((offer) => offer.offerId === offerId);
    copyArray.splice(index, 1);
    setOffers(copyArray);
  };
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
              reputation={offer.reputation}
              owned={offer.owned}
              deleteOfferFn={deleteOffer}
              usingRef={lastElementRef}
              environment={offer.environment}
              social={offer.social}
              governance={offer.governance}
              offers={[]}
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
              reputation={offer.reputation}
              owned={offer.owned}
              offerid={offer.offerId}
              deleteOfferFn={deleteOffer}
              environment={offer.environment}
              social={offer.social}
              governance={offer.governance}
              offers={[]}
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

OfferList.propTypes = {
  sortTerm: PropTypes.string.isRequired,
  direction: PropTypes.bool.isRequired,
  username: PropTypes.string.isRequired,
};
