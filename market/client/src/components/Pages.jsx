import React from 'react';
import styled from 'styled-components';
import SignUpPage from '../pages/SignUpPage';
import { Switch, Route, } from 'react-router-dom';
import LoginPage from '../pages/LoginPage';
import DashboardPage from '../pages/DashboardPage';

// the dimensions of the page
const BodyPage = styled.div`
  max-width: 1280px;
  margin: 5px auto;
`;

/**
 * The central pages contained in the application - all the pages enumerated.
 */
const Pages = () => {
  return (
    <BodyPage>
      <Switch>
        <Route path='/signup'>
          <SignUpPage />
        </Route>
        <Route path='/login'>
          <LoginPage />
        </Route>
        <Route path='/dashboard'>
          <DashboardPage />
        </Route>
      </Switch>
    </BodyPage>
  );
};

export default Pages;
