import { NextAuthOptions } from 'next-auth';
import CustomRestAdapter from "~/auth/custom-rest-adapter";
import { env } from '~/env.mjs';

export const authOptions = {
  pages: {
      // signIn: '/login'
  },
  debug: true,
  adapter: CustomRestAdapter({
    baseUrl: env.BACKEND_BASE_URL,
    apiSecret: env.NEXTAUTH_SECRET,
  })
} satisfies Omit<NextAuthOptions, 'providers'>
