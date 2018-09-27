const Store = require("../../libraries/javascript/store");
const hueClient = require("../api/HueClient");

const state = {};

const getters = {
  device: ({ state }, deviceId) => state[deviceId],
  devices: ({ state }) => state
};

const mutations = {
  setDevice({ state }, device) {
    state[device.identifier] = device;
  }
};

const actions = {
  async fetchDevice({ commit }, device) {
    const light = await hueClient.getLightById(device.hueId);
    device.applyRemoteState(light);
    commit("setDevice", device);
  },

  async saveDevice({ get, commit }, { device, state }) {
    let light = await hueClient.getLightById(device.hueId);
    applyStateToHuejay(state, light);
    light = await hueClient.saveLight(light);
    device.applyRemoteState(light);
    commit("setDevice", device);
  }
};

const store = new Store({
  state,
  getters,
  mutations,
  actions
});

const applyStateToHuejay = (state, huejayLight) => {
  for (let property in state) {
    huejayLight[property] = state[property];
  }
};

exports = module.exports = store;
