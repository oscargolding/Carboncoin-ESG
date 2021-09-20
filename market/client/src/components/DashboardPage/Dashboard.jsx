import React, { useState, } from 'react';
import { CreateOfferButton, ButtonForm, OfferRow, } from './styles/DashboardStyles';
import { useHistory, } from 'react-router';
import OfferList from './OfferList';
import PropTypes from 'prop-types';
import InputLabel from '@material-ui/core/InputLabel';
import Select from '@material-ui/core/Select';
import { Button, MenuItem, } from '@material-ui/core';

/**
 * Component represents the dashboard used inside the market. - see offers
 * @returns the Dashboard
 */
const Dashboard = (props) => {
  const history = useHistory();
  const [sortTerm, setSortTerm] = useState('');
  const [direction, setDirection] = useState(false);
  const [sortId, setSortId] = useState('');
  const { main, username, } = props;
  return (
    <>
      {main
        ? <CreateOfferButton
          variant='contained'
          color='primary'
          onClick={() => { history.push('/offer/createoffer'); }}
          fullWidth
        >
          Sell Carbon Currency
        </CreateOfferButton>
        : <></>
      }
      <OfferRow>
        <ButtonForm variant='outlined'>
          <InputLabel
            id="demo-simple-select-outlined-label"
          >
            Sort Offers
          </InputLabel>
          <Select
            labelId="demo-simple-select-outlined-label"
            id="demo-simple-select-outlined"
            label="Offer Sort"
            value={sortId}
            onChange={(event) => {
              console.log(event.target.value);
              setSortId(event.target.value);
              if (event.target.value === 0) {
                setSortTerm('reputation');
                setDirection(true);
              } else if (event.target.value === 1) {
                setSortTerm('reputation');
                setDirection(false);
              } else if (event.target.value === 2) {
                setSortTerm('price');
                setDirection(true);
              } else if (event.target.value === 3) {
                setSortTerm('price');
                setDirection(false);
              }
            }}
          >
            <MenuItem value="">
              <em>None</em>
            </MenuItem>
            <MenuItem value={0}>Carbon Reputation Ascending</MenuItem>
            <MenuItem value={1}>Carbon Reputation Descending</MenuItem>
            <MenuItem value={2}>Price Ascending</MenuItem>
            <MenuItem value={3}>Price Descending</MenuItem>
          </Select>
        </ButtonForm>
        {main
          ? <Button
            variant='contained'
            onClick={() => { history.push('/offerfinder'); }}
          >
            Offer Finder
          </Button>
          : <></>
        }
      </OfferRow>
      <OfferList
        sortTerm={sortTerm}
        direction={direction}
        username={main ? '' : username}
      />
    </>
  );
};

Dashboard.propTypes = {
  main: PropTypes.bool.isRequired,
  username: PropTypes.string,
};

export default Dashboard;
