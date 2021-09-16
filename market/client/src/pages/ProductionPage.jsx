import React from 'react';
import ProductionList from '../components/ProductionPage/ProductionList';

/**
 * Represents a level of production being recorded.
 * @returns the Production Page
 */
const ProductionPage = () => {
  return (
    <>
      <h1> Carbon Reporting </h1>
      <ProductionList />
    </>
  );
};

export default ProductionPage;
