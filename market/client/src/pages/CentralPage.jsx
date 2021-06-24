import React from 'react';
import Pages from '../components/Pages';
import NavBar from '../components/NavBar';
import { StylesProvider, } from '@material-ui/styles';
import { BrowserRouter as Router, } from 'react-router-dom';

/**
 * The application Central Page
 * @returns the Central Page to the application.
 */
const CentralPage = () => {
  return (
    <>
      <StylesProvider injectFirst>
        <Router>
          <NavBar />
          <Pages />
        </Router>
      </StylesProvider>
    </>
  );
};

export default CentralPage;
