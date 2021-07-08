import React from 'react';
import { CreateOfferButton, } from './styles/DashboardStyles';
import { useHistory, } from 'react-router';
import OfferList from './OfferList';

/**
 * Component represents the dashboard used inside the market. - see offers
 * @returns the Dashboard
 */
const Dashboard = () => {
  const history = useHistory();
  return (
    <>
      <CreateOfferButton
        variant='contained'
        color='primary'
        onClick={() => { history.push('/offer/createoffer'); }}
        fullWidth
      >
        Sell Carbon Currency
      </CreateOfferButton>
      <OfferList />
    </>
  );
};

export default Dashboard;
