import React from 'react';
import PropTypes from 'prop-types';
import Dashboard from '../DashboardPage/Dashboard';

/**
 * A way to observe individual offers for a user.
 * @returns the Individual offers for a user
 */
const IndividualOffers = (props) => {
  const { username, } = props;
  return (
    <>
      <Dashboard main={false} username={username} />
    </>
  );
};

IndividualOffers.propTypes = {
  username: PropTypes.string.isRequired,
};

export default IndividualOffers;
