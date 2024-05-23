"use client";

import { signIn, signOut, useSession } from "next-auth/react";

export default function AuthButton() {
  const session = useSession()
  console.log('🔴 session', session)

  if (session.status === 'authenticated') {
    return (
      <>
        {session?.data?.user?.name} <br />
        <button onClick={() => signOut()}>Sign out 👋</button>
      </>
    )
  }

  if (session.status === 'loading') return '🔁 Loading...'

  return (
    <>
      Not signed in <br />
      <button onClick={() => signIn()}>👉 Sign in</button>
    </>
  )
}
