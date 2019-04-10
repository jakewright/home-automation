const redis = require("redis");
const config = require("../config");

let client;

const getClient = () => {
  if (client === undefined) {
    if (!config.has("redis.host")) {
      throw new Error("No redis host defined in config");
    }

    client = redis.createClient({
      host: config.get("redis.host"),
      port: config.get("redis.port")
    });
    client.on("error", err => {
      console.error(`Redis error: ${err}`);
    });
  }

  return client;
};

const publish = (name, msg) => {
  getClient().publish(name, msg);
};

exports = module.exports = { publish };
