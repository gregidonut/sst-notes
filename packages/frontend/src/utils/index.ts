export const isDevStage = import.meta.env.ASTRO_STAGE === "dev";

export function cy(attr: string) {
    return isDevStage ? { "data-cy": attr } : {};
}

export interface User {
    username: string;
}
