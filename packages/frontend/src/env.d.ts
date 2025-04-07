/// <reference types="astro/client" />

import type { User } from "@/utils";

declare global {
    namespace App {
        interface Locals {
            user: User;
        }
    }
}
