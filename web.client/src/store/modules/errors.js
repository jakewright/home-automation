import Vue from "vue";

let id = 0;

const state = {
  all: {}
};

const getters = {
  // Convert the map of errors into an array
  allErrors: state => Object.values(state.all)
};

const actions = {
  enqueueError({ commit }, err) {
    err.id = id++;
    commit("setError", err);
  },

  removeError({ commit }, id) {
    return commit("removeError", id);
  }
};

const mutations = {
  setError(state, err) {
    Vue.set(state.all, err.id, err);
  },

  removeError(state, id) {
    Vue.delete(state.all, id);
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
