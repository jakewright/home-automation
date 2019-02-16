const config = require("../../libraries/javascript/config");
const huejay = require("huejay");
const { fromDomain, toDomain } = require("./marshaling");

class HueClient {
  getClient() {
    if (this.client === undefined) this.connect();
    return this.client;
  }

  connect() {
    console.log("Connecting to Hue Bridge");
    this.client = new huejay.Client({
      host: config.get("hueBridge.host"),
      username: config.get("hueBridge.username"),
    });
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

const hueClient = new HueClient();
exports = module.exports = hueClient;
