import { CardContent, Typography, CardActions, Button, LinearProgress, } from '@material-ui/core';
import PropTypes from 'prop-types';
import React, { useState, } from 'react';
import DoneIcon from '@material-ui/icons/Done';
import BlockIcon from '@material-ui/icons/Block';
import EvStationIcon from '@material-ui/icons/EvStation';
import {
  SpacedCard,
  OfferStatus,
} from '../DashboardPage/styles/DashboardStyles';
import { storeContext, } from '../../utils/store';
import API from '../../utils/API';
import { Alert, } from '@material-ui/lab';

/**
 * A card to show the type of production performed
 * @param {*} props to pass into a card
 * @returns the produciton card
 */
const ProductionCard = (props) => {
  const { produced, date, paid, usingRef, id, } = props;
  const { authToken: [authToken], balance: [, setBalance], } = storeContext();
  const [hasPaid, setHasPaid] = useState(paid);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const payDebt = async () => {
    setLoading(true);
    setError('');
    try {
      const response = await API.payProduction(authToken, id);
      setLoading(false);
      setHasPaid(true);
      setBalance(response.balance);
      console.log(response.balance);
    } catch (err) {
      setLoading(false);
      setError(err.message);
    }
  };
  return (
    <SpacedCard innerRef={usingRef}>
      <CardContent>
        <Typography variant='h5' component='h2'>
          <EvStationIcon /> Carbon Production on {date}
        </Typography>
        <Typography varaint="body2" component="p">
          Amount of Carbon produced: <b>{produced}</b>
        </Typography>
        {hasPaid
          ? <OfferStatus
            icon={<DoneIcon />}
            label='Paid For'
            clickable={false}
            color='primary'
          />
          : <OfferStatus
            icons={<BlockIcon />}
            label='Requires Payment with Carboncoin'
            clickable={false}
            color='secondary'
          />}
      </CardContent>
      <CardActions>
        {!hasPaid
          ? <Button onClick={() => { payDebt(); }}>
            Pay with Carboncoin
          </Button>
          : <></>}
        {error !== '' ? <Alert severity='error'>{error}</Alert> : <></>}
      </CardActions>
      {loading ? <LinearProgress /> : <></>}
    </SpacedCard>
  );
};

export default ProductionCard;

ProductionCard.propTypes = {
  usingRef: PropTypes.func,
  produced: PropTypes.number.isRequired,
  date: PropTypes.string.isRequired,
  paid: PropTypes.bool.isRequired,
  id: PropTypes.string.isRequired,
};
