import styled, { css, } from 'styled-components';
import { Card, Button, Chip, FormControl, Typography, } from '@material-ui/core';

// Envrionment for creating sale offer
export const CreateOfferButton = styled(Button)`
  flex-grow: 1;
`;

export const SpacedCard = styled(Card)`
  margin: 15px;
  &:hover {
    box-shadow: 0 4px 8px 0 rgba(0, 0, 0, 0.2), 0 6px 20px 0 rgba(0, 0, 0, 0.19);
  }
`;

export const OfferStatus = styled(Chip)`
  margin: 5px;
  ${props => props.ethical && props.label === 'Environmental' &&
    css`
      background-color: green;
  `}
  ${props => props.ethical && props.label === 'Social' &&
    css`
      background-color: pink;
  `}
  ${props => props.ethical && props.label === 'Governance' &&
    css`
      background-color: grey;
  `}
`;

export const CentralLoading = styled.div`
  margin-top: 15px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: row;
  width: 100%;
`;

export const OfferRow = styled.div`
  margin-top: 10px;
  display: flex;
  flex-direction: row;
  justify-content: space-between;
`;

export const ButtonForm = styled(FormControl)`
  min-width: 120px;
`;

export const RightStamp = styled.div`
  display: flex;
  flex-direction: row;
`;

export const HeaderCard = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
`;

export const ButtonScore = styled(Button)`
  min-width: 200px;
  ${props =>
    props.score >= -250 &&
    css`
      background: green;
      color: white;
  `};
  ${props =>
    props.score >= -500 && props.score < -250 &&
    css`
      background: palevioletred;
      color: white;
  `};
  ${props =>
    props.score <= -500 &&
    css`
      background: red;
      color: white;
  `};
`;

export const PopoverText = styled(Typography)`
  white-space: pre-line;
  padding: 15px;
`;
