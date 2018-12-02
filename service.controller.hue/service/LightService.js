const { updateDependencies } = require("../../libraries/javascript/device");

class LightService {
  constructor(store, apiClient, hueClient) {
    this.store = store;
    this.apiClient = apiClient;
    this.hueClient = hueClient;
  }

  findById(identifier) {
    return this.store.findById(identifier);
  }

  findByHueId(hueId) {
    return this.store.findAll().find(device => device.hueId == hueId);
  }

  fetchAllState() {
    return this.hueClient.getAllLights().then(hueIdToState => {
      for (const hueId in hueIdToState) {
        const device = this.findByHueId(hueId);
        if (!device) continue;
        device.applyState(hueIdToState[hueId]);
      }
    });
  }

  async applyState(device, state) {
    // Update dependencies
    await updateDependencies(state, device.dependsOn);

    // Update light
    const newState = await this.hueClient.applyState(device.hueId, state);

    // Apply new state to local device
    device.applyState(newState);

    // Emit state change events
    this.store.flush();
  }
}

exports = module.exports = LightService;
