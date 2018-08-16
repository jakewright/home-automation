import Service from '../libraries/javascript/bootstrap';
import HueLight from './domain/HueLight';
import ColorDecorator from './domain/ColorDecorator';
import ColorTempDecorator from './domain/ColorTempDecorator';
import RgbDecorator from './domain/RgbDecorator';
import HueBridgeClient from './api/HueBridgeClient';
import DeviceController from './controllers/DeviceController';
import HueDiscoveryController from './controllers/HueDiscoveryController';

const service = new Service();

// Create an instance of the Huejay client
const client = new HueBridgeClient();

// Get the devices from the registry
const devices = await service.apiClient.getDevices(service.controllerName);

let lights = [];

for (let device of devices) {
  device.client = client;
  let light = new HueLight(device);

  for (let feature of device.features) {
    switch (feature) {
      case 'color':
        ColorDecorator(light);
      case 'color-temp':
        ColorTempDecorator(light);
      case 'rgb':
        RgbDecorator(light);
      default:
        console.error(`Unknown light feature: ${feature}`);
    }
  }

  lights.push(light);
}

/* Initialise controllers to add routes */
new DeviceController(service.app, lights);
new HueDiscoveryController(service.app, client);

app.use(function (err, req, res, next) {
    console.error(err.stack);
    res.status(500);
    res.json({
        message: 'An error occurred while updating the device.',
        errors: [err.message],
    });
});

service.listen();
