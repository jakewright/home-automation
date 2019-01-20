const huejay = require("huejay");
const { fromDomain, toDomain } = require("./marshaling");

class HueClient {
  /**
   * @param {Object} config
   * @param {string} config.host Host of the Hue Bridge
   * @param {string} config.username Optional Hue Bridge username
   */
  constructor(config) {
    this.config = config;
    this.client = null;
  }

  getClient() {
    if (this.client === null) this.connect();
    return this.client;
  }

  connect() {
    console.log("Connecting to Hue Bridge");
    this.client = new huejay.Client(this.config);
  }

  discover() {
    return huejay.discover();
  }

  createUser() {
    if (this.config.username) throw new Error("User is already set");
    let user = new this.getClient().users.User();
    return this.getClient().users.create(user);
  }

  getAllUsers() {
    return this.getClient().users.getAll();
  }

  async getAllLights() {
    const lights = await this.getClient().lights.getAll();

    // Convert to a map where the keys are the Hue IDs
    return lights.reduce((map, light) => {
      map[light.id] = toDomain(light);
      return map;
    }, {});
  }

  async applyState(hueId, state) {
    state = fromDomain(state);
    let light = await this.getClient().lights.getById(hueId);

    // Apply the state to the Huejay light object
    for (let property in state) {
      light[property] = state[property];
    }

    light = await this.getClient().lights.save(light);
    return toDomain(light);
  }
}

exports = module.exports = HueClient;
