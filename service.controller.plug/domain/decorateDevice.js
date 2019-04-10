const decorateDevice = (device, decorator) => {
  if (typeof decorator.getCommands === "function") {
    const getCommands = device.getCommands;
    device.getCommands = () => {
      let commands = getCommands.call(device);
      return decorator.getCommands.call(device, commands);
    };
  }

  if (typeof decorator.state === "object") {
    // Deep clone the extra state and merge it with the device's existing state
    Object.assign(device.state, JSON.parse(JSON.stringify(decorator.state)));
  }
};

exports = module.exports = decorateDevice;
