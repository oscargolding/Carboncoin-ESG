import React from 'react';
import { Card, CardContent, } from '@material-ui/core';
import SignUp from '../components/SignUpPage/SignUp';

/**
 * Create the SignUpPage
 */
const SignUpPage = () => {
  return (
    <Card>
      <CardContent>
        <h2> Register for the Carbon Market! </h2>
        <SignUp />
      </CardContent>
    </Card>
  );
};

export default SignUpPage;
