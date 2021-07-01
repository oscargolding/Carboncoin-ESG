import styled from 'styled-components';
import { Button, } from '@material-ui/core';

// An environment for styling the offer form
export const OfferForm = styled.div`
  display: flex;
  flex-direction: column;
  max-width: 800px;
  & > * {
    margin-bottom: 15px;
  }
`;

export const SellButton = styled(Button)`
  max-height: 100px;
`;
