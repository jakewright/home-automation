class DeviceController {
  constructor(express, lights) {
    this.lights = lights;

    /* Middleware */
    express.use("/device/:deviceId", this.loadDevice.bind(this));

    /* Routes */
    express.get("/device/:deviceId", this.retrieveDevice.bind(this));
    express.patch("/device/:deviceId", this.updateDevice.bind(this));
  }

  /**
   * Middleware to load the device and update it's state
   */
  loadDevice(req, res, next) {
    req.device = this.lights[req.params.deviceId];

    if (!req.device) {
      res.status(404);
      res.json({ message: "Device not found" });
      return;
    }

    req.device
      .fetchRemoteState()
      .then(req.device.applyRemoteState.bind(req.device))
      .then(next)
      .catch(next);
  }

  /**
   * Retrieve a device's current state
   */
  retrieveDevice(req, res) {
    res.json({ message: "Device", data: req.device });
  }

  /**
   * Update a device. Only properties that are set will be updated.
   */
  updateDevice(req, res, next) {
    req.device.setState(req.body)

    req.device.save()
    .then(() => {
        res.json({ message: "Updated device", data: req.device });
      })
      .catch(next);
  }
}

exports = module.exports = DeviceController;
