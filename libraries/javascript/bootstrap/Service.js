const bodyParser = require("body-parser");
const config = require("config");
const ApiClient = require("./ApiClient");
const express = require("express");
const redis = require("redis");

class Service {
  constructor() {
    this.config = config;
    // Connect to Redis
    if (config.has("redis.host")) {
      this.redisClient = redis.createClient({
        host: config.get("redis.host"),
        port: config.get("redis.port")
      });
      this.redisClient.on("error", err => {
        console.error(`Redis error: ${err}`);
      });
    }

    this.controllerName = config.get("controllerName");
    this.apiClient = new ApiClient(config.get("apiGateway"));

    this.app = express();

    // JSON body parser
    this.app.use(bodyParser.json());

    // Request logger
    this.app.use((req, res, next) => {
      console.log(
        `${req.method} ${req.originalUrl} ${JSON.stringify(req.body)}`
      );
      next();
    });
  }

  listen() {
    this.app.listen(config.get("port"), () => {
      console.log(`Service running on port ${config.get("port")}`);
    });
  }
}

exports = module.exports = Service;
