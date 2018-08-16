import huejay from "huejay";

export default class HueBridgeClient {
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

  discover() {
    return this.getClient().discover();
  }

  connect() {
    console.log("Connecting to Hue Bridge");
    this.client = new huejay.Client(this.config);
  }

  createUser() {
    if (this.config.username) throw new Error("User is already set");
    let user = new this.client.users.User();
    return this.client.users.create(user);
  }

  getAllUsers() {
    return this.client.users.getAll();
  }

  getLightById(id) {
    return this.client.lights.getById(id);
  }

  getAllLights() {
    return this.client.lights.getAll();
  }

  saveLight(light) {
    return this.client.lights.save(light);
  }
}
