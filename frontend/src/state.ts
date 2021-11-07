import createStore from './store'

export const state = createStore({
  initialValue: {
    tab: '',
    showDeleted: false,
    serverStats: false,
    codeLanguage: 'toml' as '' | 'toml' | 'json' | 'yaml',
    seenHints: {} as Record<string, [version: number, readAt: Date]>,
    collapse: {} as Record<string, boolean>,
  },
  storage: {
    key: 'state',
  },
})
