import api from "../../api";
import Vue from "vue";

const state = {
  all: {}
};

const getters = {
  device: state => deviceId => state.all[deviceId]
};

const actions = {
  async fetchDevice({ commit }, deviceId) {
    const device = await api.fetchDevice(deviceId);
    commit("setDevice", device);
  },

  async updateDevice({ commit, getters }, { deviceId, properties }) {
    const header = getters.device(deviceId);
    console.log("controllerName", header.controllerName);
    const device = await api.updateDevice(header, properties);
    commit("setDevice", device);
  },

  async updateDeviceProperty({ dispatch }, { deviceId, name, value }) {
    await dispatch("updateDevice", { deviceId, properties: { [name]: value } });
  }
};

const mutations = {
  setDevice(state, device) {
    Vue.set(state.all, device.identifier, device);
  }
};

export default {
  state,
  getters,
  actions,
  mutations
};
