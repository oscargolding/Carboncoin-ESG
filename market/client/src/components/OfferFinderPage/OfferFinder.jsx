import { Button, Checkbox, CircularProgress, FormControlLabel, FormGroup, InputAdornment, TextField, } from '@material-ui/core';
import { Alert, } from '@material-ui/lab';
import React, { useState, useEffect, } from 'react';
import API from '../../utils/API';
import { storeContext, } from '../../utils/store';
import OfferCard from '../DashboardPage/OfferCard';
import { CentralLoading, } from '../DashboardPage/styles/DashboardStyles';
import { OfferFinderForm, } from './styles/OfferFinderStyles';
import { Link, } from 'react-router-dom';
import PropTypes from 'prop-types';

/**
 * Find offers on the open market
 */
const OfferFinder = (props) => {
  const { setOffers, offers, } = props;
  const [quantity, setQuantity] = useState(0);
  const [carbonReputation, setCarbonReputation] = useState(0);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [empty, setEmpty] = useState(false);
  const { authToken: [authToken], } = storeContext();
  useEffect(() => { }, []);
  const handleSubmit = async () => {
    setLoading(true);
    setEmpty(false);
    setOffers([]);
    setError('');
    try {
      const response = await API.getFoundOffers(authToken, quantity,
        carbonReputation);
      setLoading(false);
      setOffers(response.records);
      if (response.records.length === 0) {
        setEmpty(true);
      }
    } catch (err) {
      setLoading(false);
      setError(err.message);
    }
  };
  return (
    <>
      <p>
        Find Carboncoin offers satisfying a desired purchasing quantity.
      </p>
      <OfferFinderForm>
        <TextField
          label='Quantity of Carboncoin to Purchase'
          type='number'
          id='quantity-coin'
          value={quantity}
          onChange={(event) => {
            if (event.target.value < 0) {
              setQuantity(0);
            } else {
              setQuantity(event.target.value);
            }
          }}
          InputProps={{
            startAdornment: <InputAdornment position='start'>
              Carboncoin
            </InputAdornment>,
          }}
        >
        </TextField>
        <FormGroup row>
          <FormControlLabel
            control={<Checkbox
              checked={!carbonReputation}
              onChange={() => { setCarbonReputation(0); }}
            />}
            label='Lowest Price'
          />
          <FormControlLabel
            control={<Checkbox
              checked={carbonReputation}
              onChange={() => { setCarbonReputation(1); }}
            />}
            label='Best Carbon Reputation'
          />
        </FormGroup>
        <Button
          variant='contained'
          color='primary'
          onClick={() => { handleSubmit(); }}
        >
          Search for Offers
        </Button>
      </OfferFinderForm>
      {loading
        ? <CentralLoading> <CircularProgress /> </CentralLoading>
        : <></>}
      {error !== '' ? <Alert severity='error'>{error}</Alert> : <></>}
      {empty
        ? <h3> Sorry, no offers in the market satisfy the number of Carboncoin
          requested. Consider <Link to='/direct/purchase'>
            buying Carboncoin directly
          </Link>.
        </h3>
        : <></>}
      {offers.map((offer, i) => {
        return (
          <OfferCard
            key={i}
            producer={offer.producer}
            price={offer.amount}
            quantity={offer.tokens}
            active={offer.active}
            reputation={offer.reputation}
            offerid={offer.offerId}
            environment={offer.environment}
            social={offer.social}
            governance={offer.governance}
            offers={offers}
          />
        );
      })}
    </>
  );
};

export default OfferFinder;

OfferFinder.propTypes = {
  setOffers: PropTypes.func.isRequired,
  offers: PropTypes.any.isRequired,
};
