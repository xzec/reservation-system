import { resolve } from "./resolve"

const resolveFetch = (promise: Promise<Response>) =>
  resolve(
    promise.then((response) => {
      if (response.ok) return response
      throw new Error(`${response.status} ${response.statusText}`)
    })
  )

export default resolveFetch
