export default class DeviceController {
  constructor(express, lights) {
    this.lights = lights;

    /* Middleware */
    express.use('/device/:deviceId', this.loadDevice.bind(this));

    /* Routes */
    express.get('/device/:deviceId', this.retrieveDevice.bind(this));
    express.patch('/device/:deviceId', this.updateDevice.bind(this));
  }

  /**
   * Middleware to load the device and update it's state
   */
  async loadDevices(req, res, next) {
    req.device = this.lights[req.params.deviceId];

    if (!req.device) {
      res.status(404);
      res.json({ message: 'Device not found' })
      return;
    }

    const state = await req.device.fetchRemoteState();
    req.device.applyRemoteState(state);
    next();
  }

  /**
   * Retrieve a device's current state
   */
  retrieveDevice(req, res) {
    res.json({ message: 'Device', data: req.device });
  }

  /**
   * Update a device. Only properties that are set will be updated.
   */
  async updateDevice(req, res) {
    req.device.setState(req.body);
    await req.device.save();
    res.json({ message: 'Updated device', data: req.device });
  }
}
