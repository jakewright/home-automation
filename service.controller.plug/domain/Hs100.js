const { Device } = require("../../libraries/javascript/device");

class Hs100 extends Device {
  constructor(config) {
    super(config);

    this.state = {
      power: { type: "bool" }
    };
  }

  transform(state) {
    const t = {};

    if ("power" in state) {
      t.power = Boolean(state.power);
    }

    return t;
  }
}

exports = module.exports = Hs100;
