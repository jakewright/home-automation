class HueService {
  constructor(store, client) {
    this.store = store;
    this.client = client;
  }

  findById(identifier) {
    const light = this.store.findById(identifier);
    if (!light) throw new Error(`Light with id ${identifier} not found`);

    return light;
  }

  save(identifier, state) {
    const device = this.findById(identifier);

    return this.client
      .getLightById(device.hueId)
      .then(light => {
        applyStateToHuejay(state, light);
        return this.client.saveLight(light);
      })
      .then(light => {
        device.applyRemoteState(light);
        this.store.save(device);
      })
      .catch(err => {
        throw err;
      });
  }

  startPolling(interval = 5000) {
    this.pollingTimer = setInterval(() => {
      // Get all devices
      const devices = this.store.findAll()
      for (let id in devices) {
        // Fetch the remote state for this device
        this.client
          .getLightById(devices[id].hueId)
          .then(light => {
            // Apply the new state locally
            devices[id].applyRemoteState(light);
            this.store.save(devices[id]);
          })
          .catch(err => {
            console.error(
              `Failed to refresh state for device '${devices[id].identifier}':`,
              err
            );
          });
      }
    }, interval);
  }

  /**
   * Stop polling for state changes
   */
  stopPolling() {
    clearInterval(this.pollingTimer);
    this.pollingTimer = null;
  }
}

const applyStateToHuejay = (state, huejayLight) => {
  for (let property in state) {
    huejayLight[property] = state[property];
  }
};

exports = module.exports = HueService;
