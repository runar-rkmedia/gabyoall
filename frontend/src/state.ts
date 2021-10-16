import createStore from 'store'

export const state = createStore({
  initialValue: {
    tab: '',
    codeLanguage: 'toml' as '' | 'toml' | 'json' | 'yaml',
  },
  storage: {
    key: 'state',
  },
})
