import React, { useState, useEffect, } from 'react';
import AccountBalanceWalletIcon from '@material-ui/icons/AccountBalanceWallet';
import styled from 'styled-components';
import CircularProgress from '@material-ui/core/CircularProgress';
import { storeContext, } from '../utils/store';
import API from '../utils/API';

const BalanceText = styled.p`
  margin-left: 5px;
  margin-right: 5px;
`;

/**
 * Get the balance associated with an account.
 */
const Balance = () => {
  const [loaded, isLoaded] = useState(false);
  const [balance, setBalance] = useState('');
  const { authToken: [authToken], } = storeContext();
  useEffect(() => {
    console.log('calling use effect');
    const performRetrieve = async () => {
      try {
        console.log('getting balance');
        const response = await API.getBalance(authToken);
        console.log(response);
        setBalance(response.balance);
        isLoaded(true);
      } catch (err) {
        setBalance('Failed');
      }
    };
    performRetrieve();
  }, []);
  return (
    <>
      {loaded
        ? <>
          <AccountBalanceWalletIcon />
          <BalanceText> {balance} Tokens </BalanceText>
        </>
        : <>
          <CircularProgress color='secondary' />
        </>
      }
    </>
  );
};

export default Balance;
