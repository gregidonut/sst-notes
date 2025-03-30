describe("general", () => {
  beforeEach(() => {
    cy.signInAsUser();
  });
  it("should be able to access the homepage", () => {
    cy.visit("/");
  });
});
