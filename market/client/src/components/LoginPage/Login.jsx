import React, { useState, } from 'react';
import PropTypes from 'prop-types';
import styled from 'styled-components';
import { storeContext, } from '../../utils/store';
import { TextField, Button, LinearProgress, } from '@material-ui/core';
import Alert from '@material-ui/lab/Alert';
import { useHistory, } from 'react-router';

// progress more spaced
const SpacedProgress = styled(LinearProgress)`
  margin: 5px;
`;

/**
 * For performing a login into the network.
 */
const Login = (props) => {
  const { loginUser, } = props;
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrormessage] = useState('');
  const [loading, isLoading] = useState(false);
  const {
    authToken: [, setAuthToken],
    username: [, setUsername],
  } = storeContext();
  const history = useHistory();
  const handleSubmit = async () => {
    isLoading(true);
    try {
      const response = await loginUser(email, password);
      setAuthToken(response.token);
      setUsername(email);
      isLoading(false);
      history.push('/dashboard');
    } catch (err) {
      isLoading(false);
      setErrormessage(`Error from blockchain: ${err.message}`);
    }
  };
  return (
    <>
      <TextField label='Email' type='text' id='enter-email'
        onChange={(event) => setEmail(event.target.value)}
        fullWidth value={email} />
      <TextField label='Password' type='password' id='enter-password'
        onChange={(event) => setPassword(event.target.value)}
        value={password} fullWidth />
      <Button variant='contained' color='primary' type='submit'
        onClick={handleSubmit} fullWidth >
        Login
      </Button>
      {loading ? <SpacedProgress /> : <></>}
      {errorMessage ? <Alert severity='error'> {errorMessage} </Alert> : <></>}
    </>
  );
};

Login.propTypes = {
  loginUser: PropTypes.func.isRequired,
};

export default Login;
