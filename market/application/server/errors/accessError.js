/**
 * Invalid access from the application.
 */
class AccessError extends Error {
  constructor(message) {
    super(message);
    this.name = 'AccessError';
  }
}

module.exports = { AccessError };
