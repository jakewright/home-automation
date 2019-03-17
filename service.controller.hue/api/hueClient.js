const axios = require("axios");
const { fromDomain, toDomain } = require("./marshaling");

class HueClient {
  setHost(host) {
    this.host = host;
  }

  setUsername(username) {
    this.username = username;
  }

  getClient() {
    if (this.client === undefined) {
      this.client = axios.create({
        baseURL: `${this.host}/api/${this.username}`
      });
    }

    return this.client;
  }

  async request(method, url, data) {
    const rsp = await this.getClient().request({ method, url, data });
    if (JSON.stringify(rsp.data).includes("error")) {
      throw new Error(`Hue response included errors:
      ${method} ${url}
      ${JSON.stringify(rsp.data)}`);
    }
    return rsp;
  }

  async fetchAllState() {
    const rsp = await this.request("get", "/lights");
    const lights = {};

    for (const hueId in rsp.data) {
      lights[hueId] = toDomain(rsp.data[hueId]);
    }

    return lights;
  }

  async fetchState(hueId) {
    const rsp = await this.request("get", `/lights/${hueId}`);
    return toDomain(rsp.data);
  }

  async applyState(hueId, state) {
    state = fromDomain(state);
    await this.request("put", `/lights/${hueId}/state`, state);
    return this.fetchState(hueId);
  }
}

const hueClient = new HueClient();
exports = module.exports = hueClient;
