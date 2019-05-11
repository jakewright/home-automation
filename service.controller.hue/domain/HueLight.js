const { Device } = require("../../libraries/javascript/device");

class HueLight extends Device {
  constructor(config) {
    super(config);

    const hueId = this.attributes.hueId;
    if (typeof hueId !== "string") {
      throw new Error(`Hue ID '${hueId}' is not a string`);
    }

    this.state = {
      power: { type: "bool" },
      brightness: {
        type: "int",
        min: 0,
        max: 254,
        interpolation: "continuous"
      }
    };
  }

  validate(state) {
    if ("brightness" in state) {
      const brightness = parseInt(state.brightness, 10);

      if (isNaN(brightness)) {
        return `Brightness is not a valid number '${state.brightness}'`;
      }

      if (brightness < 0 || brightness > 254) {
        return `Invalid brightness '${state.brightness}'`;
      }

      // Brightness cannot be set unless the light is on. In most cases, we
      // can just turn the light on as part of the same request, but this
      // would be weird if you're trying to set the brightness to zero so
      // disallow this edge case.
      if (state.brightness === 0 && !this.state.power.value) {
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
      t.brightness = parseInt(state.brightness, 10);
      t.power = state.brightness > 0;
    }

    return t;
  }
}

exports = module.exports = HueLight;
