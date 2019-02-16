const {
  store,
  updateDependencies
} = require("../../libraries/javascript/device");
const hueClient = require("../api/hueClient");

const findById = identifier => {
  return store.findById(identifier);
};

const fetchAllState = () => {
  return hueClient.getAllLights().then(hueIdToState => {
    for (const hueId in hueIdToState) {
      const device = findByHueId(hueId);
      if (!device) continue;
      device.applyState(hueIdToState[hueId]);
    }
  });
};

const applyState = async (device, state) => {
  // Update dependencies
  await updateDependencies(state, device.dependsOn);

  // Update light
  const newState = await hueClient.applyState(device.hueId, state);

  // Apply new state to local device
  device.applyState(newState);

  // Emit state change events
  store.flush();
};

const findByHueId = hueId => {
  return store.findAll().find(device => device.hueId === hueId);
};

exports = module.exports = { findById, fetchAllState, applyState };
