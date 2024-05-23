"use client";

import type { FC, PropsWithChildren } from 'react';
import { SessionProvider } from "next-auth/react";

const NextAuthProvider: FC<PropsWithChildren> = ({ children }) => (
  <SessionProvider>
    {children}
  </SessionProvider>
)

export default NextAuthProvider
