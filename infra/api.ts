import { table } from "./storage";

const ClerkJWTAuthorizerAud = new sst.Secret("ClerkJWTAuthorizer");
const ClerkDevAccountDefaultEndpoint = new sst.Secret(
  "ClerkDevAccountDefaultEndpoint",
);
// Create the API
export const api = new sst.aws.ApiGatewayV2("Api", {
  link: [ClerkJWTAuthorizerAud, ClerkDevAccountDefaultEndpoint],
  transform: {
    route: {
      handler: {
        link: [table],
      },
    },
  },
  cors: true,
});

const ClerkJWTAuthorizer = "ClerkJWTAuthorizer";

const ClerkJWTAuthorizerVar = api.addAuthorizer({
  name: ClerkJWTAuthorizer,
  jwt: {
    issuer: ClerkDevAccountDefaultEndpoint.value,
    audiences: [ClerkJWTAuthorizerAud.value],
  },
});
// addProtectedRoute("POST /notes", "packages/functions/src/create.main");
// addProtectedRoute("GET /notes/{id}", "packages/functions/src/get.main");
// addProtectedRoute("GET /notes", "packages/functions/src/list.main");
// addProtectedRoute("DELETE /notes/{id}", "packages/functions/src/delete.main");
addProtectedRoute("PUT /notes/{id}", "packages/functions/src/update.main");
addProtectedGoRoute("GET /notes", "packages/functions/cmd/list/main.go");
addProtectedGoRoute("GET /notes/{id}", "packages/functions/cmd/get/main.go");
addProtectedGoRoute("POST /notes", "packages/functions/cmd/create/main.go");
addProtectedGoRoute(
  "DELETE /notes/{id}",
  "packages/functions/cmd/delete/main.go",
);

function addProtectedRoute(rawRoute: string, handler: string): void {
  api.route(rawRoute, handler, {
    auth: {
      jwt: {
        authorizer: ClerkJWTAuthorizerVar.id,
      },
    },
  });
}

function addProtectedGoRoute(rawRoute: string, handler: string): void {
  api.route(
    rawRoute,
    {
      handler,
      runtime: "go",
      environment: {
        NOTES_TABLE_NAME: table.name,
      },
    },
    {
      auth: {
        jwt: {
          authorizer: ClerkJWTAuthorizerVar.id,
        },
      },
    },
  );
}
