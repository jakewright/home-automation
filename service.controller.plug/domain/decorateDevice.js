const decorateDevice = (device, decorator) => {
  if (typeof decorator.getCommands === "function") {
    const getCommands = device.getCommands;
    device.getCommands = () => {
      let commands = getCommands.call(device);
      return decorator.getCommands.call(device, commands);
    };
  }

  if (typeof decorator.getProperties === "function") {
    const getProperties = device.getProperties;
    device.getProperties = () => {
      let properties = getProperties.call(device);
      return decorator.getProperties.call(device, properties);
    };
  }
};

exports = module.exports = decorateDevice;
