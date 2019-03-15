const dao = require("../dao");

/** Middleware to load the device */
const load = (req, res, next) => {
  req.device = dao.findById(req.params.deviceId);

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
  let err = req.device.validate(req.body);
  if (err !== undefined) {
    res.status(400);
    res.json({ message: `Invalid state: ${err}` });
    return;
  }

  const state = req.device.transform(req.body);

  try {
    await dao.applyState(req.device, state);
    res.json({ data: req.device });
  } catch (err) {
    next(err);
  }
};

exports = module.exports = { load, get, update };
