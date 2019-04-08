import api from "../../api";
import _ from "lodash";
import Vue from "vue";

const state = {
  all: {}
};

const getters = {
  // Convert the map of rooms into an array
  allRooms: state => Object.values(state.all),

  room: state => roomId => state.all[roomId]
};

const actions = {
  async fetchRooms({ commit }) {
    const rooms = await api.fetchRooms();
    commit("setRooms", rooms);
  },

  async fetchRoom({ commit }, roomId) {
    const room = await api.fetchRoom(roomId);
    commit("setRoom", room);
  }
};

const mutations = {
  setRooms(state, rooms) {
    state.all = _.keyBy(rooms, room => room.identifier);
  },
  setRoom(state, room) {
    Vue.set(state.all, room.identifier, room);
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
