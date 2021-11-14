import createStore, { isRedacted, Store } from 'store'
import type { DeepNonNullable, DeepRequired } from 'utility-types'
import type { DeepNullable } from 'ts-essentials'
import { AnyFunc, objectKeys } from 'simplytyped'
import type { Paths } from 'types'

// TODO: document why DeepNullable and DeepRequired is needed here or remove them.
type ConfigStoreType = DeepNullable<DeepRequired<ApiDef.Config>> & {
  _: {
    showDecodedJWT: boolean
    decodedToken: Array<[string, string | null]>
    impersonateUserName: boolean
  }
}
type PathKeys =
  | Paths<ApiDef.Config>
  | 'auth.impersionation_credentials.userID/userName'

const configStore = (() => {
  const s = createStore<ConfigStoreType, ApiDef.Config, PathKeys>({
    validator: (c) => configStoreToConfigPayload(c),
    initialValue: {
      _: {
        showDecodedJWT: false,
        decodedToken: [],
        impersonateUserName: false,
      },
      concurrency: null,
      request_count: null,
      response_data: null,
      ok_status_codes: [],
      secrets: {},
      auth: {
        kind: null,
        redirect_uri: null,
        client_id: null,
        client_secret: null,
        dynamic: { requests: [], headerKey: '' },
        endpoint: null,
        token: null,
        endpoint_type: null,
        header_key: null,
        impersionation_credentials: {
          user_id_to_impersonate: null,
          user_name_to_impersonate: null,
          username: null,
          password: null,
        },
      },
    },
    storage: { key: 'form-config' },
  })

  return {
    ...s,
    restore: (config?: ApiDef.Config) => {
      if (!config) {
        return s.restore()
      }
      s.restore({
        ...(config as any),
        _: {
          showDecodedJWT: false,
          decodedToken: [],
          impersonateUserName: false,
        },
      })
    },
  }
})()

/**
 * @deprecated bad pattern
 *
 * Has bugs with updating...
 *
 * */
export const configStoreToConfigPayload = (
  cfg: Store<ConfigStoreType, ApiDef.Config>
):
  | [ApiDef.Config, null]
  | [null, Partial<Record<PathKeys, string>>]
  | [null, null] => {
  if (!cfg.__didChange) {
    return [null, null]
  }

  // Pay attention to nullables
  const c: ApiDef.Config = {}
  if (cfg.concurrency) {
    if (cfg.concurrency < 0) {
      return [null, { concurrency: 'Concurrency must be positive' }]
    }
    if (!cfg.request_count) {
      return [
        null,
        {
          request_count:
            'If concurrency is set, request-count must also be set',
        },
      ]
    }
    if (cfg.concurrency > cfg.request_count) {
      return [
        null,
        {
          concurrency: 'Concurrency cannot be higher than Request-count',
        },
      ]
    }
    c.concurrency = cfg.concurrency
  }
  if (cfg.request_count !== null) {
    if (cfg.request_count < 0) {
      return [null, { request_count: '!!!Request-count must be positive' }]
    }
    if (!cfg.concurrency) {
      return [
        null,
        {
          concurrency: 'If Request-count is set, concurrency must also be set',
        },
      ]
    }
    if (cfg.request_count < cfg.concurrency) {
      return [
        null,
        { request_count: 'Concurrency cannot be higher than request-count' },
      ]
    }
    c.request_count = cfg.request_count
  }
  if (typeof cfg.response_data === 'boolean') {
    c.response_data = cfg.response_data
  }
  const a: ApiDef.Config['auth'] = {}
  const { auth } = cfg

  if (auth) {
    const { impersionation_credentials, dynamic, ...rest } = auth
    switch (auth.kind) {
      case '':
        // Does this mean to remove all props??
        break
      case 'bearer':
        if (auth.token?.trim() === '') {
          return [null, { 'auth.token': 'Token must be set if type is bearer' }]
        }
        isNeitherNullOrRadacted(rest.token, (t) => (a.token = t))
        break
      case 'impersonation':
        if (!auth.client_id) {
          return [null, { 'auth.client_id': 'client-id is missing' }]
        }
        if (!auth.endpoint_type) {
          return [
            null,
            {
              'auth.endpoint_type': 'endpoint-type must be one of: "keycloak"',
            },
          ]
        }
        if (!auth.redirect_uri) {
          return [null, { 'auth.redirect_uri': 'redirect-uri must be set' }]
        }
        if (!auth.impersionation_credentials) {
          return [
            null,
            {
              'auth.impersionation_credentials':
                'Impersonation: credentials missing',
            },
          ]
        }
        if (!auth.impersionation_credentials.username) {
          return [
            null,
            {
              'auth.impersionation_credentials':
                'Impersonation: username missing',
            },
          ]
        }
        if (
          !auth.impersionation_credentials.user_id_to_impersonate &&
          !auth.impersionation_credentials.user_name_to_impersonate
        ) {
          return [
            null,
            {
              'auth.impersionation_credentials.userID/userName':
                'Impersonation: No id or username to impersonate is set',
            },
          ]
        }
        if (!auth.impersionation_credentials) {
          auth.impersionation_credentials
        }
        for (const key of objectKeys(impersionation_credentials)) {
          isNeitherNullOrRadacted(impersionation_credentials[key], (v) => {
            auth.impersionation_credentials = {
              ...auth.impersionation_credentials,
              [key]: v,
            }
          })
        }
        for (const key of objectKeys(rest)) {
          isNeitherNullOrRadacted(rest[key], (v) => {
            auth[key] = v
          })
        }
        break
      case 'dynamic':
        if (dynamic.requests.length === 0) {
          return [
            null,
            {
              'auth.dynamic':
                'For dynamic-authentication, at least one request must be appended',
            },
          ]
        }

        for (const r of dynamic.requests) {
          if (!r.uri) {
            return [null, { 'auth.dynamic': 'Uri must be set' }]
          }
          if (!r.method) {
            return [null, { 'auth.dynamic': 'Method must be set' }]
          }
        }
        a.header_key = dynamic.headerKey || ''
        a.dynamic = {
          requests: dynamic.requests.map(({ headers, ...r }) => {
            return {
              ...r,
              json_request: !!r.json_request,
              json_response: !!r.json_response,
              method: r.method || '',
              uri: r.uri || '',
              result_jmes_path: r.result_jmes_path || '',
            }
          }),
        }
        break

      default:
        break
    }
  }
  if (!!Object.keys(a).length) {
    c.auth = a
    c.auth.kind = auth.kind!
  }
  if (!Object.keys(c).length) {
    return [null, null]
  }
  return [c, null]
}

const isNeitherNullOrRadacted = <T>(v: T, f?: (t: NonNullable<T>) => void) => {
  if (v === undefined || v === null || isRedacted(v as any)) {
    return
  }
  f?.(v as any)
  return v
}

export default configStore
