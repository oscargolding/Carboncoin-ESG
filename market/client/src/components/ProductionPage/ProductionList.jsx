import React, { useState, useRef, useCallback, } from 'react';
import API from '../../utils/API';
import { storeContext, } from '../../utils/store';
import { CentralLoading, } from '../DashboardPage/styles/DashboardStyles';
import { CircularProgress, Button, } from '@material-ui/core';
import { Alert, } from '@material-ui/lab';
import useOfferSearch from '../DashboardPage/useOfferSearch';
import ProductionCard from './ProductionCard';
import RepChart from './RepChart';
import { ReputationDiv, } from './styles/ProductionStyles';

/**
 * The production list
 * @returns the production associated with an individual on the carbon market.
 */
const ProductionList = () => {
  const { authToken: [authToken], } = storeContext();
  const [pageToken, setPageToken] = useState('');
  const [button, setButton] = useState(false);
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
        ? <>
          <ReputationDiv reputation={response.reputation}>
            <h2> Total Reputation {response.reputation}</h2>
          </ReputationDiv>
          <Button
            variant='contained'
            color='primary'
            onClick={() => {
              setButton(!button);
            }}
          >
            ESG Breakdown
          </Button>
          {button
            ? <><RepChart
              environment={response.environment}
              social={response.social}
              governance={response.governance}
            /></>
            : <></>}
        </>
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
              statistic={production.statistic}
              multiplier={production.multiplier}
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
              statistic={production.statistic}
              multiplier={production.multiplier}
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
