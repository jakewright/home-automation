const { Device } = require("../../libraries/javascript/device");

class HueLight extends Device {
  constructor(config) {
    super(config);
  }

  validate(state) {
    if ("brightness" in state) {
      if (state.brightness < 0 || state.brightness > 254) {
        return `Invalid brightness '${state.brightness}'`;
      }

      // Brightness cannot be set unless the light is on. In most cases, we
      // can just turn the light on as part of the same request, but this
      // would be weird if you're trying to set the brightness to zero so
      // disallow this edge case.
      if (state.brightness === 0 && !this.power) {
        return "Cannot set brightness to zero while light is off";
      }
    }
  }

  transform(state) {
    const t = {};

    if ("power" in state) {
      t.power = Boolean(state.power);
    }

    if ("brightness" in state) {
      t.brightness = state.brightness;
      t.power = state.brightness > 0;
    }

    return t;
  }

  getProperties() {
    return {
      power: { type: "bool" },
      brightness: {
        type: "int",
        min: 0,
        max: 254,
        interpolation: "continuous"
      }
    };
  }
}

exports = module.exports = HueLight;
