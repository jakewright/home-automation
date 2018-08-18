const EventEmitter = require("events");
const HueBridgeClient = require("../api/HueBridgeClient");

class HueLight extends EventEmitter {
  /**
   * @param {Object} config Configuration for the device
   * @param {string} config.identifier The unique identifier for the device
   * @param {string} config.name The name of the device
   * @param {string} config.type The type of the device
   * @param {string} config.controllerName The name of this controller
   * @param {number} config.hueId The ID of the lamp on the Hue bridge
   * @param {HueBridgeClient} config.client An instance of the HueBridgeClient
   */
  constructor(config) {
    super();

    this.identifier = config.identifier;
    this.name = config.name;
    this.type = config.type;
    this.controllerName = config.controllerName;
    this.hueId = config.hueId;
    this.client = config.client;

    this.power = false;
    this.brightness = 0;
  }

  /**
   * Poll the device every `interval` milliseconds
   * @param {number} interval ms
   */
  startPolling(interval = 5000) {
    this.pollingTimer = setInterval(() => {
      this.fetchRemoteState()
        .then(this.applyRemoteState)
        .catch(err => {
          console.error(
            `An error occurred while refreshing state for device ${
              this.identifier
            }: ${err.message}`
          );
        });
    }, interval);
  }

  /**
   * Stop polling for state changes
   */
  stopPolling() {
    clearInterval(this.pollingTimer);
    this.pollingTimer = null;
  }

  /**
   * Set the state of the device from the given object. An error will be thrown if validation fails.
   * @param {state} Object
   */
  setState(state) {
    if ("power" in state) this.setPower(state.power);
    if ("brightness" in state) this.setBrightness(state.brightness);
  }

  /**
   * Apply the local state to the remote bulb via the Hue Bridge
   */
  save() {
    this.prepareLight();
    return this.client.saveLight(this.light).then(this.applyRemoteState);
  }

  /**
   * Set all of the properties on this.light before sending to the hue bridge
   */
  prepareLight() {
    // The light object is returned from the Huejay library and is needed to save changes, so
    // fetchRemoteState() must be called before save() (and thus prepareLight()).
    if (this.light === null) {
      throw new Error("State must be fetched before saving is allowed");
    }

    if (this.light.on !== this.power) {
      this.light.on = this.power;
    }

    if (!this.power) return;

    const brightness = Math.floor(this.brightness * 2.54);
    if (this.light.brightness !== brightness) {
      this.light.brightness = brightness;
    }
  }

  /**
   * Get the up-to-date state of the light from the Hue bridge
   */
  fetchRemoteState() {
    return Promise.resolve(this.client.getLightById(this.hueId));
  }

  /**
   * Set the local properties based on the response from the Hue bridge
   * @param light
   */
  applyRemoteState(light) {
    // Save this for later when we want to update the light
    this.light = light;

    this.setPower(light.on);

    // Don't use the setter for brightness because it will override
    // the power state.
    this.brightness = Math.ceil(light.brightness / 2.54);

    const oldCache = this.cache;
    this.cache = this.createCache();

    // Return early if this is the first invocation of this function
    if (oldCache === undefined) return;

    // If the state has changed, emit an event
    if (JSON.stringify(oldCache) !== JSON.stringify(this.cache)) {
      super.emit("state-change", this.toJSON());
    }
  }

  /**
   * Turn the light on or off
   * @param {boolean} state
   */
  setPower(state) {
    this.power = Boolean(state);
  }

  /**
   * Set the brightness of the light
   * @param {number} value Brightness value between 0-100
   */
  setBrightness(value) {
    if (value < 0 || value > 100) {
      throw new Error(`Invalid brightness '${value}'`);
    }

    this.brightness = value;
    this.power = value > 0;
  }

  /**
   * Return an object representing this device that can be marshalled into a JSON response.
   * @return {Object}
   */
  toJSON() {
    return {
      identifier: this.identifier,
      name: this.name,
      type: this.type,
      controllerName: this.controllerName,
      availableProperties: this.getProperties(),
      power: this.power,
      brightness: this.brightness
    };
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
        max: 100,
        interpolation: "continuous"
      }
    };
  }

  /**
   * Return an object representing the current state of the device.
   */
  createCache() {
    let cache = {};
    for (const property in this.getProperties()) {
      cache[property] = this[property];
    }

    return cache;
  }
}

exports = module.exports = HueLight;
