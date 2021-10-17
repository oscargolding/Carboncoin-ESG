import React, { useState, } from 'react';
import PropTypes from 'prop-types';
import CardContent from '@material-ui/core/CardContent';
import Typography from '@material-ui/core/Typography';
import {
  CardActions, Button, IconButton, Dialog, DialogTitle, DialogContent,
  DialogActions,
  CircularProgress,
} from '@material-ui/core';
import { SpacedCard, OfferStatus, HeaderCard, } from './styles/DashboardStyles';
import { useHistory, } from 'react-router-dom';
import DoneIcon from '@material-ui/icons/Done';
import BlockIcon from '@material-ui/icons/Block';
import ReputationElement from './ReputationElement';
import DeleteForeverIcon from '@material-ui/icons/DeleteForever';
import { storeContext, } from '../../utils/store';
import { Alert, } from '@material-ui/lab';
import API from '../../utils/API';
import Link from '@mui/material/Link';
import InfoIcon from '@mui/icons-material/Info';
import RepChart from '../ProductionPage/RepChart';

/**
 * The offer card being used for the sale of carbon coin.
 * @param {*} props the props passed into the card
 * @returns the card being offered
 */
const OfferCard = (props) => {
  const [open, setOpen] = useState(false);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [openReputation, closeReputation] = useState(false);
  const { authToken: [authToken], } = storeContext();
  const {
    producer, price, quantity, active,
    offerid, reputation, usingRef, owned, deleteOfferFn,
    environment, social, governance,
  } = props;
  const history = useHistory();
  const handleClickOpen = () => {
    setOpen(true);
  };
  const handleClose = () => {
    setOpen(false);
  };
  const deleteOffer = async () => {
    setLoading(true);
    setError('');
    try {
      await API.deleteOffer(authToken, offerid);
      deleteOfferFn(offerid);
      setLoading(false);
      handleClose();
    } catch (error) {
      setError(error.message);
      setLoading(false);
    }
  };
  return (
    <SpacedCard innerRef={usingRef}>
      <CardContent>
        <HeaderCard>
          <Typography variant='h5' component='h2'>
            Carboncoin Sale by <Link
              underline='hover'
              onClick={() => history.push('/offer/user', { name: producer, })}
            >
              {producer}
            </Link>
          </Typography>
          <ReputationElement repScore={reputation} username={producer} />
        </HeaderCard>
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
        <Button
          variant="contained"
          color="default"
          startIcon={<InfoIcon />}
          onClick={() => closeReputation(!openReputation)}
        >
          Seller Reputation
        </Button>
        <Dialog
          open={open}
          onClose={handleClose}
        >
          <DialogTitle>
            {'Are you sure you want to delete this offer?'}
          </DialogTitle>
          <DialogContent>
            <p>Offer will be deleted permenantly.</p>
            {loading ? <CircularProgress /> : <></>}
            {error !== '' ? <Alert severity='error'> {error} </Alert> : <></>}
          </DialogContent>
          <DialogActions>
            <Button
              onClick={handleClose}
              color='primary'
            >
              Close
            </Button>
            <Button
              onClick={deleteOffer}
              color='primary'
              autoFocus
            >
              Delete Offer
            </Button>
          </DialogActions>
        </Dialog>
        {owned
          ? <IconButton
            aria-label='delete'
          >
            <DeleteForeverIcon onClick={handleClickOpen} />
          </IconButton>
          : <></>
        }
      </CardActions>
      {openReputation
        ? <><RepChart
          environment={environment}
          social={social}
          governance={governance}
        /> </>
        : <></>}
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
  reputation: PropTypes.number.isRequired,
  owned: PropTypes.bool.isRequired,
  usingRef: PropTypes.func,
  deleteOfferFn: PropTypes.func,
  environment: PropTypes.number.isRequired,
  social: PropTypes.number.isRequired,
  governance: PropTypes.number.isRequired,
};
