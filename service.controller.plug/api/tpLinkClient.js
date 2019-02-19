const { Client } = require("tplink-smarthome-api");

class TpLinkClient {
  constructor() {
    this.client = new Client();
    this.plugs = {};
  }

  async getStateByHost(host) {
    if (!host) {
      throw new Error("Host is not set");
    }

    const info = await this.getPlug(host).getInfo();

    return {
      power: Boolean(info.sysInfo.relay_state),
      watts: info.emeter.realtime.power
    };
  }

  async applyState(host, state) {
    if (!("power" in state)) return;
    return this.getPlug(host).setPowerState(state.power);
  }

  // Private
  getPlug(host) {
    if (!(host in this.plugs)) {
      this.plugs[host] = this.client.getPlug({ host });
    }

    return this.plugs[host];
  }
}

const tpLinkClient = new TpLinkClient();
exports = module.exports = tpLinkClient;
