const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const firehose = require("../libraries/javascript/firehose");
const router = require("../libraries/javascript/router");
const { store } = require("../libraries/javascript/device");
const dao = require("./dao");
require("./routes");

const serviceName = "service.controller.hue";
bootstrap(serviceName)
  .then(() => {
    // Subscribe to state changes from the store
    store.on("device-changed", (identifier, oldState, newState) => {
      console.log(`State changed for device ${identifier}`);
      firehose.publish(
        `device-state-changed.${identifier}`,
        JSON.stringify({ oldState, newState })
      );
    });

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
