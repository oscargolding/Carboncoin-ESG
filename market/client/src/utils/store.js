import React, { useContext, } from 'react';
import PropTypes from 'prop-types';

export const StoreContext = React.createContext(null);
export const storeContext = () => useContext(StoreContext);

/**
 * A store object.
 * @param {createStore} param0 the store
 * @returns the store
 */
const Store = ({ children, }) => {
  const [authToken, setAuthToken] = React.useState('');
  const [balance, setBalance] = React.useState('');

  const store = {
    authToken: [authToken, setAuthToken],
    balance: [balance, setBalance],
  };

  return (
    <StoreContext.Provider value={store}>{children}</StoreContext.Provider>
  );
};

Store.propTypes = {
  children: PropTypes.node.isRequired,
};

export default Store;
