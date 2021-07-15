import React from 'react';
import Typography from '@material-ui/core/Typography';
import AppBar from '@material-ui/core/AppBar';
import styled from 'styled-components';
import Toolbar from '@material-ui/core/Toolbar';
import { storeContext, } from '../utils/store';
import { Link, useHistory, } from 'react-router-dom';
import Button from '@material-ui/core/Button';
import Balance from './Balance';
import CoinLogo from '../images/svg.svg';
import IconButton from '@material-ui/core/IconButton';
import HomeIcon from '@material-ui/icons/Home';
import MenuItem from '@material-ui/core/MenuItem';
import Menu from '@material-ui/core/Menu';
import AccountCircle from '@material-ui/icons/AccountCircle';

// header div
const HeaderDiv = styled.div`
  flex-grow: 1;
`;

// application name
const AppName = styled(Typography)`
  flex-grow: 1;
`;

const Logo = styled.img`
  width: 50px;
  height: 50px;
  margin-right: 15px;
`;

/**
 * The NavBar for the carbon market application.
 * @returns the NavBar for the application.
 */
const NavBar = () => {
  // Get the store context
  const { authToken: [authToken], } = storeContext();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const open = Boolean(anchorEl);
  const history = useHistory();
  const handleMenu = (event) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };
  return (
    <HeaderDiv>
      <AppBar position='static'>
        <Toolbar>
          <Logo src={CoinLogo} />
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
              <IconButton
                aria-label="home"
                component={Link}
                to='/dashboard'
                id='home-btn'
              >
                <HomeIcon />
              </IconButton>
              <IconButton
                aria-label="account of current user"
                aria-controls="menu-appbar"
                aria-haspopup="true"
                onClick={handleMenu}
                color="inherit"
              >
                <AccountCircle />
              </IconButton>
              <Menu
                id="menu-appbar"
                anchorEl={anchorEl}
                anchorOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                keepMounted
                transformOrigin={{
                  vertical: 'top',
                  horizontal: 'right',
                }}
                open={open}
                onClose={handleClose}
              >
                <MenuItem onClick={() => history.push('/production')}>
                  Carbon Production
                </MenuItem>
                <MenuItem onClick={handleClose}>
                  Logout
                </MenuItem>
              </Menu>
            </>}
        </Toolbar>
      </AppBar>
    </HeaderDiv>
  );
};

export default NavBar;
