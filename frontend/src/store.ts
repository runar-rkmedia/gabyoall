import { Updater, writable } from 'svelte/store'
import merge from 'lodash.merge'
import debounce from 'lodash.debounce'
import cloneDeep from 'lodash.clonedeep'
import deepEqual from 'lodash.isequal'

const parseStrOrObject = <T extends {}>(s: string | T | null | undefined) => {
  if (!s) {
    return [null, null] as const
  }
  if (typeof s !== 'string') {
    return [s as T, null] as const
  }
  try {
    const j = JSON.parse(s) as T
    return [j, null] as const
  } catch (err) {
    return [null, `failed to parse string into object '${s}': ${err}'`] as const
  }
}

type Redaction = RegExp | string | ((k: string, v: unknown) => boolean)
export const redactionString = '**REDACTED**' as const

export function isRedacted(s: string, redactionStr: string = redactionString) {
  return s === redactionStr
}

function reduceNullable<T>(obj: T) {
  const o = cloneDeep(obj)
  return _reduceNullable(obj)
}

function isNull(v: any) {
  return v === null
}
function _reduceNullable<T>(obj: T) {
  if (Array.isArray(obj)) {
    for (let i = obj.length - 1; i >= 0; i--) {
      if (isNull(obj)) {
        obj.splice(i, 1)
        continue
      }
      if (typeof obj[i] === 'object') {
        obj[i] = _reduceNullable(obj[i])
        continue
      }
    }
  }
  return
}

const redactor: typeof _redacter = <T extends {}>(
  obj,
  redaction,
  replacement
) => {
  return _redacter<T>(cloneDeep(obj), redaction, replacement)
}
const _redacter = <T extends {}>(
  obj: T,
  redaction: Array<Redaction>,
  replacement: string = redactionString
) => {
  if (!redaction || !redaction.length) {
    return obj
  }
  if (typeof obj !== 'object') {
    return obj
  }
  for (const key of Object.keys(obj)) {
    if (!obj[key] || obj[key] === replacement) {
      continue
    }
    if (typeof obj[key] === 'string') {
      for (const r of redaction) {
        switch (typeof r) {
          case 'object':
            if ('test' in r) {
              if (!r.test(key)) {
                continue
              }
            } else {
              throw new Error('received an unrecognized type of redaction')
            }
            break
          case 'function':
            if (!r(key, obj[key])) {
              continue
            }
            break
          case 'string':
            if (key !== r) {
              continue
            }
            break

          default:
            break
        }
        obj[key] = replacement
      }
      continue
    }
    if (typeof obj[key] === 'object') {
      obj[key] = _redacter(obj[key], redaction)
    }
  }
  return obj
}

export type StoreState<V = null, VK extends string = string> = {
  __didChange?: boolean
  __validationMessage?: Partial<Record<VK, string>>
  __validationPayload?: V
}

export type Store<T extends {}, V = null, VK extends string = string> = T &
  StoreState

/* 
TODO: use a form-library instead. 
  This turned out too complex. 
  Keep this store simple, and only use it only for the api, and perhaps some state.
 */
function createStore<T extends {}, V = null, VK extends string = string>({
  storage: _storage,
  initialValue,
  validator,
}: {
  storage?: (AppStorage<T> | { key: string }) & {
    redactKeys?: Array<Redaction>
  }
  validator?: (
    t: Store<T, V>
  ) => [null, null] | [V, null] | [null, Partial<Record<VK, string>>]
  initialValue?: T
} = {}) {
  type S = Store<T, V, VK>
  let fromStorageValue: T | null = null
  let restoreValue = initialValue
  const storage: AppStorage<T> | null = _storage?.key
    ? {
        getItem: (key) => localStorage.getItem(key),
        // TODO: throttle saving
        setItem: (k, v) => localStorage.setItem(k, JSON.stringify(v)),
        ..._storage,
      }
    : null

  if (storage) {
    if (_storage && _storage?.redactKeys === undefined) {
      _storage.redactKeys = [/secret/i, /password/i]
    }
    const str = storage.getItem(storage.key)
    const [parsed, err] = parseStrOrObject<T>(str)
    if (err) {
      console.error(err)
    } else if (parsed) {
      fromStorageValue = initialValue ? merge({}, initialValue, parsed) : parsed
    }
  }
  const validate = (value: S): S => {
    if (!value) {
      return value
    }
    if (!validator) {
      return value
    }
    const [v, errMsg] = validator(value)
    return {
      ...value,
      __validationMessage: errMsg,
      __validationPayload: v,
    } as any
  }
  const {
    update: _update,
    subscribe,
    set: _set,
  } = writable<S>(fromStorageValue ?? (initialValue as any))
  const _saveToStorageNow = (value: T) => {
    if (!storage || !_storage?.key) {
      return
    }
    if (_storage.redactKeys && !!value) {
      value = redactor(value, _storage.redactKeys)
    }
    storage?.setItem(_storage.key, value)
  }
  const saveToStorage =
    !!storage &&
    debounce(_saveToStorageNow, 2000, {
      leading: true,
      maxWait: 5000,
    })

  function didChange(existing) {
    const {
      __didChange: _,
      __validationMessage: _2,
      __validationPayload: _3,
      storeState,
      ...restNs
    } = existing
    const changed = !deepEqual(restNs, restoreValue)
    return changed
  }

  const update = (updater: Updater<S>, storeState?: StoreState) => {
    _update((s) => {
      let ns = storeState ? { ...updater(s), storeState } : updater(s)

      if (ns === s) {
        return ns
      }
      ns = validate(ns)

      // If there is no validator, we assume we dont care about changes.
      if (!storeState && validator) {
        if (restoreValue) {
          ns.__didChange = didChange(ns)
        }
      }

      if (saveToStorage) {
        storeState ? _saveToStorageNow(ns) : saveToStorage(ns)
      }
      return ns
    })
  }

  /** Like update, but also resets all store-state  */
  const restore = (state?: S) => {
    const s = update(
      () => {
        if (!state) {
          return merge({}, restoreValue)
        }
        const ns = merge({}, initialValue, state)
        restoreValue = ns
        return ns
      },
      {
        __didChange: false,
        __validationMessage: undefined,
        __validationPayload: undefined,
      }
    )
    return s
  }
  const set = (s: S) => {
    if (saveToStorage) {
      saveToStorage(s)
    }
    s = validate(s)
    s.__didChange = didChange(s)
    _set(merge({}, s))
  }
  const reset = () => {
    const s = validate({
      __didChange: false,
      __validationMessage: undefined,
      __validationPayload: undefined,
      ...(initialValue as any),
      // ...(initialValue as any),
    })

    _set(s)
    if (storage && _storage) {
      _saveToStorageNow(s)
    }
  }

  return {
    restore,
    reset,
    subscribe,
    update,
    set,
  }
}

export interface AppStorage<T extends {}> {
  getItem: (key: string) => string | null
  setItem: (key: string, value: T) => void
  key: string
}

export default createStore
