import { ApiFetchOptions, fetchApi, methods } from './apiFetcher'
import createStore from './store'
import { objectKeys } from 'simplytyped'

/**
 * Typedefintions are created by running `yarn gen`.
 *
 * This will use the generated swagger defintions from Gobyoall-api. (which again are created from running go generate)
 */

export type DB = {
  endpoint: Record<string, ApiDef.EndpointEntity>
  request: Record<string, ApiDef.RequestEntity>
  schedule: Record<string, ApiDef.ScheduleEntity>
  stat: Record<string, ApiDef.CompactRequestStatisticsEntity>
  serverInfo: ApiDef.ServerInfo
}

export const api = {
  serverInfo: (options?: ApiFetchOptions) =>
    fetchApi<ApiDef.ServerInfo>(
      'serverInfo',
      (e) => db.update((s) => ({ ...s, serverInfo: e })),
      options
    ),
  request: CrudFactory<ApiDef.RequestPayload, 'request'>('request'),
  endpoint: CrudFactory<ApiDef.EndpointPayload, 'endpoint'>('endpoint'),
  stat: {
    list: apiGetListFactory<'stat'>('stat', 'stat'),
  },
  schedule: {
    ...CrudFactory<ApiDef.SchedulePayload, 'schedule'>('schedule'),
    update: apiUpdateFactory<ApiDef.SchedulePayload, 'schedule'>(
      'schedule',
      'schedule'
    ),
  },
} as const

export const db = createStore<DB, null>({
  initialValue: objectKeys(api).reduce((r, k) => {
    r[k] = {}
    return r
  }, {} as DB),
})

const mergeMap = <K extends DBKeyValue, V extends DB[K]>(key: K, value: V) => {
  if (!key) {
    console.error('key is required in mergeField')
    return
  }
  if (!value) {
    console.error('value is required in mergeField')
    return
  }
  db.update((s) => {
    return {
      ...s,
      [key]: {
        ...s[key],
        ...value,
      },
    }
  })
}

// Keys in in that are of type Record<string, T>
type DBKeyValue = keyof Omit<DB, 'serverInfo'>

const mergeField = <K extends DBKeyValue, V extends DB[K]['s']>(
  key: K,
  value: V,
  id: string
) => {
  if (!key) {
    console.error('key is required in mergeField')
    return
  }
  if (!value) {
    console.error('value is required in mergeField')
    return
  }
  if (!id) {
    console.error('id is required in mergeField')
    return
  }
  db.update((s) => {
    return {
      ...s,
      [key]: {
        ...s[key],
        [id]: {
          ...s[key]?.[id],
          ...value,
        },
      },
    }
  })
}

/* 
  Returns typed functions for:
  - Create
  - Get
  - List

  yes, that is not really all of the cruds...
*/
function CrudFactory<Payload extends {}, K extends DBKeyValue>(
  storeKey: K,
  subPath?: string
) {
  return {
    get: apiGetFactory(subPath || storeKey, storeKey),
    list: apiGetListFactory(subPath || storeKey, storeKey),
    create: apiCreateFactory<Payload, K>(subPath || storeKey, storeKey),
  }
}

function apiGetListFactory<K extends DBKeyValue>(subPath: string, storeKey: K) {
  return async (options?: ApiFetchOptions) => {
    const res = await fetchApi<DB[K]>(
      subPath,
      (e) => mergeMap(storeKey, e),
      options
    )
    if (!res[1]) {
      if (res[0].data) {
        db.update((s) => ({
          ...s,
          [storeKey]: { ...s[storeKey], ...res[0].data },
        }))
      }
    }
    return res
  }
}

function apiGetFactory<K extends DBKeyValue>(subPath: string, storeKey: K) {
  return (id: string, options?: ApiFetchOptions) =>
    fetchApi<DB[K]>(
      subPath + id,
      (e: any) => mergeField(storeKey, e, e.id),
      options
    )
}
function apiCreateFactory<Payload extends {}, K extends DBKeyValue>(
  subPath: string,
  storeKey: K
) {
  return (body: Payload, options?: ApiFetchOptions) =>
    fetchApi<DB[K]['s']>(subPath, (e) => mergeField(storeKey, e, e.id), {
      method: methods.POST,
      body,
      ...options,
    })
}

function apiUpdateFactory<Payload extends {}, K extends DBKeyValue>(
  subPath: string,
  storeKey: K
) {
  if (!subPath) {
    subPath = storeKey
  }
  return (id: string, body: Payload, options?: ApiFetchOptions) =>
    fetchApi<DB[K]['s']>(
      subPath + '/' + id,
      (e) => mergeField(storeKey, e, e.id),
      {
        method: methods.PUT,
        body,
        ...options,
      }
    )
}
