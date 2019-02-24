const router = require("../../libraries/javascript/router");
const device = require("./device");
const bridge = require("./bridge");

router.use("/device/:deviceId", device.load);
router.get("/device/:deviceId", device.get);
router.patch("/device/:deviceId", device.update);

router.get("/hue/discover", bridge.discover);
router.post("/hue/bridge/users", bridge.createUser);
router.get("/hue/bridge/users", bridge.getAllUsers);
router.get("/hue/bridge/lights", bridge.getAllLights);
