import React, { useState, } from 'react';
import PropTypes from 'prop-types';
import { useHistory, useLocation, } from 'react-router';
import { OfferForm, } from '../CreateOffer/styles/OfferStyles';
import {
  InputAdornment, TextField, Button,
  LinearProgress,
} from '@material-ui/core';
import { OfferDetails, } from './styles/PurchaseStyles';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';

/**
 * A component that allows for a user to purchase from a page.
 * @returns the Purchase component that allows for purchasing
 */
const Purchase = (props) => {
  const { acceptOffer, } = props;
  const location = useLocation();
  const [quantity, setQuantity] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const {
    authToken: [authToken],
    balance: [, setBalance],
  } = storeContext();
  const offerId = location.state.offerid;
  const quantityAvailable = location.state.quantity;
  const producer = location.state.producer;
  const price = location.state.price;
  const offers = location.state.offers;
  const history = useHistory();
  const [sellableQuantity, setSellableQuantity] = useState(quantityAvailable);
  const purchaseTokens = async () => {
    setLoading(true);
    setSuccess(false);
    setError('');
    try {
      const response = await acceptOffer(authToken, offerId, quantity);
      setLoading(false);
      setSuccess(true);
      setSellableQuantity((prev) => prev - quantity);
      setBalance(response.balance);
    } catch (err) {
      setLoading(false);
      setError(err.message);
    }
  };
  return (
    <OfferForm>
      {offers.length > 0
        ? <Button
          variant="contained"
          color="primary"
          size="large"
          onClick={() => {
            const offerCopy = JSON.parse(JSON.stringify(offers));
            if (sellableQuantity === 0) {
              offerCopy.splice(offerCopy.findIndex((offer) => {
                return offer.offerId === offerId;
              }), 1);
            }
            history.push('/offerfinder', { offers: offerCopy, });
          }}
        > Go Back to Offer Finder
        </Button>
        : <></>}
      <h1> Purchase Carboncoin from Seller: {producer} </h1>
      <OfferDetails elevation={3}>
        <h2> Offer Details: </h2>
        <h3> Price Offered per Carboincoin: AUD{price}</h3>
        <h3> Quantity Available for Purchase: {sellableQuantity} </h3>
      </OfferDetails>
      <TextField
        label='Quantity of Carboncoin'
        error={quantity > sellableQuantity}
        type='number'
        helperText={`Max available to purchase: ${sellableQuantity}`}
        id='standard-start-adornment'
        onChange={(event) => {
          console.log(event.target.value);
          setQuantity(event.target.value);
        }}
        InputProps={{
          startAdornment:
            <InputAdornment position="start">Carboncoin</InputAdornment>,
        }}
      />
      <Button
        variant="contained"
        color="primary"
        size="medium"
        onClick={() => { purchaseTokens(); }}
      >
        Purchase Carboncoin!
      </Button>
      {loading ? <LinearProgress /> : <></>}
      {success
        ? <Alert severity='success'> Carboncoin was purchased!</Alert>
        : <></>
      }
      {error !== '' ? <Alert severity='error'>{error}</Alert> : <></>}
    </OfferForm>
  );
};

Purchase.propTypes = {
  acceptOffer: PropTypes.func.isRequired,
  offers: PropTypes.array.isRequired,
};

export default Purchase;
