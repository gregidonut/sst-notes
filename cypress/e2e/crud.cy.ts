describe("crud", () => {
  let seedLength = 5;
  beforeEach(() => {
    cy.signInAsUser();
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
  it("should create note", () => {
    cy.visit("/create");
    const content = "pukingina mo hayup ka";
    cy.get('[data-cy="content-field"]').type(content);
    cy.get('[data-cy="create-content-button"]').click();
    cy.url().should("eq", "http://localhost:4321/notes");
    cy.get('[data-cy="note-list"] > ol')
      .as("noteList")
      .children()
      .should("have.length", seedLength + 1);
    cy.get(
      ':nth-child(1) > article > header > p > [data-cy="update-link"]',
    ).then((el) => {
      cy.log(el[0].innerHTML);
    });
  });

  it("should see specific note page", () => {
    cy.visit("/notes");
    cy.get(
      ':nth-child(1) > article > header > p > [data-cy="update-link"]',
    ).click();
  });
});
