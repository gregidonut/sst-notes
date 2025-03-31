/// <reference types="cypress" />
// ***********************************************
// This example commands.ts shows you how to
// create various custom commands and overwrite
// existing commands.
//
// For more comprehensive examples of custom
// commands please read more here:
// https://on.cypress.io/custom-commands
// ***********************************************
//
//
// -- This is a parent command --
// Cypress.Commands.add('login', (email, password) => { ... })
//
//
// -- This is a child command --
// Cypress.Commands.add('drag', { prevSubject: 'element'}, (subject, options) => { ... })
//
//
// -- This is a dual command --
// Cypress.Commands.add('dismiss', { prevSubject: 'optional'}, (subject, options) => { ... })
//
//
// -- This will overwrite an existing command --
// Cypress.Commands.overwrite('visit', (originalFn, url, options) => { ... })
//
// declare global {
//   namespace Cypress {
//     interface Chainable {
//       login(email: string, password: string): Chainable<void>
//       drag(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       dismiss(subject: string, options?: Partial<TypeOptions>): Chainable<Element>
//       visit(originalFn: CommandOriginalFn, url: string, options: Partial<VisitOptions>): Chainable<Element>
//     }
//   }
// }

import { addClerkCommands } from "@clerk/testing/cypress";
addClerkCommands({ Cypress, cy });

import { setupClerkTestingToken } from "@clerk/testing/cypress";
Cypress.Commands.add("signInAsUser", () => {
  setupClerkTestingToken();

  cy.viewport("iphone-6");
  cy.visit("/sign-in");
  cy.clerkSignIn({
    strategy: "password",
    identifier: Cypress.env("test_user"),
    password: Cypress.env("test_password"),
  });
  cy.window().then((window) => {
    window.Clerk.session.getToken().then((token) => {
      cy.wrap(token).as("clerkToken");
    });
  });
});

export {};
