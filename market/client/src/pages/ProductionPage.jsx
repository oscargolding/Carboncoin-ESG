import React from 'react';
import { useLocation, } from 'react-router';
import ProductionList from '../components/ProductionPage/ProductionList';
import { storeContext, } from '../utils/store';

/**
 * Represents a level of production being recorded.
 * @returns the Production Page
 */
const ProductionPage = () => {
  const { username: [username], } = storeContext();
  const location = useLocation();
  const name = location.state.name;
  const title = username === name
    ? 'Viewing your Reputation Breakdown'
    : `Viewing Reputation Breakdown for firm ${name}`;
  return (
    <>
      <h1> {title} </h1>
      <ProductionList username={name} />
    </>
  );
};

export default ProductionPage;
