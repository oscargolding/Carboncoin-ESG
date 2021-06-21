/**
 * Invalid input from the application calling the api
 */
class InputError extends Error {
  constructor(message) {
    super(message);
    this.name = 'InputError';
  }
}

module.exports = { InputError };
