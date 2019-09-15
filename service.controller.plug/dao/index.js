const req = require("../../libraries/javascript/http");
const tpLinkClient = require("../api/tpLinkClient");
const Hs100 = require("../domain/Hs100");
const Hs110 = require("../domain/Hs110");

const {
  store,
  updateDependencies
} = require("../../libraries/javascript/device");

const findById = identifier => {
  return store.findById(identifier);
};

const findAll = () => {
  return store.findAll();
};

const fetchAllState = async () => {
  (await req.get("service.device-registry/devices", {
    controllerName: "service.controller.plug"
  }))
    .map(instantiateDevice)
    .map(store.addDevice.bind(store));

  // Get all plug state and apply to local objects
  const promises = store.findAll().map(async device => {
    const state = await tpLinkClient.getStateByHost(device.attributes.host);
    device.applyState(state);
  });

  return Promise.all(promises);
};

const watch = interval => {
  console.log("Polling for state changes");

  setInterval(() => {
    fetchAllState().catch(err => {
      console.error("Failed to refresh state", err);
    });
  }, interval);
};

const applyState = async (device, state) => {
  // Update dependencies
  // await updateDependencies(state, device.dependsOn);

  // Update plug
  const success = await tpLinkClient.applyState(device.attributes.host, state);
  if (!success) throw new Error("Failed to apply state");

  // Apply new state to local device
  device.applyState(state);

  // Emit state change events
  store.flush();
};

const instantiateDevice = header => {
  header.controllerName = "service.controller.plug";

  if (header.attributes.energy) {
    return new Hs110(header);
  }

  return new Hs100(header);
};

exports = module.exports = {
  findById,
  findAll,
  fetchAllState,
  watch,
  applyState
};
