import { setupClerkTestingToken } from "@clerk/testing/cypress";

describe("notes page", () => {
  beforeEach(() => {
    setupClerkTestingToken();

    cy.viewport("iphone-6");
    cy.visit("/sign-in");
    // Add any other actions to test
    cy.clerkSignIn({
      strategy: "password",
      identifier: Cypress.env("test_user"),
      password: Cypress.env("test_password"),
    });
  });
  it("should be able to access the homepage", () => {
    cy.visit("/");
  });
});
