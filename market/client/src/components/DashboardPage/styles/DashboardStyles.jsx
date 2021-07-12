import styled from 'styled-components';
import { Card, Button, Chip, } from '@material-ui/core';

// Envrionment for creating sale offer
export const CreateOfferButton = styled(Button)`
  flex-grow: 1;
`;

export const SpacedCard = styled(Card)`
  margin: 15px;
`;

export const OfferStatus = styled(Chip)`
  margin-top: 10px;
`;

export const CentralLoading = styled.div`
  margin-top: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  width: 100%;
`;
