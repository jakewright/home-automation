const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const firehose = require("../libraries/javascript/firehose");
const req = require("../libraries/javascript/request");
const router = require("../libraries/javascript/router");
const { store } = require("../libraries/javascript/device");
const { fetchAllState } = require("./light");
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

    fetchAllState().catch(err => {
      console.error("Failed to fetch state", err);
    });

    // Start the server
    router.listen();

    // Poll for state changes
    if (config.get("polling.enabled", false)) {
      console.log("Polling for state changes");

      let pollingTimer = setInterval(() => {
        fetchAllState().catch(err => {
          console.error("Failed to refresh state", err);
        });
      }, config.get("polling.interval", 30000));
    }
  })
  .catch(err => {
    console.error("Error initialising service", err);
  });


