const req = require("../../libraries/javascript/request");
const HueLight = require("../domain/HueLight");
const colorDecorator = require("../domain/colorDecorator");
const colorTempDecorator = require("../domain/colorTempDecorator");
const rgbDecorator = require("../domain/rgbDecorator");
const decorateDevice = require("../domain/decorateDevice");
const {
  store,
  updateDependencies
} = require("../../libraries/javascript/device");
const hueClient = require("../api/hueClient");

const findById = identifier => {
  return store.findById(identifier);
};

const fetchAllState = async () => {
  // Get all devices from the registry and add them to the store
  (await req.get("service.registry.device/devices", {
    controllerName: "service.controller.hue"
  }))
    .map(instantiateDevice)
    .map(store.addDevice.bind(store));

  // Get all light state and apply to local objects
  const hueIdToState = await hueClient.getAllLights();
  for (const hueId in hueIdToState) {
    const device = findByHueId(hueId);
    if (!device) continue;
    device.applyState(hueIdToState[hueId]);
  }
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

const instantiateDevice = header => {
  let device = new HueLight({
    identifier: header.identifier,
    name: header.name,
    type: header.type,
    controllerName: "service.controller.hue",
    hueId: header.attributes.hueId
  });

  for (let feature of header.attributes.features) {
    switch (feature) {
      case "color":
        decorateDevice(device, colorDecorator);
        break;
      case "color-temp":
        decorateDevice(device, colorTempDecorator);
        break;
      case "rgb":
        decorateDevice(device, rgbDecorator);
        break;
      default:
        console.error(`Unknown light feature: ${feature}`);
    }
  }

  return device;
};

exports = module.exports = { findById, fetchAllState, watch, applyState };
