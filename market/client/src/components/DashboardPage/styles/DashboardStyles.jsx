import styled from 'styled-components';
import { Card, Button, } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';

// Envrionment for creating sale offer
export const CreateOfferButton = styled(Button)`
  flex-grow: 1;
`;

export const SpacedCard = styled(Card)`
  margin: 15px;
`;

export const OfferStatus = styled(Alert)`
  margin-top: 5px;
`;

export const CentralLoading = styled.div`
  margin-top: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  width: 100%;
`;
