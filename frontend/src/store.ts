import { Updater, writable } from 'svelte/store'
import merge from 'lodash.merge'
import debounce from 'lodash.debounce'
import cloneDeep from 'lodash.clonedeep'
import type { Paths } from 'types'

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
        console.log('redactor replaced', key, obj[key], replacement)
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

export type Store<T extends {}, V = null, VK extends string = string> = T & {
  __didChange?: boolean
  __validationMessage?: Partial<Record<VK, string>>
  __validationPayload?: V
}

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
      initialValue = initialValue ? merge(initialValue, parsed) : parsed
    }
  }
  const { update: _update, subscribe, set: _set } = writable<S>(initialValue)
  const saveToStorage =
    !!storage &&
    debounce(
      (value: T) => {
        if (!storage || !_storage?.key) {
          return
        }
        if (_storage.redactKeys) {
          value = redactor(value, _storage.redactKeys)
        }
        storage?.setItem(_storage.key, value)
      },
      2000,
      {
        leading: true,
        maxWait: 5000,
      }
    )

  const validate =
    !!validator &&
    debounce(
      (value: S) => {
        const [v, errMsg] = validator(value)
        if (errMsg !== value.__validationMessage) {
          update(
            ({ __validationMessage, __validationPayload, ...f }) =>
              ({
                ...f,
                __validationPayload: v,
                __validationMessage: errMsg,
              } as any)
          )
        }
      },
      150,
      { leading: true, maxWait: 2000 }
    )

  const update = (updater: Updater<S>) => {
    _update((s) => {
      const ns = updater(s)
      if (ns === s) {
        return ns
      }
      validate && validate(s)
      ;(ns as any).__didChange = true
      if (saveToStorage) {
        saveToStorage(ns)
      }
      return ns
    })
  }
  const set = (s: S) => {
    if (saveToStorage) {
      saveToStorage(s)
    }
    validate && validate(s)
    s.__didChange = true
    _set(s)
  }

  return {
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
