class DeviceController {
  constructor(express, store) {
    this.store = store;

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
    req.device = this.store.get("device", req.params.deviceId);

    if (!req.device) {
      res.status(404);
      res.json({ message: "Device not found" });
      return;
    }

    next();
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
    let state;
    try {
      state = req.device.transform(req.body);
    } catch (err) {
      res.status(400);
      res.json({ message: "Failed to transform state", error: err.message });
      return;
    }

    this.store
      .dispatch("saveDevice", { deviceId: req.device, state })
      .then(() => {
        res.json({ message: "Updated device", data: req.device });
      })
      .catch(next);
  }
}

exports = module.exports = DeviceController;
