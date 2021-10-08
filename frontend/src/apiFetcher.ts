export const contentTypes = {
  json: 'application/json',
  toml: 'application/toml',
  yaml: 'text/vnd.yaml',
} as const

export const methods = {
  POST: 'POST',
  GET: 'GET',
} as const

export const baseUrl = 'http://localhost/api'

export type ApiFetchOptions = {
  /** Body, as json. Can either be stringified or an object, in which case it will be stringified */
  body?: {} | string
  /** HTTP-method */
  method?: string
  /** Allows to use JMES-path on the result. This will be handled by the api.
   *
   * NOTE: if set, the updater-function will not run by default.
   */
  jmespath?: string
}

export async function fetchApi<T extends {}>(
  subPath: string,
  updater: (data: T) => void,
  { method = methods.GET, body, jmespath }: ApiFetchOptions = {}
) {
  const sub = subPath.replace(/^\/?/, '/').replace(/\/?$/, '/')
  const opts: RequestInit = {
    method,
    headers: {
      accept: contentTypes.json,
      'content-type': contentTypes.json,
    },
    ...(!!body && {
      body: typeof body === 'string' ? body : JSON.stringify(body),
    }),
  }
  if (jmespath) {
    opts.headers['jmes-path'] = jmespath
  }
  const url = baseUrl + sub
  const result: {
    data: T
  } = {} as any
  let response: Response
  try {
    response = await fetch(url, opts)
    const contentType = response.headers.get('content-type') || ''
    if (contentType.includes(contentTypes.json)) {
      const JSON = await response.json()
      if (response.status >= 400) {
        return [result, JSON as ApiResponses.ApiError] as const
      }
      result.data = JSON
      if (JSON && !jmespath) {
        !!JSON && updater(JSON)
      }
    }
  } catch (err) {
    console.error(`fetchApi error for ${subPath}: ${err.message}`, {
      subPath,
      url,
      opts,
      err,
      response,
    })
    return [
      result,
      {
        error: err.message as string,
        originalError: err,
        code: response?.status || 'NoStatusReceived',
      } as ApiResponses.ApiError & { originalError: Error },
    ] as const
  }
  return [result, null] as const
}
