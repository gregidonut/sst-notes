import { clerkSetup } from "@clerk/testing/cypress";
import { defineConfig } from "cypress";
import seedDynamoDB from "./cypress/tasks/seedDyanamoDB";
import emptyDynamoDB from "./cypress/tasks/emptyDynamoDB";
import { readFileSync } from "fs";
const sstOutputs = JSON.parse(readFileSync("./.sst/outputs.json", "utf-8"));

export default defineConfig({
  e2e: {
    setupNodeEvents(on, config) {
      on("task", {
        seedDyanamoDB(token: string) {
          return seedDynamoDB(sstOutputs["ApiUrl"], token);
        },

        emptyDynamoDB(token: string) {
          return emptyDynamoDB(sstOutputs["ApiUrl"], token);
        },
      });
      return clerkSetup({ config });
    },
    baseUrl: "http://localhost:4321",
  },
});
