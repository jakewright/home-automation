const EventEmitter = require("events");

class Store extends EventEmitter {
  constructor() {
    super();

    this.devices = {};
    this.cache = {};
  }

  addDevices(devices) {
    devices.forEach(this.addDevice.bind(this));
    this.updateCache();
  }

  addDevice(device) {
    if (device.identifier in this.devices)
      throw new Error(`Device ${device.identifier} already exists`);

    this.devices[device.identifier] = device;
  }

  findById(identifier) {
    return this.devices[identifier];
  }

  findAll() {
    return Object.values(this.devices);
  }

  flush() {
    for (let key in this.devices) {
      const oldState = this.cache[key];
      const newState = this.devices[key];

      if (oldState === undefined) {
        super.emit("device-added", key, oldState, newState);
      } else if (JSON.stringify(oldState) !== JSON.stringify(newState)) {
        super.emit("device-changed", key, oldState, newState);
      }
    }

    for (let key in this.cache) {
      const oldState = this.cache[key];
      const newState = this.devices[key];

      if (newState === undefined) {
        super.emit("device-removed", key, oldState, newState);
      }
    }

    this.updateCache();
  }

  updateCache() {
    for (let key in this.devices) {
      this.cache[key] = this.devices[key].toJSON();
    }
  }
}

const store = new Store();
exports = module.exports = store;
