const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const hueClient = require("./api/hueClient");
const router = require("../libraries/javascript/router");
const dao = require("./dao");
require("./routes");

bootstrap("service.controller.hue")
  .then(() => {
    // Get Hue Bridge info
    if (!config.has("hueBridge.host") || !config.has("hueBridge.username")) {
      throw new Error("Host and username must be set in config");
    }
    hueClient.setHost(config.get("hueBridge.host"));
    hueClient.setUsername(config.get("hueBridge.username"));

    return dao.fetchAllState();
  })
  .then(() => {
    router.listen();

    // Poll for state changes
    if (config.get("polling.enabled", false)) {
      dao.watch(config.get("polling.interval", 30000));
    }
  })
  .catch(err => {
    console.error("Error initialising service", err);
  });
