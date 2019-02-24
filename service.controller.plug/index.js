const bootstrap = require("../libraries/javascript/bootstrap");
const firehose = require("../libraries/javascript/firehose");
const { store } = require("../libraries/javascript/device");
const dao = require("./dao");
const router = require("../libraries/javascript/router");
const config = require("../libraries/javascript/config");
require("./routes");

// Create and initialise a Service object
const serviceName = "service.controller.plug";
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
