const req = require("../../libraries/javascript/http");
const { store } = require ("../../libraries/javascript/device");
const OnkyoHTR380 = require("../domain/OnkyoHTR380");

const findById = identifier => {
  return store.findById(identifier);
};

const fetchAllState = async () => {
  // Get all devices from the registry and add them to the store
  (await req.get("service.device-registry/devices", {
    controllerName: "service.controller.infrared"
  }))
    .map(instantiateDevice)
    .map(store.addDevice.bind(store));

  // This usually emits state change events but none of the
  // devices have any state stored locally in this service yet.
  // This is where, in the future, we should fetch state from
  // state providers.
  store.flush();
};

const instantiateDevice = header => {
  switch (header.type) {
    case "infrared.htr380":
      return new OnkyoHTR380(header);
  }

  throw new Error(`Unknown device type: ${header.type}`)
};

exports = module.exports = { findById, fetchAllState };