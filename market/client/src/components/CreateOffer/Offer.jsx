import React, { useState, } from 'react';
import PropTypes from 'prop-types';
import { OfferForm, SellButton, } from './styles/OfferStyles';
import { OutlinedInput, InputLabel, InputAdornment, LinearProgress, }
  from '@material-ui/core';
import { storeContext, } from '../../utils/store';
import Alert from '@material-ui/lab/Alert';

/**
 * Create a usable offer.
 * @returns Offer
 */
const Offer = (props) => {
  const { createOffer, } = props;
  const [dollar, setDollar] = useState('0');
  const [quantity, setQuantity] = useState('0');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState(false);
  const [error, setError] = useState('');
  const { authToken: [authToken], } = storeContext();
  const handleSubmit = async () => {
    setLoading(true);
    setSuccess(false);
    setError('');
    try {
      await createOffer(authToken, dollar, quantity);
      setLoading(false);
      setSuccess(true);
    } catch (err) {
      setLoading(false);
      setError(err.message);
    }
  };
  return (
    <OfferForm>
      <h1> Create Market Offer </h1>
      <InputLabel htmlFor="outlined-adornment-amount">Amount</InputLabel>
      <OutlinedInput
        id="outlined-adornment-amount"
        value={dollar}
        onChange={(event) => { setDollar(event.target.value); }}
        startAdornment={<InputAdornment position="start">$</InputAdornment>}
        labelWidth={60}
      />
      <InputLabel htmlFor="outlined-adornment-quantity">Quantity</InputLabel>
      <OutlinedInput
        id="outlined-adornment-quantity"
        value={quantity}
        onChange={(event) => { setQuantity(event.target.value); }}
        labelWidth={60}
      />
      {dollar !== '0' && quantity !== '0'
        ? <p>Sell {quantity} Carboncoin
          for price ${dollar} dollars </p>
        : <></>}
      <SellButton
        variant="contained"
        color="primary"
        size='medium'
        labelWidth={60}
        onClick={() => handleSubmit()}
      >
        Sell Carboncoin!
      </SellButton>
      {loading ? <LinearProgress /> : <></>}
      {success ? <Alert severity='success'> Offer was created! </Alert> : <></>}
      {error !== '' ? <Alert severity='error'>{error}</Alert> : <></>}
    </OfferForm>
  );
};

Offer.propTypes = {
  createOffer: PropTypes.func.isRequired,
};

export default Offer;
