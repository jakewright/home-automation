const router = require("../../libraries/javascript/router");
const dao = require("../dao");
const ir = require("../ir");

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
  let state;

  try {
    state = req.device.conform(req.body);
  } catch (err) {
    res.status(400);
    res.json({ message: `Failed to validate state: ${err.message}` });
    return;
  }

  try {
    await ir.execute(req.device.generateInstructions(state));
    res.json({ data: req.device });
  } catch (err) {
    next(err);
  }
};

router.use("/device/:deviceId", load);
router.get("/device/:deviceId", get);
router.patch("/device/:deviceId", update);