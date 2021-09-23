import React, { useState, useRef, useCallback, } from 'react';
import API from '../../utils/API';
import { storeContext, } from '../../utils/store';
import { CentralLoading, } from '../DashboardPage/styles/DashboardStyles';
import { CircularProgress, } from '@material-ui/core';
import { Alert, } from '@material-ui/lab';
import useOfferSearch from '../DashboardPage/useOfferSearch';
import ProductionCard from './ProductionCard';

/**
 * The production list
 * @returns the production associated with an individual on the carbon market.
 */
const ProductionList = () => {
  const { authToken: [authToken], } = storeContext();
  const [pageToken, setPageToken] = useState('');
  const {
    loading,
    error,
    offers,
    hasMore,
    paginationToken,
    response,
  } = useOfferSearch(pageToken, authToken, API.getProduction);
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
      {'records' in response
        ? <h2> Total Reputation {response.reputation}</h2>
        : <></>}
      {offers.map((production, i) => {
        if (offers.length === i + 1) {
          return (
            <ProductionCard
              key={i}
              produced={production.produced}
              date={production.date}
              paid={production.paid}
              ethical={production.ethical}
              usingRef={lastElementRef}
              id={production.productionID}
              category={production.category}
              description={production.description}
            />
          );
        } else {
          return (
            <ProductionCard
              key={i}
              produced={production.produced}
              date={production.date}
              paid={production.paid}
              ethical={production.ethical}
              id={production.productionID}
              category={production.category}
              description={production.description}
            />
          );
        }
      })}
      {loading
        ? <CentralLoading> <CircularProgress /> </CentralLoading>
        : <></>}
      {error !== ''
        ? <Alert severity='error'>{error}</Alert>
        : <></>
      }
    </>
  );
};

export default ProductionList;
