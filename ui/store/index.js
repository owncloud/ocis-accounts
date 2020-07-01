import {
  // eslint-disable-next-line camelcase
  AccountsService_ListAccounts
} from '../client/accounts'
import axios from 'axios'

const state = {
  config: null,
  initialized: false,
  accounts: {}
}

const getters = {
  config: state => state.config,
  isInitialized: state => state.initialized,
  getAccountsSorted: state => {
    // FIXME: look at data fields for available sorting options
    return Object.values(state.accounts)
  }
}

const mutations = {
  LOAD_CONFIG (state, config) {
    state.config = config
  },
  SET_INITIALIZED (state, value) {
    state.initialized = value
  },
  SET_ACCOUNTS (state, accounts) {
    state.accounts = accounts
  }
}

const actions = {
  loadConfig ({ commit }, config) {
    commit('LOAD_CONFIG', config)
  },

  async initialize ({ commit, dispatch }) {
    await dispatch('fetchAccounts')
    commit('SET_INITIALIZED', true)
  },

  async fetchAccounts ({ commit, dispatch, getters, rootGetters }) {
    injectAuthToken(rootGetters)
    const response = await AccountsService_ListAccounts({
      $domain: rootGetters.configuration.server,
      body: {}
    })
    if (response.status === 201) {
      const accounts = response.data.accounts
      commit('SET_ACCOUNTS', accounts || [])
    } else {
      dispatch('showMessage', {
        title: 'Failed to fetch accounts.',
        desc: response.statusText,
        status: 'danger'
      }, { root: true })
    }
  }
}

export default {
  namespaced: true,
  state,
  getters,
  actions,
  mutations
}

function injectAuthToken (rootGetters) {
  axios.interceptors.request.use(config => {
    if (typeof config.headers.Authorization === 'undefined') {
      const token = rootGetters.user.token
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }
    }
    return config
  })
}