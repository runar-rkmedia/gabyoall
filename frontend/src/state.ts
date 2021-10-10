import createStore from 'store'

export const state = createStore({
  initialValue: {
    tab: '',
  },
  storage: {
    key: 'state',
  },
})
