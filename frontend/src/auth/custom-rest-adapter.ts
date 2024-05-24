import type { Adapter, AdapterAccount, AdapterSession, AdapterUser, VerificationToken } from "next-auth/adapters"
import { Awaitable } from "next-auth"
import resolveFetch from "~/utils/resolveFetch"
import { resolve } from "~/utils/resolve"

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
      return client<AdapterUser>('PUT', `/users/${user.id}`, user)
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

const createFetchClient = (baseUrl: string, apiSecret: string) => {
  return async function client<T>(method: 'POST' | 'GET' | 'PUT' | 'DELETE', path: string, body?: any) {
    const url = `${baseUrl}/auth${path}`
    const headers = {
      'Content-Type': 'application/json',
      'X-Api-Secret': apiSecret,
    }

    const config: RequestInit = {
      method,
      headers,
      body: body ? JSON.stringify(body) : null
    }

    const [err, res] = await resolveFetch(fetch(url, config))
    if (err) throw new Error(`Failed to fetch the Auth API [${method}][${url}]. Error: ${err}`)

    const [jErr, jRes] = await resolve<T>(res.json())
    if (jErr) throw new Error(`Failed to parse JSON from the Auth API response [${method}][${url}]. Error: ${jErr}`)

    return jRes
  }
}
