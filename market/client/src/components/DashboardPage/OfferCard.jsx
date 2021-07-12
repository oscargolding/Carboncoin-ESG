import React from 'react';
import PropTypes from 'prop-types';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import { CardActions, Button, } from '@material-ui/core';
import { SpacedCard, OfferStatus, } from './styles/DashboardStyles';
import { useHistory, } from 'react-router-dom';
import DoneIcon from '@material-ui/icons/Done';
import BlockIcon from '@material-ui/icons/Block';

/**
 * The offer card being used for the sale of carbon coin.
 * @param {*} props the props passed into the card
 * @returns the card being offered
 */
const OfferCard = (props) => {
  const { producer, price, quantity, active, offerid, usingRef, } = props;
  const history = useHistory();
  return (
    <SpacedCard innerRef={usingRef}>
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
          ? <OfferStatus
            icon={<DoneIcon />}
            label='Active Offer'
            clickable={false}
            color='primary'
          />
          : <OfferStatus
            icon={<BlockIcon />}
            label='Finished Offer'
            clickable={false}
            color='secondary'
          />}
      </CardContent>
      <CardActions>
        <Button
          size="small"
          onClick={() => history.push('/offer/purchase',
            { offerid: offerid, quantity: quantity, producer: producer, price: price, }
          )
          }
        >
          Purchase From Offer
        </Button>
      </CardActions>
    </SpacedCard>
  );
};

export default OfferCard;

OfferCard.propTypes = {
  producer: PropTypes.string.isRequired,
  price: PropTypes.number.isRequired,
  quantity: PropTypes.number.isRequired,
  active: PropTypes.bool.isRequired,
  offerid: PropTypes.string.isRequired,
  usingRef: PropTypes.func,
};
