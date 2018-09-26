const Service = require("../libraries/javascript/bootstrap");
const HueLight = require("./domain/HueLight");
const decorate = require("./domain/decorate");
const colorDecorator = require("./domain/colorDecorator");
const colorTempDecorator = require("./domain/colorTempDecorator");
const rgbDecorator = require("./domain/rgbDecorator");
const HueClient = require("./api/HueClient");
const HueService = require("./service/HueService");
const DeviceController = require("./controllers/DeviceController");
const HueDiscoveryController = require("./controllers/HueDiscoveryController");
const Store = require("./store/Store");

(async () => {
  try {
    /* Use the bootstrap library to create a Service object */
    const service = new Service();

    /* Get the devices from the registry */
    const deviceHeaders = await service.apiClient.getDevices(
      service.controllerName
    );

    /* Create and decorate the HueLight objects */
    let devices = {};
    for (let header of deviceHeaders) {
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

      console.log(`Controlling light '${device.identifier}'`);
      devices[device.identifier] = device;
    }

    const store = new Store(devices);
    store.on("device-state-changed", state => {
      console.log(`State changed for device ${state.identifier}`);
      service.redisClient.publish(
        `device-state-changed.${state.identifier}`,
        JSON.stringify(state)
      );
    });

    const client = new HueClient({
      host: service.config.get("hueBridge.host"),
      username: service.config.get("hueBridge.username")
    });


    const hueService = new HueService(store, client);
    hueService.startPolling();

    /* Initialise controllers to add routes */
    new DeviceController(service.app, hueService);
    new HueDiscoveryController(service.app, client);

    /* Add an error handler that returns valid JSON */
    service.app.use(function(err, req, res, next) {
      console.error(err.stack);
      res.status(500);
      res.json({
        message: "An error occurred",
        error: err.message
      });
    });

    /* Start the server */
    service.listen();
  } catch (err) {
    console.error("Failed to intialise service");
    console.error(err);
  }
})();
