class Device {
  /**
   * @param {Object} config Configuration for the device
   * @param {string} config.identifier The unique identifier for the device
   * @param {string} config.name The name of the device
   * @param {string} config.type The type of the device
   * @param {string} config.controllerName The name of this controller
   * @param {Object} config.attributes Extra information about the device
   * @param {array} config.dependsOn Array of objects describing dependencies
   * @param {array} config.stateProviders Array of state provider service names
   */
  constructor(config) {
    this.identifier = config.identifier;
    this.name = config.name;
    this.type = config.type;
    this.controllerName = config.controllerName;
    this.attributes = config.attributes || {};
    this.dependsOn = config.dependsOn || [];
    this.stateProviders = config.stateProviders || [];
  }

  getCommands() {
    return {};
  }

  getProperties() {
    return {};
  }

  getStateProviderUrls() {
    return this.stateProviders.map(ctrl => `${ctrl}/device/${this.identifier}`);
  }

  applyState(state) {
    for (const property in this.getProperties()) {
      if (property in state) {
        this[property] = state[property];
      }
    }
  }

  toJSON() {
    let json = {
      identifier: this.identifier,
      name: this.name,
      type: this.type,
      controllerName: this.controllerName,
      availableProperties: this.getProperties(),
      commands: this.getCommands()
    };

    for (let property in this.getProperties()) {
      json[property] = this[property];
    }

    return json;
  }
}

exports = module.exports = Device;
