/**
 * Invalid access from the application.
 */
class AccessError extends Error {
  constructor(message) {
    super(message);
    this.name = 'AccessError';
  }
}

export default AccessError;
