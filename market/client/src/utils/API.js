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

/**
 * Create offer on the open market.
 * @param {token} token user token
 * @param {amount} amount amount in dollars
 * @param {quantity to trade} quantity trading
 * @returns
 */
API.createOffer = async (token, amount, quantity) => {
  const queryParam = `${access}/api/offer/create`;
  const response = await fetch(queryParam, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ amount: amount, tokens: quantity, }),
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

API.getOffers = async (userAuthToken, signal, paginationToken, sortTerm, direction) => {
  const orderTerm = sortTerm !== '' ? `&field=${sortTerm}` : '';
  const directionTerm = direction !== false ? '&direction=1' : '';
  const queryParam = `${access}/api/offers/list?token=${paginationToken}${orderTerm}${directionTerm}`;
  const response = await fetch(queryParam, {
    method: 'GET',
    signal: signal,
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

API.getProduction = async (userAuthToken, signal, paginationToken) => {
  const queryParam = `${access}/api/production/list?token${paginationToken}`;
  const response = await fetch(queryParam, {
    method: 'GET',
    signal: signal,
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Accept the offer for the purchase of tokens
 * @param {authToken} userAuthToken the user auth token
 * @param {offerId} offerId the offerId to purchase from
 * @param {quantity} quantity the number to purchase
 * @returns the jsonResponse
 */
API.acceptOffer = async (userAuthToken, offerId, quantity) => {
  const queryParam = `${access}/api/offer/buy`;
  console.log(offerId);
  console.log(quantity);
  const quantityBuying = String(quantity);
  const response = await fetch(queryParam, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      purchased: quantityBuying,
      id: offerId,
    }),
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Pay for the carbon production used.
 * @param {token} userAuthToken the user auth token
 * @param {*} prodId the production id of the user
 * @returns the amount of balance left
 */
API.payProduction = async (userAuthToken, prodId) => {
  const queryParam = `${access}/api/production/pay?production=${prodId}`;
  console.log(prodId);
  const response = await fetch(queryParam, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Get the direct price from the market currently offered.
 * @param {token} userAuthToken auth token
 * @returns a json blob with price
 */
API.getDirectPrice = async (userAuthToken) => {
  const queryParam = `${access}/api/direct/price`;
  const response = await fetch(queryParam, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Accept the direct price on the market.
 * @param {token} userAuthToken auth token in the market
 * @returns the json object representing the tokens
 */
API.acceptDirectPrice = async (userAuthToken, quantity) => {
  const queryParam = `${access}/api/direct/redeem`;
  const response = await fetch(queryParam, {
    method: 'POST',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      amount: quantity,
    }),
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Get the given offers as found.
 * @param {*} userAuthToken token for authentication
 * @param {*} target the target amount to get
 * @param {*} reputation the reputation being used
 * @returns the json object representing the offers
 */
API.getFoundOffers = async (userAuthToken, target, reputation) => {
  const queryParam = `${access}/api/offers/target?` +
    `target=${target}&reputation=${reputation}`;
  const response = await fetch(queryParam, {
    method: 'GET',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

/**
 * Delete the given offer from the entire chain.
 * @param {*} userAuthToken token to use
 * @param {*} offerId offerId
 */
API.deleteOffer = async (userAuthToken, offerId) => {
  const queryParam = `${access}/api/offer/delete?id=${offerId}`;
  const response = await fetch(queryParam, {
    method: 'DELETE',
    headers: {
      Authorization: `Bearer ${userAuthToken}`,
    },
  });
  const jsonResponse = await response.json();
  if (!response.ok) {
    throw new Error(jsonResponse.error);
  }
  return jsonResponse;
};

export default API;
