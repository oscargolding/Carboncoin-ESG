import { CardContent, Typography, CardActions, Button, } from '@material-ui/core';
import PropTypes from 'prop-types';
import React from 'react';
import DoneIcon from '@material-ui/icons/Done';
import BlockIcon from '@material-ui/icons/Block';
import EvStationIcon from '@material-ui/icons/EvStation';
import {
  SpacedCard,
  OfferStatus,
} from '../DashboardPage/styles/DashboardStyles';

/**
 * A card to show the type of production performed
 * @param {*} props to pass into a card
 * @returns the produciton card
 */
const ProductionCard = (props) => {
  const { produced, date, paid, usingRef, } = props;
  return (
    <SpacedCard innerRef={usingRef}>
      <CardContent>
        <Typography variant='h5' component='h2'>
          <EvStationIcon /> Carbon Production on {date}
        </Typography>
        <Typography varaint="body2" component="p">
          Amount of Carbon produced: <b>{produced}</b>
        </Typography>
        {paid
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
        {!paid
          ? <Button>
            Pay with Carboncoin
          </Button>
          : <></>}
      </CardActions>
    </SpacedCard>
  );
};

export default ProductionCard;

ProductionCard.propTypes = {
  usingRef: PropTypes.func,
  produced: PropTypes.number.isRequired,
  date: PropTypes.string.isRequired,
  paid: PropTypes.bool.isRequired,
};
