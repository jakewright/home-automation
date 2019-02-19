const plug = require("../plug");

/** Middleware to load the device */
const load = (req, res, next) => {
  req.device = plug.findById(req.params.deviceId);

  if (!req.device) {
    res.status(404);
    res.json({ message: "Device not found" });
    return;
  }

  next();
};

/** Retrieve a device's current state */
const get = (req, res) => {
  res.json({ data: req.device });
};

/** Update a device. Only properties that are set will be updated. */
const update = async (req, res, next) => {
  const state = req.device.transform(req.body);

  try {
    await plug.applyState(req.device, state);
    res.json({ data: req.device });
  } catch (err) {
    next(err);
  }
};

exports = module.exports = { load, get, update };
