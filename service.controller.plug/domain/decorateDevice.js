const decorateDevice = (device, decorator) => {
  if (typeof decorator.getCommands === "function") {
    const getCommands = device.getCommands;
    device.getCommands = () => {
      let commands = getCommands.call(device);
      return decorator.getCommands.call(device, commands);
    };
  }

  if (typeof decorator.state === "object") {
    Object.assign(device.state, decorator.state);
  }
};

exports = module.exports = decorateDevice;
