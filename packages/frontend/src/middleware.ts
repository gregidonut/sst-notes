import { defineMiddleware } from "astro:middleware";
import { clerkMiddleware, createRouteMatcher } from "@clerk/astro/server";
import { sequence } from "astro:middleware";

const validation = defineMiddleware(async function validation(_, next) {
    console.log("validation request");
    const response = await next();
    console.log("validation response");
    return response;
});
const isPublicRoute = createRouteMatcher(["/sign-in(.*)", "/sign-up(.*)"]);

const greeting = defineMiddleware(async function greeting(_, next) {
    console.log("greeting request");
    const response = await next();
    console.log("greeting response");
    return response;
});

export const onRequest = sequence(
    validation,
    clerkMiddleware(
        // @ts-ignore
        function (auth, context) {
            const { userId } = auth();

            if (!isPublicRoute(context.request) && !userId) {
                return context.redirect("/sign-in");
            }
        },
    ),
    greeting,
);
