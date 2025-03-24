// @ts-check
import { defineConfig } from "astro/config";

import compressor from "astro-compressor";
import aws from "astro-sst";

// https://astro.build/config
export default defineConfig({
    integrations: [compressor()],
    output: "server",
    adapter: aws({
        responseMode: "stream",
    }),
});

