const Service = require("../libraries/javascript/bootstrap");
const { DeviceStore } = require("../libraries/javascript/device");

const LightService = require("./service/LightService");
const HueClient = require("./api/HueClient");

const HueLight = require("./domain/HueLight");
const colorDecorator = require("./domain/colorDecorator");
const colorTempDecorator = require("./domain/colorTempDecorator");
const rgbDecorator = require("./domain/rgbDecorator");
const decorateDevice = require("./domain/decorateDevice");

const DeviceController = require("./controller/DeviceController");
const HueBridgeController = require("./controller/HueBridgeController");

// Create and initialise a Service object
const service = new Service("service.controller.hue");
service
  .init()
  .then(() => {
    // Get device headers
    return service.apiClient.get("service.registry.device/devices", {
      controllerName: service.controllerName
    });
  })
  .then(deviceHeaders => {
    // Instantiate devices and create store
    const devices = deviceHeaders.map(instantiateDevice);
    const store = new DeviceStore(devices);

    // Subscribe to state changes from the store
    store.on("device-changed", (identifier, oldState, newState) => {
      console.log(`State changed for device ${identifier}`);
      service.redisClient.publish(
        `device-state-changed.${identifier}`,
        JSON.stringify({ oldState, newState })
      );
    });

    const hueClient = new HueClient({
      host: service.config.get("hueBridge.host"),
      username: service.config.get("hueBridge.username")
    });

    const lightService = new LightService(store, service.apiClient, hueClient);
    lightService.fetchAllState().catch(err => {
      console.error("Failed to fetch state", err);
    });

    // Initialise controllers to add routes
    new DeviceController(service.app, lightService);
    new HueBridgeController(service.app, hueClient);

    // Start the server
    service.listen();

    // Poll for state changes
    if (service.config.get("polling.enabled", false)) {
      console.log("Polling for state changes");

      let pollingTimer = setInterval(() => {
        lightService.fetchAllState().catch(err => {
          console.error("Failed to refresh state", err);
        });
      }, service.config.get("polling.interval", 30000));
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
    controllerName: service.controllerName,
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
