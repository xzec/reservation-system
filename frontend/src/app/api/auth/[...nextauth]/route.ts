import CognitoProvider from "next-auth/providers/cognito";
import { env } from "~/env.mjs";
import NextAuth from "next-auth";
import { authOptions } from "~/auth/auth";

const handler = NextAuth({
  ...authOptions,
  providers: [
    CognitoProvider({
      clientId: env.COGNITO_CLIENT_ID,
      clientSecret: env.COGNITO_CLIENT_SECRET,
      issuer: env.COGNITO_ISSUER,
    }),
  ]});

export { handler as GET, handler as POST }
