const redis = require("redis");
const config = require("../config");

// Redis needs separate clients for publishing and subscribing (I think)
// https://github.com/NodeRedis/node_redis#publish--subscribe
let pubClient;
let subClient;

// A map of pattern to array of handlers
let handlers = {};

const newClient = () => {
  if (!config.has("redis.host")) {
    throw new Error("No redis host defined in config");
  }

  const client = redis.createClient({
    host: config.get("redis.host"),
    port: config.get("redis.port")
  });

  client.on("ready", () => console.log("Redis connection is ready"));
  client.on("connect", () => console.log("Redis connected"));
  client.on("reconnecting", o =>
    console.log(`Reconnecting to Redis [attempt ${o.attempt}]...`)
  );
  client.on("error", err => console.error(`Redis error: ${err}`));
  client.on("end", () => console.log(`Redis connection closed`));
  client.on("warning", warn => console.error(`Redis warning: ${warn}`));

  return client;
};

const getPubClient = () => {
  if (pubClient === undefined) pubClient = newClient();
  return pubClient;
};

const getSubClient = () => {
  if (subClient === undefined) {
    subClient = newClient();

    // Every time we get a message on a channel that matches an active
    // pattern subscription, call all of the handlers for that pattern.
    subClient.on("pmessage", (pattern, channel, message) => {
      for (const handler of handlers[pattern]) {
        handler(channel, message);
      }
    });
  }

  return subClient;
};

const publish = (name, msg) => {
  getPubClient().publish(name, msg);
};

const subscribe = (pattern, handler) => {
  getSubClient().psubscribe(pattern);

  if (pattern in handlers) {
    handlers[pattern].push(handler);
  } else {
    handlers[pattern] = [handler];
  }
};

exports = module.exports = { publish, subscribe };
