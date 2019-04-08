const req = require("../http");
const config = require("../config");

exports = module.exports = async serviceName => {
  // Initialise req for making requests to other services
  const apiGateway = process.env.API_GATEWAY;
  if (!apiGateway) throw new Error("API_GATEWAY env var not set");
  req.setApiGateway(apiGateway);

  // Load config
  const configContents = await req.get(`service.config/read/${serviceName}`);
  config.setContents(configContents);
};
