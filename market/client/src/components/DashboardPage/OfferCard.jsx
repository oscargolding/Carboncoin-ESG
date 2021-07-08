import React from 'react';
import PropTypes from 'prop-types';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import { SpacedCard, OfferStatus, } from './styles/DashboardStyles';

/**
 * The offer card being used for the sale of carbon coin.
 * @param {*} props the props passed into the card
 * @returns the card being offered
 */
const OfferCard = (props) => {
  const { producer, price, quantity, active, } = props;
  return (
    <SpacedCard>
      <CardContent>
        <Typography variant='h5' component='h2'>
          Carboncoin Sale by {producer}
        </Typography>
        <Typography variant="body2" component="p">
          Price Per Token: <b>${price}</b>
          <br />
          Quantity Offered: <b>{quantity}</b>
        </Typography>
        {active
          ? <OfferStatus severity='success'> Active Offer</OfferStatus>
          : <OfferStatus severity='warning'> Inactive Offer</OfferStatus>}
      </CardContent>
    </SpacedCard>
  );
};

export default OfferCard;

OfferCard.propTypes = {
  producer: PropTypes.string.isRequired,
  price: PropTypes.number.isRequired,
  quantity: PropTypes.number.isRequired,
  active: PropTypes.bool.isRequired,
};
