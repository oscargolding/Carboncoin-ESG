import styled, { css, } from 'styled-components';
import { Card, Button, Chip, FormControl, Typography, } from '@material-ui/core';

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

export const ButtonForm = styled(FormControl)`
  min-width: 120px;
  margin-top: 10px;
`;

export const HeaderCard = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
`;

export const ButtonScore = styled(Button)`
  min-width: 200px;
  ${props =>
    props.score < 500 &&
    css`
      background: green;
      color: white;
  `};
  ${props =>
    props.score >= 500 && props.score < 1000 &&
    css`
      background: palevioletred;
      color: white;
  `};
  ${props =>
    props.score >= 1000 &&
    css`
      background: red;
      color: white;
  `};
`;

export const PopoverText = styled(Typography)`
  white-space: pre-line;
  padding: 15px;
`;
