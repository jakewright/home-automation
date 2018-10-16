const EventEmitter = require("events");

class DeviceStore extends EventEmitter {
  constructor(devices = []) {
    super();

    this.devices = {};
    this.cache = {};

    devices.forEach(this.add.bind(this));
    this.updateCache();
  }

  findById(identifier) {
    return this.devices[identifier];
  }

  findAll() {
    return Object.values(this.devices);
  }

  add(device) {
    if (device.identifier in this.devices)
      throw new Error(`Device ${device.identifier} already exists`);

    this.devices[device.identifier] = device;
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

exports = module.exports = DeviceStore;