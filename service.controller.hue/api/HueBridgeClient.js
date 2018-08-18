const huejay = require("huejay");

class HueBridgeClient {
  /**
   * @param {Object} config
   * @param {string} config.host Host of the Hue Bridge
   * @param {string} config.username Optional: Hue Bridge username
   */
  constructor(config) {
    this.client = null;
    this.config = config;
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
    let user = new this.client.users.User();
    return this.getClient().users.create(user);
  }

  getAllUsers() {
    return this.getClient().users.getAll();
  }

  getLightById(id) {
    return this.getClient().lights.getById(id);
  }

  getAllLights() {
    return this.getClient().lights.getAll();
  }

  saveLight(light) {
    return this.getClient().lights.save(light);
  }
}

exports = module.exports = HueBridgeClient;
