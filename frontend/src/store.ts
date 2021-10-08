import { Updater, writable } from 'svelte/store'
import merge from 'lodash.merge'

function createStore<T extends {}>({
  storage: _storage,
  initialValue,
}: {
  storage?: AppStorage<T> | { key: string }
  initialValue?: T
} = {}) {
  const storage: AppStorage<T> | null = _storage?.key
    ? {
        key: _storage.key,
        getItem: (key) => localStorage.getItem(key),
        // TODO: throttle saving
        setItem: (k, v) => localStorage.setItem(k, JSON.stringify(v)),
        ..._storage,
      }
    : null

  if (storage) {
    const str = storage.getItem(storage.key)
    let parsed: T
    if (str) {
      if (typeof str === 'string') {
        try {
          const j = JSON.parse(str) as T
          parsed = j
        } catch (err) {
          console.error('failed to parse from localstorage key "store": ', err)
        }
      } else {
        parsed = str
      }
    }
    if (parsed) {
      initialValue = initialValue ? merge(initialValue, parsed) : parsed
    }
  }
  const { update: _update, subscribe } = writable<T>(initialValue)

  const update = (updater: Updater<T>) => {
    _update((s) => {
      const ns = updater(s)
      if (storage) {
        storage.setItem(storage.key, ns)
      }
      return ns
    })
  }

  return {
    subscribe,
    update,
  }
}

export interface AppStorage<T extends {}> {
  getItem: (key: string) => string
  setItem: (key: string, value: T) => void
  key: string
}

export default createStore
