import { NextAuthOptions } from 'next-auth';

export const authOptions = {
  providers: [],
  pages: {
      // signIn: '/login'
  },
  debug: true,

} satisfies NextAuthOptions

