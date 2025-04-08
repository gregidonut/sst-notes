import type { Note } from "@models";
describe("crud", () => {
  const seedLength = 5;
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

  it("should display notes sorted by datetime", () => {
    cy.visit("/notes");
    // check of notelist from the api response and from the dom are sorted and match exactly
    cy.get("@clerkToken").then((token) => {
      cy.task("getApiURL").then((t) => {
        cy.request({
          method: "GET",
          url: `${t}/notes`,
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }).then((resp) => {
          cy.get('[data-cy="note-list"] > ol')
            .as("noteList")
            .then((noteListEl) => {
              const timeEls: NodeListOf<HTMLTimeElement> =
                noteListEl[0].querySelectorAll(
                  '[data-cy="date-paragraph"] time',
                );
              const dateTimes: Array<number> = Array.from(timeEls).map(
                function (timeEl): number {
                  const timeStringParts = timeEl
                    .getAttribute("datetime")
                    .split(":");
                  return Number(
                    timeStringParts[timeStringParts.length - 1].split("Z")[0],
                  );
                },
              );
              const sorted = [...dateTimes].sort(function (a, b) {
                return b - a;
              });

              const apiResponse: Array<number> = (
                JSON.parse(resp.body) as Array<Note>
              ).map(function (respItem): number {
                const timeStringParts = respItem.createdAt.split(":");
                return Number(
                  timeStringParts[timeStringParts.length - 1].split("Z")[0],
                );
              });
              const tableData = dateTimes.map((dt, i) => ({
                api: apiResponse[i],
                dt,
                s: sorted[i],
              }));
              console.table(tableData);

              expect(apiResponse).to.deep.equal(dateTimes);
              expect(apiResponse).to.deep.equal(sorted);
              expect(dateTimes).to.deep.equal(sorted);
            });
        });
      });
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
    cy.get('[data-cy="content-from-list"]').first().contains(content);
  });

  it("should see specific note page and update note", () => {
    cy.visit("/notes");
    cy.get('[data-cy="content-from-list"]')
      .first()
      .as("latestNote")
      .then((el) => {
        return el[0].innerText;
      })
      .then((elText) => {
        cy.get('[data-cy="note-page-link"]').first().click();
        cy.get('[data-cy="content-specific-note"]').contains(elText);
        cy.get('[data-cy="update-link"]').click();
        cy.get('[data-cy="update-form-content-field"]')
          .as("contentField")
          .then((contentFieldEl) => {
            expect((contentFieldEl[0] as HTMLInputElement).value).to.equal(
              elText,
            );
          });
        const newContent = "betlog";
        cy.get("@contentField").click().clear().type(newContent);
        cy.get('[data-cy="update-form-submit"]').click();
        cy.get('[data-cy="content-specific-note"]').contains(newContent);

        cy.visit("/notes");
        cy.get("@latestNote").contains(newContent);
      });
  });

  it("should delete note", () => {
    cy.visit("/notes");
    cy.get('[data-cy="delete-link"]')
      .as("deleteLinks")
      .should("have.length", 5);
    cy.get("@deleteLinks").first().click();
    cy.get("@deleteLinks").should("have.length", 4);
  });
});
