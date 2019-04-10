const decorateDevice = (device, decorator) => {
  if (typeof decorator.validate === "function") {
    const validate = device.validate;
    device.validate = state => {
      const err = validate.call(device, state);
      if (err !== undefined) return err;
      return decorator.validate.call(device, state);
    };
  }

  if (typeof decorator.transform === "function") {
    const transform = device.transform;
    device.transform = state => {
      let t = transform.call(device, state);
      return decorator.transform.call(device, state, t);
    };
  }

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
