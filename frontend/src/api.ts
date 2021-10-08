import { ApiFetchOptions, fetchApi, methods } from './apiFetcher'
import createStore from './store'

/**
 * Typedefintions are created by running `yarn gen`.
 *
 * This will use the generated swagger defintions from Gobyoall-api. (which again are created from running go generate)
 */

export type DB = {
  endpoint: Record<string, ApiDef.EndpointEntity>
  request: Record<string, ApiDef.RequestEntity>
}

export const api = {
  request: CrudFactory<ApiDef.RequestPayload, 'request'>('request', 'request'),
  endpoint: CrudFactory<ApiDef.EndpointPayload, 'endpoint'>(
    'endpoint',
    'endpoint'
  ),
}

const withStorage = !!false

export const db = createStore<DB>({
  initialValue: {
    endpoint: {},
    request: {},
  },
  ...(withStorage && {
    storage: {
      key: 'store',
    },
  }),
})

const mergeMap = <K extends keyof DB, V extends DB[K]>(key: K, value: V) => {
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
const mergeField = <K extends keyof DB, V extends DB[K]['s']>(
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
function CrudFactory<Payload extends {}, K extends keyof DB>(
  subPath: string,
  storeKey: K
) {
  return {
    get: apiGetFactory(subPath, storeKey),
    list: apiGetListFactory(subPath, storeKey),
    create: apiCreateFactory<Payload, K>(subPath, storeKey),
  }
}

function apiGetListFactory<K extends keyof DB>(subPath: string, storeKey: K) {
  return async (options?: ApiFetchOptions) => {
    const res = await fetchApi<DB[K]>(
      subPath,
      (e: any) => mergeMap(storeKey, e),
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

function apiGetFactory<K extends keyof DB>(subPath: string, storeKey: K) {
  return (id: string, options?: ApiFetchOptions) =>
    fetchApi<DB[K]>(
      subPath + id,
      (e: any) => mergeField(storeKey, e, e.id),
      options
    )
}
function apiCreateFactory<Payload extends {}, K extends keyof DB>(
  subPath: string,
  storeKey: K
) {
  return (body: Payload, options?: ApiFetchOptions) =>
    fetchApi<DB[K]['s']>(subPath, (e: any) => mergeField(storeKey, e, e.id), {
      method: methods.POST,
      body,
      ...options,
    })
}
