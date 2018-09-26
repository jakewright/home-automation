const EventEmitter = require("events");

class Store extends EventEmitter {
  constructor(devices) {
    super();

    this.devices = devices;
    this.cache = {};
  }

  findAll() {
    return this.devices;
  }

  findById(identifier) {
    return this.devices[identifier];
  }

  save(device) {
    // Update the map to hold the new state
    this.devices[device.identifier] = device;

    // Get the old hash from the cache
    const old = this.cache[device.identifier];

    // Update the cache
    this.cache[device.identifier] = hash(device);

    // If the device didn't exist before, return early.
    if (!old) return;

    // If the state has changed, emit an event.
    if (old !== hash(device)) super.emit("device-state-changed", device);
  }
}

const hash = device => {
  if (!device) return "";
  return JSON.stringify(device.toJSON());
};

exports = module.exports = Store;
