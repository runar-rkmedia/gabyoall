import createStore from 'store'

export const state = createStore({
  initialValue: {
    tab: '',
    codeLanguage: 'toml' as '' | 'toml' | 'json' | 'yaml',
    seenHints: {} as Record<string, [version: number, readAt: Date]>,
  },
  storage: {
    key: 'state',
  },
})
