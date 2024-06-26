import type { Adapter, AdapterAccount, AdapterSession, AdapterUser, VerificationToken } from "next-auth/adapters"
import { Awaitable } from "next-auth"
import { createFetchClient } from "~/auth/fetch-client"

type Options = {
  baseUrl: string
  apiSecret: string
}

export default function CustomRestAdapter({ baseUrl, apiSecret }: Options): Adapter {
  const client = createFetchClient(baseUrl, apiSecret)

  return {
    createUser(user) {
      console.log('createUser', user)
      return client<AdapterUser>('POST', '/users', user)
    },
    getUser(id) {
      console.log('getUser', id)
      return client<AdapterUser | null>('GET', `/users/${id}`)
    },
    getUserByEmail(email) {
      console.log('getUserByEmail', email)
      return client<AdapterUser | null>('GET', `/users/email/${email}`)
    },
    getUserByAccount({ provider, providerAccountId }) {
      console.log('getUserByAccount', { provider, providerAccountId })
      return client<AdapterUser | null>('GET', `/users/account/${provider}/${providerAccountId}`)
    },
    updateUser(user) {
      console.log('updateUser', user)
      const { id, ...fields } = user
      return client<AdapterUser>('PUT', `/users/${id}`, fields)
    },
    deleteUser(userId: string) {
      console.log('deleteUser', userId)
      return client<AdapterUser | null | undefined>('DELETE', `/users/${userId}`)
    },
    linkAccount(account) {
      console.log('linkAccount', account)
      return client<AdapterAccount | null | undefined>('POST', '/accounts', account)
    },
    unlinkAccount({ provider, providerAccountId }) {
      console.log('unlinkAccount', { provider, providerAccountId })
      return client<AdapterAccount | undefined>('DELETE', `/accounts/${provider}/${providerAccountId}`)
    },
    createSession(session) {
      console.log('createSession', session)
      return client<AdapterSession>('POST', '/sessions', session)
    },
    getSessionAndUser(sessionToken) {
      console.log('getSessionAndUser', sessionToken)
      return client<{ session: AdapterSession; user: AdapterUser } | null>('GET', `/sessions/${sessionToken}`)
    },
    updateSession({ sessionToken, userId, expires }) {
      console.log('updateSession', { sessionToken, userId, expires  })
      return client<AdapterSession | null | undefined>('PUT', `/sessions/${sessionToken}`, { userId, expires })
    },
    deleteSession(sessionToken) {
      console.log('deleteSession', sessionToken)
      return client<AdapterSession | null | undefined>('DELETE', `/sessions/${sessionToken}`)
    },
    createVerificationToken(verificationToken) {
      console.log('createVerificationToken', verificationToken)
      return client<VerificationToken | null | undefined>('POST', '/verification-tokens', verificationToken)
    },
    useVerificationToken(params): Awaitable<VerificationToken | null> {
      console.log('useVerificationToken', params)
      return client<VerificationToken | null>('POST', '/verification-tokens/use', params)
    }
  }
}

