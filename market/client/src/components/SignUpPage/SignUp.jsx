import React, { useState, } from 'react';
import TextField from '@material-ui/core/TextField';
import Button from '@material-ui/core/Button';
import Alert from '@material-ui/lab/Alert';
import API from '../../utils/API';
import { storeContext, } from '../../utils/store';
import LinearProgress from '@material-ui/core/LinearProgress';
import styled from 'styled-components';
import { useHistory, } from 'react-router';

const SpacedProgress = styled(LinearProgress)`
  margin: 5px;
`;

/**
 * Creating a SignUp application.
 */
const SignUp = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [name, setName] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [loading, setLoading] = useState(false);
  const { authToken: [, setAuthToken], } = storeContext();
  const history = useHistory();
  const handleCreation = async () => {
    setLoading(true);
    try {
      const response = await API.registerUser(email, name, password);
      setAuthToken(response.token);
      setLoading(false);
      history.push('/dashboard');
    } catch (err) {
      setLoading(false);
      setErrorMessage(`Error from blockchain: ${err.message}`);
    }
  };
  return (
    <>
      <TextField label='Name' type='text' id='enter-name'
        onChange={(event) => setName(event.target.value)} fullWidth />
      <TextField label='Email' type='text' id='enter-email'
        onChange={(event) => setEmail(event.target.value)} fullWidth />
      <TextField label='Password' type='password' id='enter-password'
        onChange={(event) => setPassword(event.target.value)} fullWidth />
      <Button variant='contained' color='primary' type='submit'
        onClick={handleCreation} fullWidth >
        Create Account
      </Button>
      {loading ? <SpacedProgress /> : <></>}
      {errorMessage
        ? <Alert severity='error'> {errorMessage} </Alert>
        : <></>
      }
    </>
  );
};

export default SignUp;
