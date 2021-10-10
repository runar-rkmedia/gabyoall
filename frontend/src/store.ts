import { Updater, writable } from 'svelte/store'
import merge from 'lodash.merge'
import debounce from 'lodash.debounce'

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

function createStore<T extends {}>({
  storage: _storage,
  initialValue,
}: {
  storage?: AppStorage<T> | { key: string }
  initialValue?: T
} = {}) {
  const storage: AppStorage<T> | null = _storage?.key
    ? {
        getItem: (key) => localStorage.getItem(key),
        // TODO: throttle saving
        setItem: (k, v) => localStorage.setItem(k, JSON.stringify(v)),
        ..._storage,
      }
    : null

  if (storage) {
    const str = storage.getItem(storage.key)
    const [parsed, err] = parseStrOrObject<T>(str)
    if (err) {
      console.error(err)
    } else if (parsed) {
      initialValue = initialValue ? merge(initialValue, parsed) : parsed
    }
  }
  const { update: _update, subscribe, set: _set } = writable<T>(initialValue)
  const saveToStorage = debounce(
    (value: T) => {
      if (storage && _storage?.key) {
        storage?.setItem(_storage.key, value)
      }
    },
    2000,
    {
      leading: true,
      maxWait: 5000,
    }
  )

  const update = (updater: Updater<T>) => {
    _update((s) => {
      const ns = updater(s)
      if (storage) {
        saveToStorage(ns)
      }
      return ns
    })
  }
  const set = (value: T) => {
    if (storage && _storage?.key) {
      saveToStorage(value)
    }
    _set(value)
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
