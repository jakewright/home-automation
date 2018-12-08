const bodyParser = require("body-parser");
const ApiClient = require("../ApiClient");
const express = require("express");
const redis = require("redis");
const Config = require("./Config");

class Service {
  constructor(controllerName) {
    this.controllerName = controllerName;
  }

  async init() {
    /* Initialise API ApiClient */
    const apiGateway = process.env.API_GATEWAY;
    if (!apiGateway) throw new Error("API_GATEWAY env var not set");
    this.apiClient = new ApiClient(apiGateway);

    /* Load config */
    let config = await this.apiClient.get(
      `service.config/read/${this.controllerName}`
    );
    this.config = new Config(config);

    /* Connect to Redis */
    if (this.config.has("redis.host")) {
      this.redisClient = redis.createClient({
        host: this.config.get("redis.host"),
        port: this.config.get("redis.port")
      });
      this.redisClient.on("error", err => {
        console.error(`Redis error: ${err}`);
      });
    }

    /* Create Express app */
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
    // Add an error handler that returns valid JSON
    this.app.use(function(err, req, res, next) {
      console.error(err.stack);
      res.status(500);
      res.json({ message: err.message });
    });

    const port = this.config.get("port", 80);
    this.app.listen(port, () => {
      console.log(`Service running on port ${port}`);
    });
  }
}

exports = module.exports = Service;
