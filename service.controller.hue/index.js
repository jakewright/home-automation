const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const firehose = require("../libraries/javascript/firehose");
const router = require("../libraries/javascript/router");
const { store } = require("../libraries/javascript/device");
const light = require("./light");
require("./handler/router");

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

    return light.fetchAllState();
  })
  .then(() => {
    router.listen();

    // Poll for state changes
    if (config.get("polling.enabled", false)) {
      light.watch(config.get("polling.interval", 30000));
    }
  })
  .catch(err => {
    console.error("Error initialising service", err);
  });
