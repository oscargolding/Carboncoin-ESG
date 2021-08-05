import { Button, InputAdornment, LinearProgress, TextField, } from '@material-ui/core';
import React, { useState, useEffect, } from 'react';
import PropTypes from 'prop-types';
import { OfferForm, } from '../CreateOffer/styles/OfferStyles';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';
import { OfferDetails, } from '../PurchaseOfferPage/styles/PurchaseStyles';
import { Link, } from 'react-router-dom';

/**
 * Ability to directly purchase from the market.
 * @param {*} props props passed into the functional component
 * @returns the direct purchase offer.
 */
const PurchaseDirect = (props) => {
  const { getOffer, redeemOffer, } = props;
  const [loading, setLoading] = useState(true);
  const [price, setPrice] = useState('');
  const [error, setError] = useState('');
  const [redeemLoad, setRedeemLoad] = useState(false);
  const [success, setSuccess] = useState(false);
  const [quantity, setQuantity] = useState('');
  const { authToken: [authToken], balance: [, setBalance], } = storeContext();
  // Run on the first time the component is called
  useEffect(() => {
    const callMarket = async () => {
      setLoading(true);
      setSuccess(false);
      try {
        const price = await getOffer(authToken);
        setPrice(price.price);
        setLoading(false);
      } catch (err) {
        setLoading(false);
        setError(err.message);
      }
    };
    callMarket();
  }, []);
  const purchaseTokens = async () => {
    setRedeemLoad(true);
    setError('');
    try {
      const response = await redeemOffer(authToken, quantity);
      setSuccess(true);
      setRedeemLoad(false);
      setBalance(response.balance);
    } catch (err) {
      setRedeemLoad(false);
      setError(err.message);
    }
  };
  return (
    <OfferForm>
      <h1> Direclty Purchase Carboncoin from the Market! </h1>
      {loading
        ? <LinearProgress />
        : <>
          <OfferDetails elevation={3}>
            <h2> Price Offered per Carboncoin: ${price} </h2>
            <h4>
              Note: price offered is based on the highest offer
              in the market with a margin added. Consider going to
              the <Link to='/dashboard'>open market</Link> for a cheaper price.
            </h4>
          </OfferDetails>
          <TextField
            label='Quantity of Carboncoin'
            type='number'
            id='standard-start-adornment'
            onChange={(event) => {
              setQuantity(event.target.value);
            }}
            InputProps={{
              startAdornment:
                <InputAdornment position='start'>Carboncoin</InputAdornment>,
            }}
          />
          <Button
            variant="contained"
            color="primary"
            size="large"
            onClick={() => purchaseTokens()}
          >
            Purchase Carboncoin!
          </Button>
          {redeemLoad ? <LinearProgress /> : <></>}
          {success
            ? <Alert severity='success'> Carboncoin directly purchased! </Alert>
            : <></>}
        </>}
      {error !== '' ? <Alert severity='error'> {error} </Alert> : <></>}
    </OfferForm >
  );
};

export default PurchaseDirect;

PurchaseDirect.propTypes = {
  getOffer: PropTypes.func.isRequired,
  redeemOffer: PropTypes.func.isRequired,
};
