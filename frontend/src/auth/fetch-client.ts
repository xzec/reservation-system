import { z } from "zod"
import resolveFetch from "~/utils/resolveFetch"
import { resolve } from "~/utils/resolve"

export const createFetchClient = (baseUrl: string, apiSecret: string) => {
  return async function client<T>(method: 'POST' | 'GET' | 'PUT' | 'DELETE', path: string, body?: any) {
    const url = `${baseUrl}/auth${path}`
    const headers = {
      'Content-Type': 'application/json',
      'X-Api-Key': apiSecret,
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

    return isObject(jRes) ? formatRes<T>(jRes) : jRes
  }
}

const isObject = (value: unknown): value is Record<string, unknown> => typeof value === 'object' && !Array.isArray(value) && value !== null

const formatRes = <T, >(res: Record<string, unknown>): T => {
  const recursiveFormatRes = (obj: Record<string, unknown>) => {
    return Object.entries(obj).reduce<Record<string, unknown>>((result, [key, value]) => {
      if (isObject(value))
        result[key] = recursiveFormatRes(value)
      else {
        // attempt to convert every value to Date and replace it when succeeded
        const dateValue = z.string().datetime({ offset: true }).transform((value) => new Date(value)).safeParse(value)
        result[key] = dateValue.success ? dateValue.data : value
      }
      return result
    }, {}) as T
  }

  return recursiveFormatRes(res)
}