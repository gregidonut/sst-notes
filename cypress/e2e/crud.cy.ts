describe("crud", () => {
  beforeEach(() => {
    cy.signInAsUser();
  });
  it("should find an initialized database", () => {
    cy.visit("/notes");
  });
});
