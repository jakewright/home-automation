class DeviceController {
  constructor(express, store) {
    this.store = store;

    /* Middleware */
    express.use("/device/:deviceId", this.loadDevice.bind(this));

    /* Routes */
    express.get("/device/:deviceId", this.readDevice.bind(this));
    express.patch("/device/:deviceId", this.updateDevice.bind(this));
  }

  /** Middleware to load the device and update it's state */
  loadDevice(req, res, next) {
    req.device = this.store.get("device", req.params.deviceId);

    if (!req.device) {
      res.status(404);
      res.json({ message: "Device not found" });
      return;
    }

    next();
  }

  /** Retrieve a device's current state */
  readDevice(req, res) {
    res.json({ data: req.device });
  }

  /** Update a device. Only properties that are set will be updated. */
  async updateDevice(req, res, next) {
    let state;
    try {
      state = req.device.transform(req.body);
    } catch (err) {
      res.status(400);
      res.json({ message: `Failed to transform state: ${err.message}` });
      return;
    }

    try {
      await this.store.dispatch("saveDevice", { device: req.device, state });
      res.json({ data: req.device });
    } catch (err) {
      next(err);
    }
  }
}

exports = module.exports = DeviceController;
