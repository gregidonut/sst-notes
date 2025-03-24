// @ts-check
import { defineConfig } from "astro/config";

import compressor from "astro-compressor";
import aws from "astro-sst";

import clerk from "@clerk/astro";

import { shadesOfPurple } from "@clerk/themes";
// https://astro.build/config
export default defineConfig({
    integrations: [
        compressor(),
        clerk({
            appearance: {
                baseTheme: shadesOfPurple,
            },
        }),
    ],
    output: "server",
    adapter: aws({
        responseMode: "stream",
    }),
});
