import React from 'react';
import { useLocation, } from 'react-router';
import IndividualOffers from '../components/IndividualOfferPage/IndividualOffers';
import { storeContext, } from '../utils/store';

/**
 * Ability to view offers related to a user.
 * @returns a page to view all offers related to a user on.
 */
const OfferViewPage = () => {
  const { username: [username], } = storeContext();
  const location = useLocation();
  const name = location.state.name;
  const title = username === name
    ? 'Viewing your active offers'
    : `Viewing offers for user ${name}`;
  return (
    <>
      <h1> {title} </h1>
      <IndividualOffers username={username} />
    </>
  );
};

export default OfferViewPage;
