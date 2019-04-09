const router = require("../../libraries/javascript/router");
const device = require("./device");

router.use("/device/:deviceId", device.load);
router.get("/device/:deviceId", device.get);
router.patch("/device/:deviceId", device.update);
