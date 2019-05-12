const dao = require("../dao");
const _ = require("lodash");

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
  const state = req.device.transform(req.body);

  try {
    await dao.applyState(req.device, state);
    res.json({ data: req.device });
  } catch (err) {
    next(err);
  }
};

const provideState = (req, res) => {
  for (const device of dao.findAll()) {
    // Skip if it doesn't have devices in its attributes
    if (!("devices" in device.attributes)) continue;

    // When devices are retrieved from the device registry, all of the keys in
    // the response are converted to CamelCase, which annoyingly includes the
    // devices for which we provide state. The easy fix here is to also convert
    // the device ID in the request to CamelCase before looking it up.
    const id = _.camelCase(req.params.deviceId);

    // Skip if the ID of the device for which we're providing state
    // doesn't appear in the devices attribute
    if (!(id in device.attributes.devices)) continue;

    const power = device.combination.includes(id);
    res.json({
      data: {
        state_provider: device.identifier,
        // Use the original ID in the request, not the CamelCase version.
        identifier: req.params.deviceId,
        power: power
      }
    });
    return;
  }

  res.status(404);
  res.json({ message: "ID not found in any device's attributes" });
};

exports = module.exports = { load, get, update, provideState };
