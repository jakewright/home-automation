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

    this.state = {};
  }

  getCommands() {
    return {};
  }

  getStateProviderUrls() {
    return this.stateProviders.map(
      ctrl => `${ctrl}/provide-state/${this.identifier}`
    );
  }

  applyState(state) {
    for (const property in this.state) {
      if (property in state) {
        this.state[property].value = state[property];
      }
    }
  }

  toJSON() {
    return {
      identifier: this.identifier,
      name: this.name,
      type: this.type,
      controllerName: this.controllerName,
      state: this.state,
      commands: this.getCommands()
    };
  }
}

exports = module.exports = Device;
