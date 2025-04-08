declare namespace Cypress {
  interface Chainable {
    signInAsUser(conf?: { user2: boolean }): Chainable;
  }
}
