const Service = require("../libraries/javascript/bootstrap");
const HueLight = require("./domain/HueLight");
const decorate = require("./domain/decorate");
const colorDecorator = require("./domain/colorDecorator");
const colorTempDecorator = require("./domain/colorTempDecorator");
const rgbDecorator = require("./domain/rgbDecorator");
const hueClient = require("./api/hueClient");
const DeviceController = require("./controllers/DeviceController");
const HueDiscoveryController = require("./controllers/HueBridgeController");
const HueBridgeController = require("./controllers/HueBridgeController");
const store = require("./store/devices");

/* Use the bootstrap library to create a Service object */
const service = new Service("service.controller.hue");
service.init()
  .then(() => {
    hueClient.setConfig({
      host: service.config.get("hueBridge.host"),
      username: service.config.get("hueBridge.username"),
    });

    /* Subscribe to state changes from the store */
    store.on("key-changed", key => {
      const device = store.get("device", key);
      console.log(`State changed for device ${device.identifier}`);
      service.redisClient.publish(
        `device-state-changed.${device.identifier}`,
        JSON.stringify(device)
      );
    });

    /* Initialise the devices */
    service.apiClient
      .get("service.registry.device/devices", { controllerName: service.controllerName })
      .then(deviceHeaders => {
        for (header of deviceHeaders) {
          const device = instantiateDevice(header);
          console.log(`Controlling light '${device.identifier}'`);
          store.commit("setDevice", device);
        }
      })
      .catch(err => {
        console.error("Failed to intialise devices:", err);
      });

    /* Initialise controller to add routes */
    new DeviceController(service.app, store);
    new HueBridgeController(service.app, hueClient);

    /* Add an error handler that returns valid JSON */
    service.app.use(function(err, req, res, next) {
      console.error(err.stack);
      res.status(500);
      res.json({message: err.message});
    });

    /* Start the server */
    service.listen();

    /* Poll for state changes */
    if (service.config.get("polling.enabled", false)) {
      console.log("Polling for state changes");

      let pollingTimer = setInterval(() => {
        const devices = store.get("devices");
        for (let id in devices) {
          console.log(`Refreshing state for '${devices[id].identifier}'`);
          store.dispatch("fetchDevice", devices[id]).catch(err => {
            console.error(
              `Failed to refresh state for '${devices[id].identifier}':`,
              err
            );
          });
        }
      }, service.config.get("polling.interval", 30000));
    }
  }).catch(err => {
    console.error("Error initialising service", err);
  });

const instantiateDevice = header => {
  let device = new HueLight({
    identifier: header.identifier,
    name: header.name,
    type: header.type,
    controllerName: service.controllerName,
    hueId: header.attributes.hueId
  });

  for (let feature of header.attributes.features) {
    switch (feature) {
      case "color":
        decorate(device, colorDecorator);
        break;
      case "color-temp":
        decorate(device, colorTempDecorator);
        break;
      case "rgb":
        decorate(device, rgbDecorator);
        break;
      default:
        console.error(`Unknown light feature: ${feature}`);
    }
  }

  return device;
};
