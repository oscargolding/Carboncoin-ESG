import styled from 'styled-components';
import { Typography, } from '@material-ui/core';

export const HomeContainer = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  height: 50vh;
`;

export const LargeLog = styled.img`
  width: 200px;
  height: 200px;
`;

export const Title = styled(Typography)`
  margin-top: 5px;
`;
