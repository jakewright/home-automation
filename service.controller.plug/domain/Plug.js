const { Device } = require("../../libraries/javascript/device");

class Plug extends Device {
  transform(state) {
    const t = {};

    if ("power" in state) {
      t.power = Boolean(state.power);
    }

    return t;
  }

  getProperties() {
    return {
      power: { type: "bool" }
    };
  }
}

exports = module.exports = Plug;
