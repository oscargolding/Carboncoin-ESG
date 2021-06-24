const API = {};

const access = 'http://localhost:3000';

/**
 * Register the user with the blockchain.
 * @param {email} email for registering the email
 * @param {firm} firm the firm to register
 * @param {password} password the password of the registering firm
 * @returns token on success
 */
API.registerUser = async (email, firm, password) => {
  const queryParam = `${access}/api/admin/auth/register`;
  const response = await fetch(queryParam, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json', },
    body: JSON.stringify({ email: email, firm: firm, password: password, }),
  });
  if (!response.ok) {
    const errorResponse = await response.json();
    throw new Error(errorResponse.error);
  }
  const json = await response.json();
  return json;
};

/**
 * Get the balance of the user in the system.
 * @param {token} token for auth.
 * @returns the balance of the user in json format
 */
API.getBalance = async (token) => {
  console.log('trying to get balance');
  const queryParam = `${access}/api/token/balance`;
  console.log('calling fn');
  const response = await fetch(queryParam, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Login the user into the server / blockchain.
 * @param {email} email email of the user
 * @param {password} password password of the user
 * @returns the json response
 */
API.loginUser = async (email, password) => {
  const queryParam = `${access}/api/admin/auth/login`;
  const response = await fetch(queryParam, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ email: email, password: password, }),
  });
  if (!response.ok) {
    const errorResponse = await response.json();
    throw new Error(errorResponse.error);
  }
  const json = await response.json();
  return json;
};

export default API;
