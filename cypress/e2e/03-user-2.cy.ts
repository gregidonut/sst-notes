import type { Note } from "@models";
describe("crud", () => {
  const seedLength = 5;
  beforeEach(() => {
    cy.signInAsUser({ user2: true });
    cy.visit("/notes");
    cy.get('[data-cy="note-list"]')
      .children()
      .then((element) => {
        return element[0].tagName === "OL" ? false : true;
      })
      .then((dbEmpty) => {
        switch (dbEmpty) {
          case true:
            cy.log("db is empty");
            cy.get("@clerkToken").then((token) => {
              cy.task("seedDyanamoDB", token);
            });
            cy.visit("/notes");

            cy.get('[data-cy="note-list"] > ol')
              .as("noteList")
              .children()
              .should("have.length", seedLength);
            break;
          case false:
            cy.log("db is not empty");
            cy.get("@clerkToken").then((token) => {
              cy.task("emptyDynamoDB", token);
            });
            cy.visit("/notes");
            cy.get('[data-cy="note-list"] > p').contains("no items yet");

            cy.get("@clerkToken").then((token) => {
              cy.task("seedDyanamoDB", token);
            });
            cy.visit("/notes");

            cy.get('[data-cy="note-list"] > ol')
              .as("noteList")
              .children()
              .should("have.length", seedLength);
            break;
          default:
        }
      });
  });
  afterEach(() => {
    cy.visit("/notes");
    cy.get('[data-cy="note-list"]')
      .children()
      .then((element) => {
        return element[0].tagName === "OL" ? false : true;
      })
      .then((dbEmpty) => {
        if (dbEmpty) return;
        cy.log("db is not empty");
        cy.get("@clerkToken").then((token) => {
          cy.task("emptyDynamoDB", token);
        });
        cy.visit("/notes");
        cy.get('[data-cy="note-list"] > p').contains("no items yet");
      });
  });
  it("should see notes only for logged in user", () => {
    cy.visit("/notes");
    cy.get('[data-cy="username"]')
      .should("have.length", 5)
      .each((el) => {
        expect(el[0].innerText === Cypress.env("test_user2"));
      });
  });
});
