const Service = require("../libraries/javascript/bootstrap");
const HueLight = require("./domain/HueLight");
const ColorDecorator = require("./domain/ColorDecorator");
const ColorTempDecorator = require("./domain/ColorTempDecorator");
const RgbDecorator = require("./domain/RgbDecorator");
const HueBridgeClient = require("./api/HueBridgeClient");
const DeviceController = require("./controllers/DeviceController");
const HueDiscoveryController = require("./controllers/HueDiscoveryController");

(async () => {
  try {
    /* Use the bootstrap library to create a Service object */
    const service = new Service();

    /* Create an instance of the Hue Bridge client */
    const client = new HueBridgeClient({
      host: service.config.get("hueBridge.host"),
      username: service.config.get("hueBridge.username")
    });

    /* Get the devices from the registry */
    const devices = await service.apiClient.getDevices(service.controllerName);

    /* Create and decorate the HueLight objects */
    let lights = {};
    for (let device of devices) {
      let light = new HueLight({
        identifier: device.identifier,
        name: device.name,
        type: device.type,
        controllerName: service.controllerName,
        hueId: device.attributes.hueId,
        client: client
      });

      for (let feature of device.attributes.features) {
        switch (feature) {
          case "color":
            ColorDecorator(light);
            break;
          case "color-temp":
            ColorTempDecorator(light);
            break;
          case "rgb":
            RgbDecorator(light);
            break;
          default:
            console.error(`Unknown light feature: ${feature}`);
        }
      }

      console.log(`Controlling light '${light.identifier}'`);
      lights[light.identifier] = light;
    }

    /* Initialise controllers to add routes */
    new DeviceController(service.app, lights);
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
