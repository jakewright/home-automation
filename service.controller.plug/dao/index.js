const req = require("../../libraries/javascript/http");
const tpLinkClient = require("../api/tpLinkClient");
const Index = require("../domain/Plug");
const energyDecorator = require("../domain/energyDecorator");
const decorateDevice = require("../domain/decorateDevice");

const {
  store,
  updateDependencies
} = require("../../libraries/javascript/device");

const findById = identifier => {
  return store.findById(identifier);
};

const fetchAllState = async () => {
  (await req.get("service.registry.device/devices", {
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
  await updateDependencies(state, device.dependsOn);

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
  let device = new Index(header);

  if (header.attributes.energy) {
    decorateDevice(device, energyDecorator);
  }

  return device;
};

exports = module.exports = { findById, fetchAllState, watch, applyState };
