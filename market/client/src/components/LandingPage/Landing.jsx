import React from 'react';
import { HomeContainer, LargeLog, Title, } from './styles/LandingStyles';
import CoinLogo from '../../images/svg.svg';

/**
 * Landing page containing the logo and all other details.
 * @returns the Landing Page
 */
const Landing = () => {
  return (
    <>
      <HomeContainer>
        <LargeLog src={CoinLogo} />
        <Title variant='h4' gutterBottom>
          Carboncoin
        </Title>
        <Title variant='subtitle1' gutterBottom>
          Powered by Blockchain for Emissions Trading
        </Title>
      </HomeContainer>
    </>
  );
};

export default Landing;
