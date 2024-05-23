import type {Adapter, AdapterAccount, AdapterSession, AdapterUser, VerificationToken} from "next-auth/adapters"
import {Awaitable} from "next-auth"
import {undefined, util} from "zod"
import resolveFetch from "~/utils/resolveFetch"
import Omit = util.Omit;

type Options = {
  baseUrl: string
  apiSecret: string
}

export default function CustomRestAdapter({ baseUrl, apiSecret }: Options): Adapter {
  const fetcher = createFetcher(baseUrl, apiSecret)

  return {
    async createUser(user: Omit<AdapterUser, "id">) {
      console.log('createUser', user)
      const res = await fetcher<AdapterUser>('POST', '/user', user)
      return res
    },
    getUser(id: string): Awaitable<AdapterUser | null> {
      console.log('getUser', id)
      return undefined
    },
    getUserByEmail(email: string): Awaitable<AdapterUser | null> {
      console.log('getUserByEmail', email)
      return undefined
    },
    getUserByAccount(providerAccountId: Pick<AdapterAccount, "provider" | "providerAccountId">): Awaitable<AdapterUser | null> {
      console.log('getUserByAccount', providerAccountId)
      return undefined
    },
    updateUser(user: Partial<AdapterUser> & Pick<AdapterUser, "id">): Awaitable<AdapterUser> {
      console.log('updateUser', user)
      return undefined
    },
    deleteUser(userId: string): Promise<void> | Awaitable<AdapterUser | null | undefined> {
      console.log('deleteUser', userId)
      return undefined
    },
    linkAccount(account: AdapterAccount): Promise<void> | Awaitable<AdapterAccount | null | undefined> {
      console.log('linkAccount', account)
      return undefined
    },
    unlinkAccount(providerAccountId: Pick<AdapterAccount, "provider" | "providerAccountId">): Promise<void> | Awaitable<AdapterAccount | undefined> {
      console.log('unlinkAccount', providerAccountId)
      return undefined
    },
    async createSession(session: { sessionToken: string; userId: string; expires: Date }) {
      console.log('createSession', session)
      const res = await fetcher<AdapterSession>('POST', '/user', session)
      return res
    },
    getSessionAndUser(sessionToken: string): Awaitable<{ session: AdapterSession; user: AdapterUser } | null> {
      console.log('getSessionAndUser', sessionToken)
      return undefined
    },
    updateSession(session: Partial<AdapterSession> & Pick<AdapterSession, "sessionToken">): Awaitable<AdapterSession | null | undefined> {
      console.log('updateSession', session)
      return undefined
    },
    deleteSession(sessionToken: string): Promise<void> | Awaitable<AdapterSession | null | undefined> {
      console.log('deleteSession', sessionToken)
      return undefined
    },
    createVerificationToken(verificationToken: VerificationToken): Awaitable<VerificationToken | null | undefined> {
      console.log('createVerificationToken', verificationToken)
      return undefined
    },
    useVerificationToken(params: { identifier: string; token: string }): Awaitable<VerificationToken | null> {
      console.log('useVerificationToken', params)
      return undefined
    }
  }
}

const createFetcher = (baseUrl: string, apiSecret: string) => {
  return async <T>(method: string, endpoint: string, body: any) => {
    const url = `${baseUrl}/auth${endpoint}`
    const headers = {
      'Content-Type': 'application/json',
      'x-api-secret': apiSecret,
    }

    const config: RequestInit = {
      method,
      headers,
      body: body ? JSON.stringify(body) : null
    }

    const [err, res] = await resolveFetch(fetch(url, config))
    if (err) throw new Error(`Response was not ok. Error: ${err}`)

    return res.json() as T
  }
}
