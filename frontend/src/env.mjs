import { createEnv } from "@t3-oss/env-nextjs";
import { z } from "zod";

export const env = createEnv({
    server: {
        NEXTAUTH_SECRET: z.string().min(1),
        NEXTAUTH_URL: z.string().url(),
        COGNITO_CLIENT_ID: z.string().min(1),
        COGNITO_CLIENT_SECRET: z.string().min(1),
        COGNITO_ISSUER: z.string().url(),
        BACKEND_BASE_URL: z.string().url(),
    },
    client: {
    },
    runtimeEnv: {
        NEXTAUTH_SECRET: process.env.NEXTAUTH_SECRET,
        NEXTAUTH_URL: process.env.NEXTAUTH_URL,
        COGNITO_CLIENT_ID: process.env.COGNITO_CLIENT_ID,
        COGNITO_CLIENT_SECRET: process.env.COGNITO_CLIENT_SECRET,
        COGNITO_ISSUER: process.env.COGNITO_ISSUER,
        BACKEND_BASE_URL: process.env.BACKEND_BASE_URL,
    },
});
