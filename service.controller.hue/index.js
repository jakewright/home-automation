const bootstrap = require("../libraries/javascript/bootstrap");
const config = require("../libraries/javascript/config");
const firehose = require("../libraries/javascript/firehose");
const req = require("../libraries/javascript/request");
const router = require("../libraries/javascript/router");
const { store } = require("../libraries/javascript/device");

const HueLight = require("./domain/HueLight");
const colorDecorator = require("./domain/colorDecorator");
const colorTempDecorator = require("./domain/colorTempDecorator");
const rgbDecorator = require("./domain/rgbDecorator");
const decorateDevice = require("./domain/decorateDevice");
const { fetchAllState } = require("./light");
require("./handler/router");

const serviceName = "service.controller.hue";
bootstrap(serviceName)
  .then(() => {
    // Get device headers
    return req.get("service.registry.device/devices", {
      controllerName: serviceName
    });
  })
  .then(deviceHeaders => {
    // Instantiate devices and add to store
    const devices = deviceHeaders.map(instantiateDevice);
    store.addDevices(devices);

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

const instantiateDevice = header => {
  let device = new HueLight({
    identifier: header.identifier,
    name: header.name,
    type: header.type,
    controllerName: serviceName,
    hueId: header.attributes.hueId
  });

  for (let feature of header.attributes.features) {
    switch (feature) {
      case "color":
        decorateDevice(device, colorDecorator);
        break;
      case "color-temp":
        decorateDevice(device, colorTempDecorator);
        break;
      case "rgb":
        decorateDevice(device, rgbDecorator);
        break;
      default:
        console.error(`Unknown light feature: ${feature}`);
    }
  }

  return device;
};
