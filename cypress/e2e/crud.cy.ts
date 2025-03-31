describe("crud", () => {
  beforeEach(() => {
    cy.signInAsUser();
    cy.get("@clerkToken").then((token) => {
      cy.task("seedDyanamoDB", token);
    });
  });
  afterEach(() => {
    cy.get("@clerkToken").then((token) => {
      cy.task("emptyDynamoDB", token);
    });
  });

  it("should find an initialized database", () => {
    cy.visit("/notes");
  });
});
