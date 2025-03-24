import { api } from "./api";
import { bucket } from "./storage";

const region = aws.getRegionOutput().name;

export const frontend = new sst.aws.Astro("Frontend", {
  path: "packages/frontend",
  environment: {
    ASTRO_REGION: region,
    ASTRO_API_URL: api.url,
    ASTRO_BUCKET: bucket.name,
    // VITE_USER_POOL_ID: userPool.id,
    // VITE_IDENTITY_POOL_ID: identityPool.id,
    // VITE_USER_POOL_CLIENT_ID: userPoolClient.id,
  },
});
