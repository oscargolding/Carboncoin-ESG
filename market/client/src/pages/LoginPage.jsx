import React from 'react';
import { Card, CardContent, } from '@material-ui/core';
import Login from '../components/LoginPage/Login';
import API from '../utils/API';

/**
 * Login page for the application.
 */
const LoginPage = () => {
  return (
    <Card>
      <CardContent>
        <h2> Login </h2>
        <Login loginUser={API.loginUser} />
      </CardContent>
    </Card>
  );
};

export default LoginPage;
