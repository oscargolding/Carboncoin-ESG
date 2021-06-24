import React from 'react';
import Typography from '@material-ui/core/Typography';
import AppBar from '@material-ui/core/AppBar';
import styled from 'styled-components';
import Toolbar from '@material-ui/core/Toolbar';
import { storeContext, } from '../utils/store';
import { Link, } from 'react-router-dom';
import Button from '@material-ui/core/Button';
import Balance from './Balance';

// header div
const HeaderDiv = styled.div`
  flex-grow: 1;
`;

// application name
const AppName = styled(Typography)`
  flex-grow: 1;
`;

/**
 * The NavBar for the carbon market application.
 * @returns the NavBar for the application.
 */
const NavBar = () => {
  // Get the store context
  const { authToken: [authToken], } = storeContext();
  return (
    <HeaderDiv>
      <AppBar position='static'>
        <Toolbar>
          <AppName variant="h6">
            Blockchain Carbon Market
          </AppName>
          {!authToken
            ? <>
              <Button color="inherit" component={Link} to='/login'>
                Login
              </Button>
              <Button color="inherit" component={Link} to='/signup'>
                Sign Up
              </Button>
            </>
            : <>
              <Balance />
            </>}
        </Toolbar>
      </AppBar>
    </HeaderDiv>
  );
};

export default NavBar;
