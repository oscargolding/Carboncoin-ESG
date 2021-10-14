import React from 'react';
import styled from 'styled-components';
import SignUpPage from '../pages/SignUpPage';
import { Switch, Route, } from 'react-router-dom';
import LoginPage from '../pages/LoginPage';
import DashboardPage from '../pages/DashboardPage';
import CreateOfferPage from '../pages/CreateOfferPage';
import Landing from './LandingPage/Landing';
import PurchaseOfferPage from '../pages/PurchaseOfferPage';
import ProductionPage from '../pages/ProductionPage';
import DirectPurchasePage from '../pages/DirectPurchasePage';
import OfferFinderPage from '../pages/OfferFinderPage';
import OfferViewPage from '../pages/OfferViewPage';
import FAQPage from '../pages/FAQPage';

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
        <Route path='/production'>
          <ProductionPage />
        </Route>
        <Route path='/offer/createoffer'>
          <CreateOfferPage />
        </Route>
        <Route path='/offer/purchase'>
          <PurchaseOfferPage />
        </Route>
        <Route path='/offer/user'>
          <OfferViewPage />
        </Route>
        <Route path='/direct/purchase'>
          <DirectPurchasePage />
        </Route>
        <Route path='/offerfinder'>
          <OfferFinderPage />
        </Route>
        <Route path='/faq'>
          <FAQPage />
        </Route>
        <Route path='/'>
          <Landing />
        </Route>
      </Switch>
    </BodyPage>
  );
};

export default Pages;
