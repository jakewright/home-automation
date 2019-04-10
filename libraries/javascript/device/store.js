const EventEmitter = require("events");
const firehose = require("../firehose");

class Store extends EventEmitter {
  constructor() {
    super();

    this.devices = {};
    this.cache = {};
  }

  addDevice(device) {
    if (device.identifier in this.devices) return;

    this.devices[device.identifier] = device;
    this.updateCache();
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
      } else if (oldState !== JSON.stringify(newState)) {
        console.log("Device state changed");
        super.emit("device-changed", key, newState);
      }
    }

    for (let key in this.cache) {
      const newState = this.devices[key];

      if (newState === undefined) {
        super.emit("device-removed", key, newState);
      }
    }

    this.updateCache();
  }

  updateCache() {
    for (let key in this.devices) {
      // The object needs to be stringified otherwise you get
      // shallow copy where the state mutates with the real device.
      this.cache[key] = JSON.stringify(this.devices[key].toJSON());
    }
  }
}

const store = new Store();

// Emit firehose events when device state is changed
store.on("device-changed", (identifier, newState) => {
  console.log(`State changed for device ${identifier}`);
  firehose.publish(
    `device-state-changed.${identifier}`,
    JSON.stringify(newState)
  );
});

exports = module.exports = store;
