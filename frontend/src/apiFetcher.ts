export const contentTypes = {
  json: 'application/json',
  toml: 'application/toml',
  yaml: 'text/vnd.yaml',
} as const

export const methods = {
  POST: 'POST',
  GET: 'GET',
  PUT: 'PUT',
  DELETE: 'DELETE',
} as const

export const baseUrl = `${window.location.protocol}//${window.location.host}/api`
export const wsUrl = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws/`

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
      ...(!!jmespath && {
        'jmes-path': jmespath,
      }),
    },
    ...(!!body && {
      body: typeof body === 'string' ? body : JSON.stringify(body),
    }),
  }
  const url = baseUrl + sub
  const result: {
    data: T
  } = {} as any
  let response: Response | null = null
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

export function serializeDate(date: Date) {
  return date.toISOString()
}
export function deserializeDate(dateStr: string) {
  return new Date(dateStr)
}

let wsDisconnects = 0
let wsFails = 0
export const wsSubscribe = (options: {
  onMessage: (msg: WsMessage) => void,
  onClose?: () => void,
  autoReconnect: boolean
}) => {
  const {
    onMessage,
    onClose,
    autoReconnect,
  } = options
  if (!window["WebSocket"]) {
    console.error("Your browser does not support WebSocket")
    return
  }
  try {

    const conn = new WebSocket(wsUrl);
    conn.onclose = function(evt) {
      console.log('[ws]: connection closed')
      onClose?.()
      wsDisconnects++
      if (autoReconnect) {
        setTimeout(
          () => wsSubscribe(options), 1000 * (wsDisconnects + wsFails)
        )
      }

    };
    conn.onmessage = function(evt) {
      console.log('[ws]: got message', evt.data)
      try {
        const json = JSON.parse(evt.data)
        onMessage(json)
      } catch (err) {
        console.error('Failed to parse json-message\n', err)
      }
    };
  } catch (err) {
    console.error('Failed in wsSubscribe ', err)
    wsFails++

  }
}

type Ws<K extends string, V extends string, T> = {
  kind: K
  variant: V
  contents: T
}


type WsCreateEndpoit = Ws<'endpoint', 'create', ApiDef.EndpointEntity>
type WsUpdateEndpoit = Ws<'endpoint', 'update', ApiDef.EndpointEntity>
type WsDeleteEndpoit = Ws<'endpoint', 'soft-delete', ApiDef.EndpointEntity>

type WsCreateRequest = Ws<'request', 'create', ApiDef.RequestEntity>
type WRequestRequest = Ws<'request', 'update', ApiDef.RequestEntity>
type WsDeleteRequest = Ws<'request', 'soft-delete', ApiDef.RequestEntity>

type WsCreateSchedule = Ws<'schedule', 'create', ApiDef.StatEntity>
type WRequestSchedule = Ws<'schedule', 'update', ApiDef.StatEntity>
type WsDeleteSchedule = Ws<'schedule', 'soft-delete', ApiDef.StatEntity>

type WsCreateStat = Ws<'stat', 'create', ApiDef.StatEntity>
type WRequestStat = Ws<'stat', 'update', ApiDef.StatEntity>
type WsClearAllStats = Ws<'stat', 'clean', ApiDef.StatEntity>

type WsMessage =
  | WsCreateEndpoit
  | WsUpdateEndpoit
  | WsDeleteEndpoit

  | WsCreateRequest
  | WRequestRequest
  | WsDeleteRequest

  | WsCreateSchedule
  | WRequestSchedule
  | WsDeleteSchedule

  | WsCreateStat
  | WRequestStat
  | WsClearAllStats
