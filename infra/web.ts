import { api } from "./api";
import { bucket } from "./storage";

const region = aws.getRegionOutput().name;

const clerkPublic = new sst.Secret("ClerkPublicKey");
const clerkSecret = new sst.Secret("ClerkSecretKey");
const stage = process.env.SST_STAGE;
export const frontend = new sst.aws.Astro("Frontend", {
  path: "packages/frontend",
  link: [clerkPublic, clerkSecret],
  environment: {
    ASTRO_STAGE: stage,
    ASTRO_REGION: region,
    ASTRO_API_URL: api.url,
    ASTRO_BUCKET: bucket.name,
    PUBLIC_CLERK_PUBLISHABLE_KEY: clerkPublic.value,
    CLERK_SECRET_KEY: clerkSecret.value,
    CLERK_SIGN_IN_URL: "/sign-in",
  },
});
