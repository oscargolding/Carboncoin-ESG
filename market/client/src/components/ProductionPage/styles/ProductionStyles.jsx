import styled, { css, } from 'styled-components';

export const PaperList = styled.div`
  margin: 5px;
  display: flex;
  flex-direction: row;
  justify-content: flex-start;
  align-items: center;
  flex-wrap: wrap;
  list-style: none;
`;

export const ReputationDiv = styled.div`
  width: fit-content;
  border-radius: 20px;
  ${props => props.reputation >= 0 && css`border: 3px solid #ADD8E6;`}
  ${props => props.reputation < 0 && css`border: 3px solid #F08080;`}
  padding: 1px;
  margin-bottom: 10px;
`;
