class DeviceController {
  constructor(express, lightService) {
    this.lightService = lightService;

    // Middleware
    express.use("/device/:deviceId", this.loadDevice.bind(this));

    // Routes
    express.get("/device/:deviceId", this.readDevice.bind(this));
    express.patch("/device/:deviceId", this.updateDevice.bind(this));
  }

  /** Middleware to load the device */
  loadDevice(req, res, next) {
    req.device = this.lightService.findById(req.params.deviceId);

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
    let err = req.device.validate(req.body);
    if (err !== undefined) {
      res.status(422);
      res.json({ message: `Invalid state: ${err}` });
      return;
    }

    const state = req.device.transform(req.body);

    try {
      await this.lightService.applyState(req.device, state);
      res.json({ data: req.device });
    } catch (err) {
      next(err);
    }
  }
}

exports = module.exports = DeviceController;
