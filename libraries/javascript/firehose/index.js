const redis = require("redis");
const config = require("../config");

class Firehose {
  getClient() {
    if (this.client === undefined) {
      if (!config.has("redis.host")) {
        throw new Error("No redis host defined in config");
      }

      this.client = redis.createClient({
        host: config.get("redis.host"),
        port: config.get("redis.port")
      });
      this.client.on("error", err => {
        console.error(`Redis error: ${err}`);
      });
    }

    return this.client;
  }

  publish(name, msg) {
    this.getClient().publish(name, msg);
  }
}

const firehose = new Firehose();
exports = module.exports = firehose;
