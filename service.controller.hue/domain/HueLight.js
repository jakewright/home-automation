class HueLight {
  /**
   * @param {Object} config Configuration for the device
   * @param {string} config.identifier The unique identifier for the device
   * @param {string} config.name The name of the device
   * @param {string} config.type The type of the device
   * @param {string} config.controllerName The name of this controller
   * @param {number} config.hueId The ID of the lamp on the Hue bridge
   */
  constructor(config) {
    this.identifier = config.identifier;
    this.name = config.name;
    this.type = config.type;
    this.controllerName = config.controllerName;
    this.hueId = config.hueId;
  }

  /**
   * Transform validates and manipulates the state into a form that is
   * ready to apply to the light.
   */
  transform(state) {
    const t = {};

    if ("power" in state) {
      t.on = Boolean(state.power);
    }

    if ("brightness" in state) {
      if (state.brightness < 0 || state.brightness > 254) {
        throw new Error(`Invalid brightness '${state.brightness}'`);
      }

      // Brightness cannot be set unless the light is on. In most cases, we
      // can just turn the light on as part of the same request, but this
      // would be weird if you're trying to set the brightness to zero so
      // disallow this edge case.
      if (state.brightness == 0 && !this.power) {
        throw new Error(`Cannot set brightness to zero while light is off`);
      }

      t.brightness = state.brightness;
      t.on = t.brightness > 0;
    }

    return t;
  }

  /**
   * Set the state of the device from the given object. An error will be thrown if validation fails.
   * @param {state} Object
   */
  applyRemoteState(state) {
    this.power = state.on;
    this.brightness = state.brightness;
  }

  /**
   * Return an object representing this device that can be marshalled into a JSON response.
   * @return {Object}
   */
  toJSON() {
    let json = {
      identifier: this.identifier,
      name: this.name,
      type: this.type,
      controllerName: this.controllerName,
      availableProperties: this.getProperties()
    };

    for (let property in this.getProperties()) {
      json[property] = this[property];
    }

    return json;
  }

  /**
   * Return an object where the keys are the setable properties on this device and the values
   * represent the values that the properties can take.
   */
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
