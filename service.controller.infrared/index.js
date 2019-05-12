const bootstrap = require("../libraries/javascript/bootstrap");

const serviceName = "service.controller.infrared";
bootstrap(serviceName)
  .catch(err => {
    console.error("Error initialising service", err)
  });