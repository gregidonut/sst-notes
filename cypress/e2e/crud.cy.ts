describe("crud", () => {
  beforeEach(() => {
    cy.signInAsUser();
    cy.get("@clerkToken").then((token) => {
      cy.task("seedDyanamoDB", token);
    });
    cy.visit("/notes");
    cy.get('[data-cy="note-list"] > ol')
      .as("noteList")
      .children()
      .should("have.length", 5);
  });
  afterEach(() => {
    cy.get("@clerkToken").then((token) => {
      cy.task("emptyDynamoDB", token);
    });
    cy.visit("/notes");

    cy.get("@noteList").should("not.exist");
    cy.get('[data-cy="note-list"] > p').contains("no items yet");
  });

  it("should create note", () => {
    cy.visit("/create");
  });
});
