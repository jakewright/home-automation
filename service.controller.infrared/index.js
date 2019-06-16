const bootstrap = require("../libraries/javascript/bootstrap");
const router = require("../libraries/javascript/router");
const dao = require("./dao");
require("./routes");

const serviceName = "service.controller.infrared";
bootstrap(serviceName)
  .then(() => {
    return dao.fetchAllState()
  })
  .then(() => {
    router.listen();
  })
  .catch(err => {
    console.error("Error initialising service", err)
  });